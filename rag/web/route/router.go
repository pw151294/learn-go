package route

import (
	"github.com/gin-gonic/gin"
	"iflytek.com/weipan4/learn-go/rag/api/service"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/mysql"
	"iflytek.com/weipan4/learn-go/rag/web/controller"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	docRepo := models.NewDocumentRepo(mysql.DB)
	docService := service.NewDocumentService(docRepo)
	datasetRepo := models.NewDatasetRepo(mysql.DB)
	datasetService := service.NewDatasetServiceImpl(datasetRepo)

	datasetController := controller.NewDatasetController(datasetService)
	docController := controller.NewDocumentController(docService, datasetService)

	api := r.Group("/api/rag")
	{
		api.POST("/document/upload", docController.Upload)
		api.POST("/document/segments/create", func(c *gin.Context) {
			// 创建文档分段处理逻辑
			c.JSON(200, gin.H{"msg": "create document segments"})
		})
		api.POST("/dataset/create", datasetController.Create)
	}

	return r
}
