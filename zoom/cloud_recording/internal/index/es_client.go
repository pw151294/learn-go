package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
)

// ESClient 封装原生 HTTP 的 ES 客户端
type ESClient struct {
	cfg    *configs.ESConfig
	client *http.Client
}

// NewESClient 创建新的 ES 客户端
func NewESClient(cfg *configs.ESConfig) *ESClient {
	return &ESClient{cfg: cfg, client: &http.Client{}}
}

// EnsureIndex 确保索引存在，不存在则创建
func (e *ESClient) EnsureIndex(ctx context.Context) error {
	_, status, err := e.doRequest(ctx, http.MethodHead, "/"+e.cfg.IndexName, nil)
	if err != nil {
		return fmt.Errorf("check index: %w", err)
	}
	if status == http.StatusOK {
		return nil
	}

	mappings := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"recording_id":     map[string]string{"type": "keyword"},
				"title":            map[string]string{"type": "text"},
				"user_id":          map[string]string{"type": "keyword"},
				"room_id":          map[string]string{"type": "keyword"},
				"bucket":           map[string]string{"type": "keyword"},
				"object_key":       map[string]string{"type": "keyword"},
				"size":             map[string]string{"type": "long"},
				"duration":         map[string]string{"type": "integer"},
				"status":           map[string]string{"type": "keyword"},
				"created_at":       map[string]string{"type": "date"},
				"last_access_at":   map[string]string{"type": "date"},
				"tier_migrated_at": map[string]string{"type": "date"},
			},
		},
	}

	_, status, err = e.doRequest(ctx, http.MethodPut, "/"+e.cfg.IndexName, mappings)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	if status/100 != 2 {
		return fmt.Errorf("create index failed with status %d", status)
	}
	return nil
}

// doRequest 发送 HTTP 请求，返回响应 body、状态码和错误
func (e *ESClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	url := e.cfg.URL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if e.cfg.Username != "" {
		req.SetBasicAuth(e.cfg.Username, e.cfg.Password)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response body: %w", err)
	}
	return respBody, resp.StatusCode, nil
}

// IndexDocument 写入或覆盖文档
func (e *ESClient) IndexDocument(ctx context.Context, id string, doc interface{}) error {
	path := fmt.Sprintf("/%s/_doc/%s", e.cfg.IndexName, id)
	_, status, err := e.doRequest(ctx, http.MethodPut, path, doc)
	if err != nil {
		return fmt.Errorf("index document: %w", err)
	}
	if status/100 != 2 {
		return fmt.Errorf("index document failed with status %d", status)
	}
	return nil
}

// GetDocument 获取文档，返回 _source 字段
func (e *ESClient) GetDocument(ctx context.Context, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/%s/_doc/%s", e.cfg.IndexName, id)
	data, status, err := e.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status/100 != 2 {
		return nil, fmt.Errorf("get document failed with status %d", status)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	source, _ := result["_source"].(map[string]interface{})
	return source, nil
}

// UpdateDocument 局部更新文档字段
func (e *ESClient) UpdateDocument(ctx context.Context, id string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/%s/_update/%s", e.cfg.IndexName, id)
	body := map[string]interface{}{"doc": fields}
	_, status, err := e.doRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	if status/100 != 2 {
		return fmt.Errorf("update document failed with status %d", status)
	}
	return nil
}

// DeleteDocument 删除文档
func (e *ESClient) DeleteDocument(ctx context.Context, id string) error {
	path := fmt.Sprintf("/%s/_doc/%s", e.cfg.IndexName, id)
	_, status, err := e.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	if status/100 != 2 {
		return fmt.Errorf("delete document failed with status %d", status)
	}
	return nil
}

// Search 执行搜索，返回 _source 列表和总命中数
func (e *ESClient) Search(ctx context.Context, query map[string]interface{}) ([]map[string]interface{}, int, error) {
	path := fmt.Sprintf("/%s/_search", e.cfg.IndexName)
	data, status, err := e.doRequest(ctx, http.MethodPost, path, query)
	if err != nil {
		return nil, 0, fmt.Errorf("search: %w", err)
	}
	if status/100 != 2 {
		return nil, 0, fmt.Errorf("search failed with status %d", status)
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source map[string]interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, 0, fmt.Errorf("unmarshal search result: %w", err)
	}

	sources := make([]map[string]interface{}, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		sources = append(sources, hit.Source)
	}
	return sources, result.Hits.Total.Value, nil
}
