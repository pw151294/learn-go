package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"iflytek.com/weipan4/learn-go/rag/api/service"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
	"iflytek.com/weipan4/learn-go/rag/pkg/constants"
)

type DocumentController struct {
	DocService     service.DocumentService
	DatasetService service.DatasetService
}

func NewDocumentController(docService service.DocumentService, datasetService service.DatasetService) *DocumentController {
	return &DocumentController{
		DocService:     docService,
		DatasetService: datasetService,
	}
}

// Upload 处理文档上传
func (dc *DocumentController) Upload(c *gin.Context) {
	// 0. 参数校验
	subjectId := c.PostForm("subjectId")
	datasetId := c.PostForm("datasetId")
	if subjectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subjectId不能为空"})
	}
	dataset, err := dc.DatasetService.GetByID(datasetId)
	if err != nil || dataset == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "知识库不存在"})
	}

	// 1. 解析文件和基本信息（小驼峰命名）
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件未上传"})
		return
	}

	// 2. 保存文件到本地（使用常量）
	saveDir := constants.UploadDir
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建存储目录"})
		return
	}
	fileId := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	fileName := strings.TrimSuffix(file.Filename, ext)
	savePath := filepath.Join(saveDir, fileId+ext)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 3. 保存文档记录
	doc := &models.Document{
		ID:        fileId,
		Name:      fileName,
		SubjectID: subjectId,
		DatasetID: datasetId,
		DocForm:   ext,
		Status:    constants.DocumentWaitingIndexing,
		Enabled:   true,
	}
	if err := dc.DocService.Create(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文档记录保存失败"})
		return
	}

	// 4. 返回小驼峰命名的json
	c.JSON(http.StatusOK, gin.H{
		"msg":     "上传成功",
		"fileId":  fileId,
		"fileUrl": savePath,
	})
}
