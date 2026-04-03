// zoom/distributed_encoding_platform/config/config.go
package config

const (
	KafkaBroker = "localhost:9092"

	TopicHigh   = "task.high"
	TopicNormal = "task.normal"
	TopicLow    = "task.low"
	TopicDLQ    = "task.dlq"

	GroupID     = "transcoding-group"
	MaxWorkers  = 3
	MaxRetries  = 3
	FailureRate = 0.3 // 30% 模拟失败率
)

// PriorityTopics 按优先级顺序排列，Scheduler 依次尝试
var PriorityTopics = []string{TopicHigh, TopicNormal, TopicLow}
