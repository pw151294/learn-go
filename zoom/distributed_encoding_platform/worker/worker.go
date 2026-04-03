// zoom/distributed_encoding_platform/worker/worker.go
package worker

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/config"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/dlq"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/model"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/producer"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/store"
)

// semaphore 用带缓冲 channel 限制并发转码数
var semaphore = make(chan struct{}, config.MaxWorkers)

// Submit 将任务提交到 Worker Pool 异步执行（非阻塞入口）
func Submit(ctx context.Context, task *model.Task) {
	go func() {
		// 获取 semaphore slot，超过 MaxWorkers 时阻塞
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		process(ctx, task)
	}()
}

// process 执行单个任务的转码逻辑
func process(ctx context.Context, task *model.Task) {
	log.Printf("[Worker] start transcoding task %s (priority=%s, retries=%d)",
		task.ID, task.Priority, task.Retries)

	// 状态: pending → transcoding
	task.TransitionTo(model.StatusTranscoding)
	store.Save(task)

	// 模拟转码耗时 1~5 秒
	duration := time.Duration(1+rand.Intn(5)) * time.Second
	select {
	case <-time.After(duration):
	case <-ctx.Done():
		log.Printf("[Worker] task %s cancelled", task.ID)
		return
	}

	// 模拟失败（30% 概率）
	if rand.Float64() < config.FailureRate {
		handleFailure(ctx, task, "simulated transcoding error")
		return
	}

	// 转码成功
	task.TransitionTo(model.StatusSuccess)
	store.Save(task)
	log.Printf("[Worker] task %s succeeded", task.ID)
}

// handleFailure 处理转码失败：重试或进入 DLQ
func handleFailure(ctx context.Context, task *model.Task, reason string) {
	task.Retries++
	task.Error = reason
	log.Printf("[Worker] task %s failed (retries=%d): %s", task.ID, task.Retries, reason)

	if task.Retries >= config.MaxRetries {
		// 超过最大重试次数，进入死信队列
		task.TransitionTo(model.StatusFailed)
		store.Save(task)
		dlq.Push(task, fmt.Sprintf("exceeded max retries(%d): %s", config.MaxRetries, reason))
		return
	}

	// 重试：重新入队（回到 pending，写回同优先级 Topic）
	task.TransitionTo(model.StatusPending)
	store.Save(task)
	if err := producer.Send(ctx, task); err != nil {
		log.Printf("[Worker] re-queue task %s error: %v", task.ID, err)
	} else {
		log.Printf("[Worker] task %s re-queued for retry %d", task.ID, task.Retries)
	}
}
