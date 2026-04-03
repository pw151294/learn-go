// zoom/distributed_encoding_platform/consumer/consumer.go
package consumer

import (
	"context"
	"log"
	"time"

	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/scheduler"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/worker"
)

// Start 启动消费者主循环，直到 ctx 取消
// 持续从 Scheduler 拉取任务，提交给 Worker Pool 执行
func Start(ctx context.Context) {
	log.Println("[Consumer] starting consumer group loop...")
	for {
		select {
		case <-ctx.Done():
			log.Println("[Consumer] stopping consumer loop")
			return
		default:
		}

		task := scheduler.FetchNext(ctx)
		if task == nil {
			// 所有 Topic 均无消息，短暂休眠避免空转 CPU
			time.Sleep(50 * time.Millisecond)
			continue
		}

		// 非阻塞提交给 Worker Pool（内部有 semaphore 控制并发）
		worker.Submit(ctx, task)
	}
}
