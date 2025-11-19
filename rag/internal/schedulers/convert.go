package schedulers

import (
	"time"

	"github.com/google/uuid"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/weaviate"
	"iflytek.com/weipan4/learn-go/rag/pkg/constants"
	"iflytek.com/weipan4/learn-go/rag/pkg/utils"
)

func AssembleDocumentSegment(document models.Document, text string) models.DocumentSegment {
	segment := models.DocumentSegment{}
	segment.ID = uuid.New().String()
	segment.DatasetID = document.DatasetID
	segment.DocumentID = document.ID
	segment.SubjectID = document.SubjectID
	segment.Content = text
	segment.IndexNodeID = uuid.New().String()
	segment.IndexNodeHash = utils.Sha256Hex(text)
	segment.Enabled = true
	segment.Status = constants.SegmentIndexingInit
	segment.CreatedAt = time.Now()
	segment.UpdatedAt = time.Now()
	segment.Scope = document.Scope
	segment.Env = document.Env

	return segment
}

func AssembleEmbeddingRecord(segment models.DocumentSegment, embeddings []float32) weaviate.EmbeddingRecord {
	record := weaviate.EmbeddingRecord{}
	record.DatasetID = segment.DatasetID
	record.DocumentID = segment.DocumentID
	record.SegmentID = segment.ID
	record.IndexNodeID = segment.IndexNodeID
	record.IndexNodeHash = segment.IndexNodeHash
	record.Text = segment.Content
	record.Vector = embeddings

	return record
}
