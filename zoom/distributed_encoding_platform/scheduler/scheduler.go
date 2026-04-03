// zoom/distributed_encoding_platform/scheduler/scheduler.go
package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/config"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/model"
)

// readers 按 topic 缓存 Kafka Reader
var readers map[string]*kafka.Reader

// Init 初始化三个优先级 Topic 的 Reader
func Init() {
	readers = make(map[string]*kafka.Reader)
	for _, topic := range config.PriorityTopics {
		readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{config.KafkaBroker},
			Topic:          topic,
			GroupID:        config.GroupID,
			MinBytes:       1,
			MaxBytes:       10e6,
			CommitInterval: time.Second,
		})
	}
}

// FetchNext 按 high→normal→low 顺序尝试拉取一条消息。
// 使用短超时非阻塞尝试，所有 Topic 均无消息时返回 nil。
func FetchNext(ctx context.Context) *model.Task {
	for _, topic := range config.PriorityTopics {
		task := tryFetch(ctx, topic)
		if task != nil {
			return task
		}
	}
	return nil
}

// tryFetch 对单个 Topic 做带超时的非阻塞拉取，超时返回 nil
func tryFetch(ctx context.Context, topic string) *model.Task {
	reader := readers[topic]
	// 100ms 超时：有消息立即返回，无消息快速放弃
	fetchCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	msg, err := reader.FetchMessage(fetchCtx)
	if err != nil {
		// 超时或 context 取消均属正常（无消息）
		return nil
	}

	var task model.Task
	if err := json.Unmarshal(msg.Value, &task); err != nil {
		log.Printf("[Scheduler] unmarshal error in topic %s: %v", topic, err)
		// 提交 offset 跳过损坏消息
		reader.CommitMessages(ctx, msg)
		return nil
	}

	// 提交 offset，确认消费
	if err := reader.CommitMessages(ctx, msg); err != nil {
		log.Printf("[Scheduler] commit error: %v", err)
	}

	log.Printf("[Scheduler] fetched task %s from topic %s", task.ID, topic)
	return &task
}

// Close 关闭所有 Reader
func Close() {
	for _, r := range readers {
		r.Close()
	}
}
