package schedulers

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"iflytek.com/weipan4/learn-go/rag/api/service"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/embedding"
	"iflytek.com/weipan4/learn-go/rag/pkg/constants"
	"iflytek.com/weipan4/learn-go/rag/pkg/loggers"
)

type Embedding struct {
	DocumentId string
	segments   []EmbeddingSegment
}
type EmbeddingSegment struct {
	segmentId  string
	embeddings []float32
}
type EmbeddingResp struct {
	segmentId  string
	status     string
	embeddings []float32
}

type EmbeddingScheduler struct {
	embeddingClient embedding.EmbeddingClient
	docService      service.DocumentService
	segmentService  service.DocumentSegmentService
	messageChan     <-chan Message
	embeddingChan   chan<- Embedding
}

func NewEmbeddingScheduler(embeddingClient embedding.EmbeddingClient, docService service.DocumentService, segmentService service.DocumentSegmentService, messageChan <-chan Message, embeddingChan chan<- Embedding) *EmbeddingScheduler {
	return &EmbeddingScheduler{
		embeddingClient: embeddingClient,
		docService:      docService,
		segmentService:  segmentService,
		messageChan:     messageChan,
		embeddingChan:   embeddingChan,
	}
}

func (es *EmbeddingScheduler) startEmbeddingTextTask(messageChan <-chan Message) {
	for message := range messageChan {
		es.embeddingText(message)
	}
}

func (es *EmbeddingScheduler) embeddingText(message Message) {
	// 1. 遍历文本进行向量化
	docId := message.DocumentId
	messageSegments := message.segments
	segmentIds := make([]string, 0, len(messageSegments))
	for _, msgSeg := range messageSegments {
		segmentIds = append(segmentIds, msgSeg.SegmentId)
	}
	// 2. 更新Document还有所有segment的状态为Embedding
	if err := es.docService.UpdateStatus(docId, constants.DocumentEmbedding); err != nil {
		loggers.Error("update document status to embedding failed", zap.String("docId", docId), zap.Error(err))
		return
	}
	loggers.Info("document embedding task start", zap.String("docId", docId))
	if err := es.segmentService.BatchUpdateStatus(segmentIds, constants.SegmentEmbedding); err != nil {
		loggers.Error("update segments status to embedding failed", zap.Error(err))
		return
	}
	loggers.Info("segments embedding task start", zap.Strings("segmentIds", segmentIds))

	// 3. 对每个Segment进行并发向量化
	var wg sync.WaitGroup
	respCh := make(chan EmbeddingResp)
	for _, msgSeg := range messageSegments {
		if msgSeg.Content == "" {
			respCh <- EmbeddingResp{
				segmentId:  msgSeg.SegmentId,
				status:     constants.SegmentEmbeddingCompleted,
				embeddings: make([]float32, 0),
			}
			continue
		}
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancelFunc()
			embeddings, err := es.embeddingClient.MockEmbedding(ctx, text)
			if err != nil {
				loggers.Error("embedding for text failed", zap.Error(err))
				respCh <- EmbeddingResp{
					segmentId:  msgSeg.SegmentId,
					status:     constants.SegmentEmbeddingFailed,
					embeddings: nil,
				}
				return
			}
			respCh <- EmbeddingResp{
				segmentId:  msgSeg.SegmentId,
				status:     constants.SegmentEmbeddingCompleted,
				embeddings: embeddings,
			}
		}(msgSeg.Content)
	}
	go func() {
		wg.Wait()
		close(respCh)
	}()

	var completeIds, failIds []string
	embeddingSegments := make([]EmbeddingSegment, 0)
	for resp := range respCh {
		switch resp.status {
		case constants.SegmentEmbeddingCompleted:
			completeIds = append(completeIds, resp.segmentId)
			embeddingSegments = append(embeddingSegments, EmbeddingSegment{
				segmentId:  resp.segmentId,
				embeddings: resp.embeddings,
			})
		case constants.SegmentEmbeddingFailed:
			failIds = append(failIds, resp.segmentId)
		}
	}
	loggers.Info("segments embedding task end", zap.Strings("completed segment ids", completeIds),
		zap.Strings("failed segment ids", failIds))

	// 4. 更新Document和Segment的状态
	if err := es.segmentService.BatchUpdateStatus(completeIds, constants.SegmentEmbeddingCompleted); err != nil {
		loggers.Error("update segments status to completed failed", zap.Error(err))
		return
	}
	if err := es.segmentService.BatchUpdateStatus(failIds, constants.SegmentEmbeddingFailed); err != nil {
		loggers.Error("update segments status to failed failed", zap.Error(err))
		return
	}
	var newStatus = constants.DocumentEmbeddingCompleted
	if len(failIds) > 0 {
		newStatus = constants.DocumentEmbeddingFailed
	}
	if err := es.docService.UpdateStatus(docId, newStatus); err != nil {
		loggers.Error("update document status failed",
			zap.String("docId", docId), zap.String("new status", newStatus), zap.Error(err))
		return
	}

	// 5. 将embedding结果存入消息队列中
	es.embeddingChan <- Embedding{
		DocumentId: docId,
		segments:   embeddingSegments,
	}
	loggers.Info("segments embedding finished", zap.String("docId", docId), zap.Strings("segmentIds", completeIds))
}
