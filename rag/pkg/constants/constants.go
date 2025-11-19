package constants

const UploadDir = "./uploads"

const (
	DocumentWaitingIndexing    = "WAITING"
	DocumentIndexingInit       = "INDEXING_INIT"
	DocumentEmbedding          = "EMBEDDING"
	DocumentEmbeddingCompleted = "EMBEDDING_COMPLETED"
	DocumentEmbeddingFailed    = "EMBEDDING_FAILED"
	DocumentIndexingFinished   = "INDEXING_FINISHED"
	DocumentIndexingFailed     = "INDEXING_FAILED"
)

const (
	SegmentIndexingInit       = "INDEXING_INIT"
	SegmentEmbedding          = "EMBEDDING"
	SegmentEmbeddingCompleted = "EMBEDDING_COMPLETED"
	SegmentEmbeddingFailed    = "EMBEDDING_FAILED"
	SegmentIndexingFinished   = "INDEXING_FINISHED"
	SegmentIndexingFailed     = "INDEXING_FAILED"
)

const (
	BatchDocumentSize = 10
)
