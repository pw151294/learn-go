package index

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/models"
)

// RecordingIndex 提供 Recording 的 ES 索引操作
type RecordingIndex struct {
	es *ESClient
}

// NewRecordingIndex 创建新的 RecordingIndex
func NewRecordingIndex(es *ESClient) *RecordingIndex {
	return &RecordingIndex{es: es}
}

// recordingToMap 将 Recording 转为 map，时间字段使用 RFC3339 格式
func recordingToMap(rec *models.Recording) map[string]interface{} {
	return map[string]interface{}{
		"recording_id":     rec.RecordingID,
		"title":            rec.Title,
		"user_id":          rec.UserID,
		"room_id":          rec.RoomID,
		"bucket":           rec.Bucket,
		"object_key":       rec.ObjectKey,
		"size":             rec.Size,
		"duration":         rec.Duration,
		"status":           rec.Status,
		"created_at":       rec.CreatedAt.Format(time.RFC3339),
		"last_access_at":   rec.LastAccessAt.Format(time.RFC3339),
		"tier_migrated_at": rec.TierMigratedAt.Format(time.RFC3339),
	}
}

// Create 在 ES 中创建录制文件索引
func (r *RecordingIndex) Create(ctx context.Context, rec *models.Recording) error {
	doc := recordingToMap(rec)
	if err := r.es.IndexDocument(ctx, rec.RecordingID, doc); err != nil {
		return fmt.Errorf("create recording index: %w", err)
	}
	return nil
}

// Update 局部更新索引字段
func (r *RecordingIndex) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if err := r.es.UpdateDocument(ctx, id, fields); err != nil {
		return fmt.Errorf("update recording index: %w", err)
	}
	return nil
}

// UpdateTierStatus 更新存储层级相关字段
func (r *RecordingIndex) UpdateTierStatus(ctx context.Context, id, status, bucket string) error {
	fields := map[string]interface{}{
		"status":           status,
		"bucket":           bucket,
		"tier_migrated_at": time.Now().Format(time.RFC3339),
	}
	if err := r.es.UpdateDocument(ctx, id, fields); err != nil {
		return fmt.Errorf("update tier status: %w", err)
	}
	return nil
}

// Search 根据查询参数搜索录制文件
func (r *RecordingIndex) Search(ctx context.Context, query *models.SearchQuery) ([]*models.Recording, int, error) {
	must := []map[string]interface{}{}
	filter := []map[string]interface{}{}

	if query.Keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query.Keyword,
				"fields": []string{"title"},
			},
		})
	}
	if query.UserID != "" {
		filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"user_id": query.UserID}})
	}
	if query.RoomID != "" {
		filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"room_id": query.RoomID}})
	}
	if query.Status != "" {
		filter = append(filter, map[string]interface{}{"term": map[string]interface{}{"status": query.Status}})
	}

	page, size := query.Page, query.Size
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		"from": (page - 1) * size,
		"size": size,
	}

	sources, total, err := r.es.Search(ctx, esQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("search recordings: %w", err)
	}

	recordings := make([]*models.Recording, 0, len(sources))
	for _, src := range sources {
		rec, err := sourceToRecording(src)
		if err != nil {
			return nil, 0, fmt.Errorf("parse recording: %w", err)
		}
		recordings = append(recordings, rec)
	}
	return recordings, total, nil
}

// ListForMigration 查询指定状态且创建时间早于 beforeTime 的录制文件
func (r *RecordingIndex) ListForMigration(ctx context.Context, status string, beforeTime time.Time, from, size int) ([]*models.Recording, error) {
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{"term": map[string]interface{}{"status": status}},
					{"range": map[string]interface{}{
						"created_at": map[string]interface{}{
							"lt": beforeTime.Format(time.RFC3339),
						},
					}},
				},
			},
		},
		"from": from,
		"size": size,
	}

	sources, _, err := r.es.Search(ctx, esQuery)
	if err != nil {
		return nil, fmt.Errorf("list for migration: %w", err)
	}

	recordings := make([]*models.Recording, 0, len(sources))
	for _, src := range sources {
		rec, err := sourceToRecording(src)
		if err != nil {
			return nil, fmt.Errorf("parse recording: %w", err)
		}
		recordings = append(recordings, rec)
	}
	return recordings, nil
}

// sourceToRecording 将 ES _source map 反序列化为 Recording 结构体
func sourceToRecording(src map[string]interface{}) (*models.Recording, error) {
	data, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("marshal source: %w", err)
	}
	rec := &models.Recording{}
	if err := json.Unmarshal(data, rec); err != nil {
		return nil, fmt.Errorf("unmarshal recording: %w", err)
	}
	return rec, nil
}
