package embedding

import (
	"context"
	"errors"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"iflytek.com/weipan4/learn-go/rag/configs"
)

// EmbeddingClient defines the interface for creating text embeddings.
type EmbeddingClient interface {
	// CreateEmbedding creates a vector embedding for the given text.
	CreateEmbedding(ctx context.Context, text string) ([]float32, error)
	// CreateEmbeddings creates vector embeddings for the given texts in batch.
	CreateEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
	// GetEmbeddingEngine returns the name of the embedding engine/model.
	GetEmbeddingEngine() string
	// MockEmbedding simulates generating an embedding for the given text.
	MockEmbedding(ctx context.Context, text string) ([]float32, error)
}

// openAIEmbeddingClient is an implementation of EmbeddingClient using the OpenAI API.
type openAIEmbeddingClient struct {
	client *openai.Client
	model  openai.EmbeddingModel
}

// NewOpenAIEmbeddingClient creates a new client for OpenAI embeddings.
func NewOpenAIEmbeddingClient() (EmbeddingClient, error) {
	conf := configs.GetConfig()
	if conf == nil {
		return nil, fmt.Errorf("config is not initialized")
	}
	// 初始化openai客户端
	llmConf := conf.LLM
	client := openai.NewClient(
		option.WithAPIKey(llmConf.ApiKey),
		option.WithBaseURL(llmConf.BaseURL))

	if &client == nil {
		return nil, errors.New("LLM client not initialized")
	}
	return &openAIEmbeddingClient{
		client: &client,
	}, nil
}

func (e *openAIEmbeddingClient) MockEmbedding(ctx context.Context, text string) ([]float32, error) {
	vec := make([]float32, 768)
	for i := range vec {
		vec[i] = float32(i)
	}
	return vec, nil
}

// CreateEmbedding creates a vector embedding for a single text.
func (e *openAIEmbeddingClient) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := e.CreateEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, fmt.Errorf("embedding create failed: %w", err)
	}
	if len(embeddings) == 0 {
		return nil, errors.New("embedding generation returned no vectors")
	}
	return embeddings[0], nil
}

// CreateEmbeddings creates vector embeddings for a slice of texts.
func (e *openAIEmbeddingClient) CreateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, errors.New("no texts provided for embedding")
	}

	// 构建输入参数 - 转换为接口切片
	inputs := make([]any, len(texts))
	for i, text := range texts {
		inputs[i] = text
	}

	// 调用 OpenAI Embedding API
	params := openai.EmbeddingNewParams{
		Input:          openai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: texts},
		Model:          openai.EmbeddingModelTextEmbeddingAda002,
		Dimensions:     openai.Int(1),
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
		User:           openai.String("user-1234"),
	}

	embedding, err := e.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}

	// 提取向量数据
	if len(embedding.Data) != len(texts) {
		return nil, fmt.Errorf("expected %d embeddings, got %d", len(texts), len(embedding.Data))
	}

	result := make([][]float32, len(embedding.Data))
	for i, data := range embedding.Data {
		// 将 float64 转换为 float32
		vec := make([]float32, len(data.Embedding))
		for j, val := range data.Embedding {
			vec[j] = float32(val)
		}
		result[i] = vec
	}

	return result, nil
}

// GetEmbeddingEngine returns the model name used for embeddings.
func (e *openAIEmbeddingClient) GetEmbeddingEngine() string {
	return e.model
}
