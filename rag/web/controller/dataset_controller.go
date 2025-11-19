package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"iflytek.com/weipan4/learn-go/rag/api/service"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
)

type DatasetController struct {
	DatasetService service.DatasetService
}

func NewDatasetController(ds service.DatasetService) *DatasetController {
	return &DatasetController{
		DatasetService: ds,
	}
}

func (dc *DatasetController) Create(c *gin.Context) {
	var req models.Dataset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数解析失败"})
		return
	}
	req.ID = uuid.New().String()
	if err := dc.DatasetService.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "知识库创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "知识库创建成功",
		"id":  req.ID,
	})
}
