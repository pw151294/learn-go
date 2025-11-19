package schedulers

import (
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"iflytek.com/weipan4/learn-go/rag/api/service"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
	"iflytek.com/weipan4/learn-go/rag/pkg/constants"
	"iflytek.com/weipan4/learn-go/rag/pkg/loggers"
	"iflytek.com/weipan4/learn-go/rag/pkg/utils"
)

type MessageSegment struct {
	SegmentId string
	Content   string
}
type Message struct {
	DocumentId string
	segments   []MessageSegment
}

// SplittingScheduler 文本切分调度器
type SplittingScheduler struct {
	docService     service.DocumentService
	segmentService service.DocumentSegmentService
	messageChan    chan<- Message
}

func NewSplittingScheduler(docService service.DocumentService, segmentService service.DocumentSegmentService, messageChan chan<- Message) *SplittingScheduler {
	return &SplittingScheduler{
		docService:     docService,
		segmentService: segmentService,
		messageChan:    messageChan,
	}
}

func (ss *SplittingScheduler) StartSplittingTextTask(period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			loggers.Info("begin text splitting task")
			ss.splitText()
		}
	}
}

func (ss *SplittingScheduler) splitText() {
	// 1. 查询出所有Waiting文档的元数据
	docs, err := ss.docService.WaitingDocuments(constants.BatchDocumentSize)
	if err != nil {
		loggers.Error("query waiting documents failed", zap.Error(err))
		return
	}

	// 2. 查询出文档的元数据
	for _, doc := range docs {
		savePath := filepath.Join(constants.UploadDir, doc.ID+doc.DocForm)
		file, err := utils.ReadFile(savePath)
		if err != nil {
			loggers.Error("read file failed", zap.String("path", savePath), zap.Error(err))
			continue
		}
		content := string(file)
		texts := utils.SplitText(content)
		documentSegments := make([]models.DocumentSegment, 0, len(texts))
		messageSegments := make([]MessageSegment, 0, len(texts))
		for _, text := range texts {
			docSeg := AssembleDocumentSegment(doc, text)
			documentSegments = append(documentSegments, docSeg)
		}
		if err := ss.segmentService.InsertBatch(documentSegments); err != nil {
			loggers.Error("batch create document segments failed", zap.Error(err))
			continue
		}
		loggers.Info("create document segments success", zap.String("docName", doc.Name),
			zap.Int("segment count", len(documentSegments)))

		// 3. 采集已经切割好的文本
		for _, docSeg := range documentSegments {
			messageSegments = append(messageSegments, MessageSegment{
				SegmentId: docSeg.ID,
				Content:   docSeg.Content,
			})
		}
		if err := ss.docService.UpdateStatus(doc.ID, constants.DocumentIndexingInit); err != nil {
			loggers.Error("update document status failed", zap.Error(err))
			continue
		}
		loggers.Info("start document indexing", zap.String("docId", doc.ID), zap.String("docName", doc.Name))

		// 4. 将message分段文本投递到chan中
		ss.messageChan <- Message{doc.ID, messageSegments}
	}
}
