// Package repository 提供对 upload_chunks 表的 CRUD 操作。
package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/model"
)

// ChunkRepo 封装 upload_chunks 表的数据库操作。
type ChunkRepo struct {
	db *gorm.DB
}

func NewChunkRepo(db *gorm.DB) *ChunkRepo {
	return &ChunkRepo{db: db}
}

// FindByTaskID 查询某任务所有已成功上传的分片，按 chunk_index 升序。
func (r *ChunkRepo) FindByTaskID(taskID string) ([]model.UploadChunk, error) {
	var chunks []model.UploadChunk
	err := r.db.Where("task_id = ? AND status = ?", taskID, model.ChunkStatusUploaded).
		Order("chunk_index ASC").
		Find(&chunks).Error
	return chunks, err
}

// FindChunk 查询特定分片，不存在时返回 nil。
func (r *ChunkRepo) FindChunk(taskID string, chunkIndex int) (*model.UploadChunk, error) {
	var chunk model.UploadChunk
	err := r.db.Where("task_id = ? AND chunk_index = ?", taskID, chunkIndex).First(&chunk).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &chunk, err
}

// Upsert 写入或更新分片记录（幂等）。
// 冲突键为 (task_id, chunk_index)，冲突时更新 etag、status、updated_at。
func (r *ChunkRepo) Upsert(chunk *model.UploadChunk) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "task_id"}, {Name: "chunk_index"}},
		DoUpdates: clause.AssignmentColumns([]string{"etag", "status", "updated_at"}),
	}).Create(chunk).Error
}

// CountUploaded 统计某任务已成功上传的分片数量。
func (r *ChunkRepo) CountUploaded(taskID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.UploadChunk{}).
		Where("task_id = ? AND status = ?", taskID, model.ChunkStatusUploaded).
		Count(&count).Error
	return count, err
}
