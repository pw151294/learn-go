// Package repository 提供对 upload_tasks 表的 CRUD 操作。
package repository

import (
	"errors"

	"gorm.io/gorm"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/model"
)

// UploadTaskRepo 封装 upload_tasks 表的数据库操作。
type UploadTaskRepo struct {
	db *gorm.DB
}

func NewUploadTaskRepo(db *gorm.DB) *UploadTaskRepo {
	return &UploadTaskRepo{db: db}
}

// FindByMD5 通过文件 MD5 查找任务，返回 nil 表示不存在。
func (r *UploadTaskRepo) FindByMD5(fileMD5 string) (*model.UploadTask, error) {
	var task model.UploadTask
	err := r.db.Where("file_md5 = ?", fileMD5).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}

// FindByTaskID 通过 task_id 查找任务，返回 nil 表示不存在。
func (r *UploadTaskRepo) FindByTaskID(taskID string) (*model.UploadTask, error) {
	var task model.UploadTask
	err := r.db.Where("task_id = ?", taskID).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}

// Create 写入新的上传任务记录。
func (r *UploadTaskRepo) Create(task *model.UploadTask) error {
	return r.db.Create(task).Error
}

// UpdateStatus 仅更新任务状态字段。
func (r *UploadTaskRepo) UpdateStatus(taskID, status string) error {
	return r.db.Model(&model.UploadTask{}).
		Where("task_id = ?", taskID).
		Update("status", status).Error
}

// UpdateComplete 将任务标记为已完成，同时记录最终文件 URL。
func (r *UploadTaskRepo) UpdateComplete(taskID, fileURL string) error {
	return r.db.Model(&model.UploadTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"status":   model.TaskStatusCompleted,
			"file_url": fileURL,
		}).Error
}
