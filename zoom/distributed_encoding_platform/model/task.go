// zoom/distributed_encoding_platform/model/task.go
package model

import "time"

// TaskStatus 表示转码任务的状态
type TaskStatus string

const (
	StatusPending     TaskStatus = "pending"     // 已提交，等待消费
	StatusTranscoding TaskStatus = "transcoding" // 正在转码
	StatusSuccess     TaskStatus = "success"     // 转码成功
	StatusFailed      TaskStatus = "failed"      // 超过最大重试，已入死信队列
)

// Task 是转码任务的完整描述
type Task struct {
	ID        string     `json:"task_id"`
	File      string     `json:"file"`
	Format    string     `json:"format"`
	Priority  string     `json:"priority"` // "high" | "normal" | "low"
	Status    TaskStatus `json:"status"`
	Retries   int        `json:"retries"`
	Error     string     `json:"error,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TransitionTo 状态机：执行合法的状态转换，返回是否成功
func (t *Task) TransitionTo(next TaskStatus) bool {
	valid := map[TaskStatus][]TaskStatus{
		StatusPending:     {StatusTranscoding},
		StatusTranscoding: {StatusSuccess, StatusFailed, StatusPending}, // pending = 重试重新入队
	}
	for _, allowed := range valid[t.Status] {
		if allowed == next {
			t.Status = next
			t.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}
