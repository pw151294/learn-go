package weaviate

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

// VectorConstant 对应 Java 代码中的常量
const (
	DatasetIDProperty     = "dataset_id"
	DocumentIDProperty    = "document_id"
	SegmentIDProperty     = "segment_id"
	IndexNodeIDProperty   = "index_node_id"
	IndexNodeHashProperty = "index_node_hash"
	TextProperty          = "text"
)

// EmbeddingRecord 用于批量插入的数据结构
type EmbeddingRecord struct {
	DatasetID     string
	DocumentID    string
	SegmentID     string
	IndexNodeID   string
	IndexNodeHash string
	Text          string
	Vector        []float32
}

// EmbeddingMatch 用于向量搜索返回的结果结构
type EmbeddingMatch struct {
	Score         float64
	DatasetID     string
	DocumentID    string
	SegmentID     string
	IndexNodeID   string
	IndexNodeHash string
	Text          string
	ID            string
	Vector        []float32
}

// existCollection 判断名为 collectionName 的 collection 是否存在
func existCollection(collectionName string) bool {
	client := GetWeaviateClient()
	if client == nil {
		return false
	}

	exists, err := client.Schema().ClassExistenceChecker().WithClassName(collectionName).Do(context.Background())
	if err != nil {
		return false
	}

	return exists
}

// CreateCollection 在 Weaviate 中创建一个新的集合（Class）
func CreateCollection(collectionName string) error {
	client := GetWeaviateClient()
	if client == nil {
		return fmt.Errorf("weaviate client is not initialized")
	}

	// 判断 collection 是否存在，存在则直接返回
	if existCollection(collectionName) {
		return nil
	}

	// 定义 Class
	class := &models.Class{
		Class:      collectionName,
		Properties: defaultProperties(),
		Vectorizer: "none", // 不启用向量化，由外部生成后直接存入
	}

	// 创建 Class
	err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create class '%s': %w", collectionName, err)
	}

	return nil
}

// InsertBatch 批量插入 EmbeddingRecord 到 Weaviate
func InsertBatch(collectionName string, records []EmbeddingRecord) error {
	// 1. 参数校验
	if len(records) == 0 {
		return nil
	}

	client := GetWeaviateClient()
	if client == nil {
		return fmt.Errorf("weaviate client is not initialized")
	}

	// 2. 构建对象并批量添加
	objects := make([]*models.Object, len(records))
	for i, record := range records {
		objects[i] = &models.Object{
			Class: collectionName,
			Properties: map[string]interface{}{
				DatasetIDProperty:     record.DatasetID,
				DocumentIDProperty:    record.DocumentID,
				SegmentIDProperty:     record.SegmentID,
				IndexNodeIDProperty:   record.IndexNodeID,
				IndexNodeHashProperty: record.IndexNodeHash,
				TextProperty:          record.Text,
			},
			Vector: record.Vector,
		}
	}

	// 3. 执行批量插入
	batchResult, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		return fmt.Errorf("批量插入到 weaviate 中失败: %w", err)
	}

	// 检查批量操作中是否有单独的错误
	for _, res := range batchResult {
		if res.Result != nil && res.Result.Errors != nil {
			return fmt.Errorf("批量插入时发生错误: %v", res.Result.Errors.Error)
		}
	}

	return nil
}

// SearchByVector 根据向量进行搜索
func SearchByVector(collectionName string, queryVector []float32, topK int, scoreThreshold float64) ([]EmbeddingMatch, error) {
	client := GetWeaviateClient()
	if client == nil {
		return nil, fmt.Errorf("weaviate client is not initialized")
	}

	// 1. 构建查询参数
	nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(queryVector)

	// 定义需要返回的字段
	fields := []graphql.Field{
		{Name: TextProperty},
		{Name: DatasetIDProperty},
		{Name: DocumentIDProperty},
		{Name: SegmentIDProperty},
		{Name: IndexNodeIDProperty},
		{Name: IndexNodeHashProperty},
		{Name: "_additional", Fields: []graphql.Field{
			{Name: "id"},
			{Name: "vector"},
			{Name: "distance"}, // distance 是 weaviate 的一个度量，值越小越相似
		}},
	}

	// 2. 执行查询
	response, err := client.GraphQL().Get().
		WithClassName(collectionName).
		WithFields(fields...).
		WithNearVector(nearVector).
		WithLimit(topK).
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}

	// 3. 解析和过滤结果
	var matches []EmbeddingMatch
	if get, ok := response.Data["Get"].(map[string]interface{}); ok {
		if classData, ok := get[collectionName].([]interface{}); ok {
			for _, item := range classData {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				additional, _ := itemMap["_additional"].(map[string]interface{})
				distance, _ := additional["distance"].(float64)

				// Weaviate 的 distance 越小越好。scoreThreshold 通常指相似度（越大越好）。
				// 这里我们假设 scoreThreshold 是 distance 的上限。
				if distance > scoreThreshold {
					continue
				}

				vectorData, _ := additional["vector"].([]interface{})
				vector := make([]float32, len(vectorData))
				for i, v := range vectorData {
					if val, ok := v.(float64); ok {
						vector[i] = float32(val)
					}
				}

				match := EmbeddingMatch{
					Score:         1 - distance, // 将 distance 转换为 0-1 的相似度分数
					ID:            additional["id"].(string),
					Vector:        vector,
					Text:          itemMap[TextProperty].(string),
					DatasetID:     itemMap[DatasetIDProperty].(string),
					DocumentID:    itemMap[DocumentIDProperty].(string),
					SegmentID:     itemMap[SegmentIDProperty].(string),
					IndexNodeID:   itemMap[IndexNodeIDProperty].(string),
					IndexNodeHash: itemMap[IndexNodeHashProperty].(string),
				}
				matches = append(matches, match)
			}
		}
	}

	return matches, nil
}

// defaultProperties 返回集合的默认属性列表
func defaultProperties() []*models.Property {
	return []*models.Property{
		{
			Name:        DatasetIDProperty,
			DataType:    []string{"text"},
			Description: "知识库id",
		},
		{
			Name:        DocumentIDProperty,
			DataType:    []string{"text"},
			Description: "归属的文档id",
		},
		{
			Name:        SegmentIDProperty,
			DataType:    []string{"text"},
			Description: "归属的文档块id",
		},
		{
			Name:        IndexNodeIDProperty,
			DataType:    []string{"text"},
			Description: "分段唯一标识",
		},
		{
			Name:        IndexNodeHashProperty,
			DataType:    []string{"text"},
			Description: "索引块hash值",
		},
		{
			Name:        TextProperty,
			DataType:    []string{"text"},
			Description: "文本内容",
		},
	}
}
