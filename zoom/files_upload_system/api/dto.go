// Package api 定义 HTTP 请求/响应的数据传输对象（DTO）。
package api

// InitUploadRequest 初始化上传请求体。
type InitUploadRequest struct {
	FileName    string `json:"file_name"    binding:"required"`
	FileSize    int64  `json:"file_size"    binding:"required,min=1"`
	FileMD5     string `json:"file_md5"     binding:"required,len=32"` // 32 位十六进制 MD5
	TotalChunks int    `json:"total_chunks" binding:"required,min=1"`
}

// CompleteUploadRequest 完成上传请求体。
type CompleteUploadRequest struct {
	TaskID string `json:"task_id" binding:"required"`
}

// Response 统一 JSON 响应包装。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 返回成功响应。
func OK(data interface{}) Response {
	return Response{Code: 0, Message: "ok", Data: data}
}

// Fail 返回失败响应。
func Fail(msg string) Response {
	return Response{Code: -1, Message: msg}
}
