// zoom/distributed_encoding_platform/producer/producer.go
package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/config"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/model"
)

// writers 按 topic 缓存 Kafka Writer
var writers = map[string]*kafka.Writer{}

// Init 初始化三个优先级 Topic 的 Writer
func Init() {
	for _, topic := range config.PriorityTopics {
		writers[topic] = &kafka.Writer{
			Addr:     kafka.TCP(config.KafkaBroker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
	}
}

// Send 将任务序列化后发送到对应优先级的 Kafka Topic
func Send(ctx context.Context, task *model.Task) error {
	topic, err := priorityToTopic(task.Priority)
	if err != nil {
		return err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("producer marshal error: %w", err)
	}

	return writers[topic].WriteMessages(ctx, kafka.Message{
		Key:   []byte(task.ID),
		Value: data,
	})
}

// priorityToTopic 将优先级字符串映射到 Topic 名
func priorityToTopic(priority string) (string, error) {
	switch priority {
	case "high":
		return config.TopicHigh, nil
	case "normal":
		return config.TopicNormal, nil
	case "low":
		return config.TopicLow, nil
	default:
		return "", fmt.Errorf("unknown priority: %s", priority)
	}
}

// Close 关闭所有 Writer
func Close() {
	for _, w := range writers {
		w.Close()
	}
}
