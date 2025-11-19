package schedulers

import (
	"strings"

	"go.uber.org/zap"
	"iflytek.com/weipan4/learn-go/rag/api/service"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/weaviate"
	"iflytek.com/weipan4/learn-go/rag/pkg/constants"
	"iflytek.com/weipan4/learn-go/rag/pkg/loggers"
)

type EmbeddingStoreScheduler struct {
	docService     service.DocumentService
	segmentService service.DocumentSegmentService
	embeddingChan  <-chan Embedding
}

func NewEmbeddingStoreScheduler(docService service.DocumentService, segmentService service.DocumentSegmentService, embeddingChan <-chan Embedding) *EmbeddingStoreScheduler {
	return &EmbeddingStoreScheduler{
		docService:     docService,
		segmentService: segmentService,
		embeddingChan:  embeddingChan,
	}
}

func (ess *EmbeddingStoreScheduler) startEmbeddingStoreTask(embeddingChan <-chan Embedding) {
	for embedding := range embeddingChan {
		ess.embeddingStore(embedding)
	}
}

func (ess *EmbeddingStoreScheduler) embeddingStore(embedding Embedding) {
	// 1. 从embedding里解析出参数
	docId := embedding.DocumentId
	embeddingSegments := embedding.segments
	segIds := make([]string, 0, len(embeddingSegments))
	embeddingsBySegId := make(map[string][]float32, 0)
	for _, embdSeg := range embeddingSegments {
		segIds = append(segIds, embdSeg.segmentId)
		embeddingsBySegId[embdSeg.segmentId] = embdSeg.embeddings
	}
	documentSegments, err := ess.segmentService.ListByIds(segIds)
	if err != nil {
		loggers.Error("query segments by ids failed", zap.Strings("ids", segIds), zap.Error(err))
		return
	}
	if len(documentSegments) == 0 {
		return
	}

	// 2. 根据Segment构建出records
	datasetID := documentSegments[0].DatasetID
	records := make([]weaviate.EmbeddingRecord, 0, len(documentSegments))
	for _, docSeg := range documentSegments {
		embeddings := embeddingsBySegId[docSeg.ID]
		record := AssembleEmbeddingRecord(docSeg, embeddings)
		records = append(records, record)
	}

	// 3. 存储embedding的结果
	collectionName := "DATASET_" + strings.ReplaceAll(datasetID, "-", "")
	if err := weaviate.CreateCollection(collectionName); err != nil {
		loggers.Error("create weaviate collection failed", zap.String("name", collectionName), zap.Error(err))
		return
	}
	docNewStatus := constants.DocumentIndexingFinished
	segNewStatus := constants.SegmentIndexingFinished
	err = weaviate.InsertBatch(collectionName, records)
	if err != nil {
		loggers.Error("insert embedding records failed", zap.String("name", collectionName), zap.Error(err))
		docNewStatus = constants.DocumentIndexingFailed
		segNewStatus = constants.SegmentIndexingFailed
	}

	// 4. 更新Segment还有document的状态
	if err := ess.docService.UpdateStatus(docId, docNewStatus); err != nil {
		loggers.Error("update document status failed", zap.String("datasetId", datasetID),
			zap.String("docId", docId), zap.String("new status", docNewStatus), zap.Error(err))
		return
	}
	if err := ess.segmentService.BatchUpdateStatus(segIds, segNewStatus); err != nil {
		loggers.Error("update segments status failed", zap.String("datasetId", datasetID),
			zap.Strings("segmentIds", segIds), zap.String("new status", segNewStatus), zap.Error(err))
	}

	loggers.Info("embedding store end", zap.String("datasetId", datasetID), zap.String("docId", docId),
		zap.Strings("segmentIds", segIds))
}
