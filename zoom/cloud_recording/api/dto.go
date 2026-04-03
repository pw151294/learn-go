package api

import "iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/models"

// UploadRequest 上传录制文件请求（multipart form）
type UploadRequest struct {
	Title    string `form:"title" binding:"required"`
	UserID   string `form:"user_id" binding:"required"`
	RoomID   string `form:"room_id" binding:"required"`
	Duration int    `form:"duration"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	RecordingID string `json:"recording_id"`
	PlayURL     string `json:"play_url"`
}

// PlayURLRequest 生成播放 URL 请求
type PlayURLRequest struct {
	Expiry int `json:"expiry"` // 秒，0 表示使用默认值
}

// PlayURLResponse 播放 URL 响应
type PlayURLResponse struct {
	URL string `json:"url"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Total   int                `json:"total"`
	Records []*models.Recording `json:"records"`
}
