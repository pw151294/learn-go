package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/auth"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/cdn"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/index"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/models"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/storage"
)

type Handler struct {
	cfg       *configs.Config
	minio     *storage.MinIOClient
	recIdx    *index.RecordingIndex
	signer    *auth.Signer
	validator *auth.Validator
	rewriter  *cdn.URLRewriter
}

func NewHandler(
	cfg *configs.Config,
	minio *storage.MinIOClient,
	recIdx *index.RecordingIndex,
	signer *auth.Signer,
	validator *auth.Validator,
	rewriter *cdn.URLRewriter,
) *Handler {
	return &Handler{cfg: cfg, minio: minio, recIdx: recIdx, signer: signer, validator: validator, rewriter: rewriter}
}

// UploadRecording POST /api/recordings
// 接收 multipart form：file, title, user_id, room_id, duration
func (h *Handler) UploadRecording(c *gin.Context) {
	var req UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	recordingID := uuid.New().String()
	objectKey := req.UserID + "/" + req.RoomID + "/" + recordingID

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "video/mp4"
	}

	if err := h.minio.UploadToHot(c.Request.Context(), objectKey, file, header.Size, contentType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	rec := &models.Recording{
		RecordingID:    recordingID,
		Title:          req.Title,
		UserID:         req.UserID,
		RoomID:         req.RoomID,
		Bucket:         h.cfg.MinIO.HotBucket,
		ObjectKey:      objectKey,
		Size:           header.Size,
		Duration:       req.Duration,
		Status:         models.TierHot,
		CreatedAt:      now,
		LastAccessAt:   now,
		TierMigratedAt: now,
	}

	if err := h.recIdx.Create(c.Request.Context(), rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expiry := h.cfg.Auth.DefaultExpiry
	playURL := h.signer.GeneratePlayURL(recordingID, expiry)

	c.JSON(http.StatusOK, UploadResponse{RecordingID: recordingID, PlayURL: playURL})
}

// SearchRecordings GET /api/recordings
func (h *Handler) SearchRecordings(c *gin.Context) {
	var query models.SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 10
	}

	recs, total, err := h.recIdx.Search(c.Request.Context(), &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SearchResponse{Total: total, Records: recs})
}

// GetRecording GET /api/recordings/:id
func (h *Handler) GetRecording(c *gin.Context) {
	id := c.Param("id")

	// 更新 last_access_at
	h.recIdx.Update(c.Request.Context(), id, map[string]interface{}{
		"last_access_at": time.Now().Format(time.RFC3339),
	})

	recs, _, err := h.recIdx.Search(c.Request.Context(), &models.SearchQuery{Keyword: id, Page: 1, Size: 1})
	if err != nil || len(recs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, recs[0])
}

// DeleteRecording DELETE /api/recordings/:id
func (h *Handler) DeleteRecording(c *gin.Context) {
	id := c.Param("id")

	recs, _, err := h.recIdx.Search(c.Request.Context(), &models.SearchQuery{Keyword: id, Page: 1, Size: 1})
	if err != nil || len(recs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	rec := recs[0]
	if err := h.minio.DeleteObject(c.Request.Context(), rec.Bucket, rec.ObjectKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.recIdx.UpdateTierStatus(c.Request.Context(), id, "deleted", "")
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GeneratePlayURL POST /api/recordings/:id/play-url
func (h *Handler) GeneratePlayURL(c *gin.Context) {
	id := c.Param("id")
	var req PlayURLRequest
	c.ShouldBindJSON(&req)

	expiry := req.Expiry
	if expiry <= 0 {
		expiry = h.cfg.Auth.DefaultExpiry
	}

	playURL := h.signer.GeneratePlayURL(id, expiry)
	c.JSON(http.StatusOK, PlayURLResponse{URL: playURL})
}

// PlayRecording GET /play/:id
// 验证 token，重定向到实际播放 URL
func (h *Handler) PlayRecording(c *gin.Context) {
	id := c.Param("id")
	token := c.Query("token")
	expires := c.Query("expires")

	if err := h.validator.ValidateToken(id, token, expires); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	recs, _, err := h.recIdx.Search(c.Request.Context(), &models.SearchQuery{Keyword: id, Page: 1, Size: 1})
	if err != nil || len(recs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	rec := recs[0]
	presignedURL, err := h.minio.GetPresignedURL(c.Request.Context(), rec.Bucket, rec.ObjectKey, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	finalURL := h.rewriter.RewriteURL(presignedURL)
	c.Redirect(http.StatusFound, finalURL)
}
