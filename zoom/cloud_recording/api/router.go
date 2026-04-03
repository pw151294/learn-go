package api

import "github.com/gin-gonic/gin"

func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/recordings")
	{
		api.POST("", h.UploadRecording)
		api.GET("", h.SearchRecordings)
		api.GET("/:id", h.GetRecording)
		api.DELETE("/:id", h.DeleteRecording)
		api.POST("/:id/play-url", h.GeneratePlayURL)
	}
	r.GET("/play/:id", h.PlayRecording)

	return r
}
