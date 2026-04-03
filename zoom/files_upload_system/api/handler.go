// Package api 实现 HTTP 路由处理器。
package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/service"
)

// UploadHandler 将 HTTP 请求映射到 UploadService 方法。
type UploadHandler struct {
	svc *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// RegisterRoutes 注册所有上传相关路由到 /api/v1/upload 前缀下。
func (h *UploadHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/upload")
	{
		v1.POST("/init", h.InitUpload)
		v1.POST("/chunk", h.UploadChunk)
		v1.POST("/complete", h.CompleteUpload)
		v1.GET("/progress/:task_id", h.GetProgress)
	}
}

// InitUpload godoc
// POST /api/v1/upload/init
// Body: InitUploadRequest（JSON）
// 返回 task_id、upload_id、已上传分片列表（断点续传）或秒传标志。
func (h *UploadHandler) InitUpload(c *gin.Context) {
	var req InitUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Fail("参数错误: "+err.Error()))
		return
	}

	result, err := h.svc.InitUpload(req.FileName, req.FileSize, req.FileMD5, req.TotalChunks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, OK(result))
}

// UploadChunk godoc
// POST /api/v1/upload/chunk
// Content-Type: multipart/form-data
// 字段: task_id（string）、chunk_index（int）、chunk_md5（string）、file（binary）
func (h *UploadHandler) UploadChunk(c *gin.Context) {
	taskID := c.PostForm("task_id")
	chunkIndexStr := c.PostForm("chunk_index")
	chunkMD5 := c.PostForm("chunk_md5")

	if taskID == "" || chunkIndexStr == "" || chunkMD5 == "" {
		c.JSON(http.StatusBadRequest, Fail("缺少必要参数: task_id, chunk_index, chunk_md5"))
		return
	}

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || chunkIndex < 0 {
		c.JSON(http.StatusBadRequest, Fail("chunk_index 必须为非负整数"))
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Fail("读取分片文件失败: "+err.Error()))
		return
	}
	defer file.Close()

	result, err := h.svc.UploadChunk(taskID, chunkIndex, chunkMD5, file, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, OK(result))
}

// CompleteUpload godoc
// POST /api/v1/upload/complete
// Body: CompleteUploadRequest（JSON）
func (h *UploadHandler) CompleteUpload(c *gin.Context) {
	var req CompleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Fail("参数错误: "+err.Error()))
		return
	}

	fileURL, err := h.svc.CompleteUpload(req.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, OK(gin.H{"file_url": fileURL}))
}

// GetProgress godoc
// GET /api/v1/upload/progress/:task_id
func (h *UploadHandler) GetProgress(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, Fail("task_id 不能为空"))
		return
	}

	progress, err := h.svc.GetProgress(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, OK(progress))
}
