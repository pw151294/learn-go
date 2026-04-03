// zoom/distributed_encoding_platform/dlq/dlq.go
package dlq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/config"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/model"
)

// DLQEntry 是死信队列中的一条记录
type DLQEntry struct {
	Task     *model.Task `json:"task"`
	Reason   string      `json:"reason"`
	FailedAt time.Time   `json:"failed_at"`
}

var (
	entries []DLQEntry
	mu      sync.RWMutex
	writer  *kafka.Writer
)

// Init 初始化 DLQ Kafka Writer
func Init() {
	writer = &kafka.Writer{
		Addr:     kafka.TCP(config.KafkaBroker),
		Topic:    config.TopicDLQ,
		Balancer: &kafka.LeastBytes{},
	}
}

// Push 将失败任务写入 DLQ Topic，并缓存到内存列表
func Push(task *model.Task, reason string) error {
	entry := DLQEntry{Task: task, Reason: reason, FailedAt: time.Now()}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("dlq marshal error: %w", err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(task.ID),
		Value: data,
	})
	if err != nil {
		log.Printf("[DLQ] write error: %v", err)
		// 即使 Kafka 写失败，仍保留内存记录
	}

	mu.Lock()
	entries = append(entries, entry)
	mu.Unlock()

	log.Printf("[DLQ] task %s entered dead-letter queue, reason: %s", task.ID, reason)
	return nil
}

// List 返回所有死信记录（内存快照）
func List() []DLQEntry {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]DLQEntry, len(entries))
	copy(result, entries)
	return result
}

// Close 关闭 DLQ Writer
func Close() {
	if writer != nil {
		writer.Close()
	}
}
