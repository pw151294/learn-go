// zoom/distributed_encoding_platform/store/store.go
package store

import (
	"sync"

	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/model"
)

var taskMap sync.Map

// Save 保存或更新任务
func Save(task *model.Task) {
	taskMap.Store(task.ID, task)
}

// Get 按 ID 查询任务，第二个返回值表示是否存在
func Get(id string) (*model.Task, bool) {
	val, ok := taskMap.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*model.Task), true
}

// List 返回所有任务的快照切片
func List() []*model.Task {
	var tasks []*model.Task
	taskMap.Range(func(_, val any) bool {
		tasks = append(tasks, val.(*model.Task))
		return true
	})
	return tasks
}
