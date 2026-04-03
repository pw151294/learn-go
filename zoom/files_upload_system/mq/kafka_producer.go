// Package mq 封装 Kafka 生产者，用于发送转码通知消息。
package mq

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/config"
)

// TranscodeMessage 是发送到 Kafka 的转码通知消息结构。
type TranscodeMessage struct {
	TaskID    string    `json:"task_id"`
	FileURL   string    `json:"file_url"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}

var producer sarama.SyncProducer

// InitKafkaProducer 初始化 Kafka 同步生产者。
// 使用同步模式确保消息可靠投递（等待所有副本确认）。
func InitKafkaProducer() error {
	cfg := config.Global
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	saramaCfg.Producer.Retry.Max = 3

	var err error
	producer, err = sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}
	log.Printf("[kafka] producer connected: brokers=%s topic=%s", cfg.KafkaBrokers, cfg.KafkaTopic)
	return nil
}

// SendTranscodeMessage 发送转码通知消息到 Kafka。
// 使用 task_id 作为消息 key，确保同一任务的消息路由到同一分区（有序）。
func SendTranscodeMessage(msg TranscodeMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	pm := &sarama.ProducerMessage{
		Topic: config.Global.KafkaTopic,
		Key:   sarama.StringEncoder(msg.TaskID), // 同一任务路由到同一分区
		Value: sarama.ByteEncoder(body),
	}

	partition, offset, err := producer.SendMessage(pm)
	if err != nil {
		return fmt.Errorf("SendMessage: %w", err)
	}
	log.Printf("[kafka] sent transcode msg: task_id=%s partition=%d offset=%d", msg.TaskID, partition, offset)
	return nil
}

// Close 关闭 Kafka 生产者连接。
func Close() {
	if producer != nil {
		_ = producer.Close()
	}
}
