// Package service 实现上传业务的核心逻辑：秒传、断点续传、分片校验、状态流转。
package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/model"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/mq"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/repository"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/storage"
)

type UploadService struct {
	taskRepo  *repository.UploadTaskRepo
	chunkRepo *repository.ChunkRepo
	s3        *storage.S3Client
}

func NewUploadService(taskRepo *repository.UploadTaskRepo, chunkRepo *repository.ChunkRepo, s3 *storage.S3Client) *UploadService {
	return &UploadService{taskRepo: taskRepo, chunkRepo: chunkRepo, s3: s3}
}

// InitUploadResult 初始化上传的返回结果。
type InitUploadResult struct {
	TaskID         string `json:"task_id"`
	UploadID       string `json:"upload_id"`
	AlreadyExists  bool   `json:"already_exists"` // true 表示秒传
	FileURL        string `json:"file_url,omitempty"`
	ExistingChunks []int  `json:"existing_chunks"` // 已上传的分片索引（断点续传）
}

// InitUpload 初始化上传任务，处理三种场景：
// 1. 秒传：file_md5 已存在且状态为 completed
// 2. 断点续传：file_md5 已存在但未完成，返回已上传分片列表
// 3. 新建任务：获取分布式锁，创建 S3 Multipart Upload
func (s *UploadService) InitUpload(fileName string, fileSize int64, fileMD5 string, totalChunks int) (*InitUploadResult, error) {
	// 1. 查找是否已有相同 MD5 的任务
	task, err := s.taskRepo.FindByMD5(fileMD5)
	if err != nil {
		return nil, fmt.Errorf("FindByMD5: %w", err)
	}

	if task != nil {
		// 秒传：文件已上传完成
		if task.Status == model.TaskStatusCompleted {
			log.Printf("[service] fast upload: task_id=%s file_md5=%s", task.TaskID, fileMD5)
			return &InitUploadResult{
				TaskID:        task.TaskID,
				UploadID:      task.UploadID,
				AlreadyExists: true,
				FileURL:       task.FileURL,
			}, nil
		}

		// 断点续传：任务存在但未完成，返回已上传的分片
		uploadedChunks, err := s.chunkRepo.FindByTaskID(task.TaskID)
		if err != nil {
			return nil, fmt.Errorf("FindByTaskID: %w", err)
		}
		existingIndexes := make([]int, 0, len(uploadedChunks))
		for _, c := range uploadedChunks {
			existingIndexes = append(existingIndexes, c.ChunkIndex)
		}
		log.Printf("[service] resume upload: task_id=%s uploaded=%d/%d", task.TaskID, len(existingIndexes), totalChunks)
		return &InitUploadResult{
			TaskID:         task.TaskID,
			UploadID:       task.UploadID,
			ExistingChunks: existingIndexes,
		}, nil
	}

	// 2. 新建任务：获取分布式锁防止并发重复创建
	mutex, err := storage.AcquireLock(fileMD5, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("并发上传冲突: %w", err)
	}
	defer func() {
		if _, err := mutex.Unlock(); err != nil {
			log.Printf("[service] unlock error: %v", err)
		}
	}()

	// 双重检查：获取锁后再查一次
	task, err = s.taskRepo.FindByMD5(fileMD5)
	if err != nil {
		return nil, err
	}
	if task != nil {
		// 另一个请求已抢先创建
		return &InitUploadResult{TaskID: task.TaskID, UploadID: task.UploadID}, nil
	}

	// 3. 创建 S3 Multipart Upload
	taskID := uuid.New().String()
	objectKey := fmt.Sprintf("videos/%s/%s", taskID, fileName)
	uploadID, err := s.s3.CreateMultipartUpload(objectKey)
	if err != nil {
		return nil, fmt.Errorf("CreateMultipartUpload: %w", err)
	}

	// 4. 持久化任务到数据库
	newTask := &model.UploadTask{
		TaskID:      taskID,
		FileName:    fileName,
		FileSize:    fileSize,
		FileMD5:     fileMD5,
		UploadID:    uploadID,
		Bucket:      s.s3.Bucket(),
		ObjectKey:   objectKey,
		TotalChunks: totalChunks,
		Status:      model.TaskStatusUploading,
	}
	if err = s.taskRepo.Create(newTask); err != nil {
		// 回滚 S3 multipart
		_ = s.s3.AbortMultipartUpload(objectKey, uploadID)
		return nil, fmt.Errorf("Create task: %w", err)
	}

	log.Printf("[service] new upload task: task_id=%s total_chunks=%d", taskID, totalChunks)
	return &InitUploadResult{
		TaskID:         taskID,
		UploadID:       uploadID,
		ExistingChunks: []int{},
	}, nil
}

// UploadChunkResult 上传分片的返回结果。
type UploadChunkResult struct {
	ChunkIndex int    `json:"chunk_index"`
	ETag       string `json:"etag"`
	Status     string `json:"status"`
}

// UploadChunk 处理单个分片上传：
// 1. 幂等检查（分片已上传则直接返回）
// 2. 校验分片 MD5 保证完整性
// 3. 上传到 S3 并获取 ETag
// 4. 持久化分片记录
func (s *UploadService) UploadChunk(taskID string, chunkIndex int, chunkMD5 string, data io.Reader, size int64) (*UploadChunkResult, error) {
	// 获取任务信息
	task, err := s.taskRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("FindByTaskID: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	if task.Status == model.TaskStatusCompleted {
		return nil, fmt.Errorf("task already completed: %s", taskID)
	}

	// 幂等检查：分片已存在则直接返回（断点续传时前端可能重传）
	existing, err := s.chunkRepo.FindChunk(taskID, chunkIndex)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Status == model.ChunkStatusUploaded {
		log.Printf("[service] chunk already uploaded: task_id=%s chunk=%d", taskID, chunkIndex)
		return &UploadChunkResult{ChunkIndex: chunkIndex, ETag: existing.ETag, Status: existing.Status}, nil
	}

	// 读取分片数据并校验 MD5
	chunkData, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("read chunk data: %w", err)
	}
	h := md5.New()
	h.Write(chunkData)
	actualMD5 := hex.EncodeToString(h.Sum(nil))
	if actualMD5 != chunkMD5 {
		return nil, fmt.Errorf("MD5 校验失败: expected=%s actual=%s", chunkMD5, actualMD5)
	}

	// 上传到 S3
	etag, err := s.s3.UploadPart(task.ObjectKey, task.UploadID, chunkIndex,
		newBytesReader(chunkData), int64(len(chunkData)))
	if err != nil {
		return nil, fmt.Errorf("UploadPart: %w", err)
	}

	// 持久化分片记录（Upsert 保证幂等）
	chunk := &model.UploadChunk{
		TaskID:     taskID,
		ChunkIndex: chunkIndex,
		ChunkSize:  int64(len(chunkData)),
		ChunkMD5:   chunkMD5,
		ETag:       etag,
		Status:     model.ChunkStatusUploaded,
	}
	if err = s.chunkRepo.Upsert(chunk); err != nil {
		return nil, fmt.Errorf("Upsert chunk: %w", err)
	}

	log.Printf("[service] chunk uploaded: task_id=%s chunk=%d etag=%s", taskID, chunkIndex, etag)
	return &UploadChunkResult{ChunkIndex: chunkIndex, ETag: etag, Status: model.ChunkStatusUploaded}, nil
}

// CompleteUpload 完成上传：合并所有分片 + 更新状态 + 发送 Kafka 转码通知。
func (s *UploadService) CompleteUpload(taskID string) (string, error) {
	task, err := s.taskRepo.FindByTaskID(taskID)
	if err != nil {
		return "", fmt.Errorf("FindByTaskID: %w", err)
	}
	if task == nil {
		return "", fmt.Errorf("task not found: %s", taskID)
	}
	// 幂等：已完成的任务直接返回
	if task.Status == model.TaskStatusCompleted {
		return task.FileURL, nil
	}

	// 校验所有分片已上传
	uploadedCount, err := s.chunkRepo.CountUploaded(taskID)
	if err != nil {
		return "", err
	}
	if int(uploadedCount) != task.TotalChunks {
		return "", fmt.Errorf("分片不完整: uploaded=%d total=%d", uploadedCount, task.TotalChunks)
	}

	// 构造 S3 CompletePart 列表
	chunks, err := s.chunkRepo.FindByTaskID(taskID)
	if err != nil {
		return "", err
	}
	parts := make([]minio.CompletePart, len(chunks))
	for i, c := range chunks {
		parts[i] = minio.CompletePart{
			PartNumber: c.ChunkIndex + 1, // S3 partNumber 从 1 开始
			ETag:       c.ETag,
		}
	}

	// 调用 S3 合并分片
	fileURL, err := s.s3.CompleteMultipartUpload(task.ObjectKey, task.UploadID, parts)
	if err != nil {
		_ = s.taskRepo.UpdateStatus(taskID, model.TaskStatusFailed)
		return "", fmt.Errorf("CompleteMultipartUpload: %w", err)
	}

	// 更新数据库状态
	if err = s.taskRepo.UpdateComplete(taskID, fileURL); err != nil {
		return "", fmt.Errorf("UpdateComplete: %w", err)
	}

	// 发送 Kafka 转码通知（失败不阻断主流程，等待补偿重试）
	msg := mq.TranscodeMessage{
		TaskID:    taskID,
		FileURL:   fileURL,
		FileName:  task.FileName,
		FileSize:  task.FileSize,
		CreatedAt: time.Now(),
	}
	if err = mq.SendTranscodeMessage(msg); err != nil {
		log.Printf("[service] WARNING: kafka send failed task_id=%s err=%v", taskID, err)
	}

	log.Printf("[service] upload completed: task_id=%s file_url=%s", taskID, fileURL)
	return fileURL, nil
}

// GetProgress 查询上传任务的当前进度。
func (s *UploadService) GetProgress(taskID string) (map[string]interface{}, error) {
	task, err := s.taskRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	uploadedCount, err := s.chunkRepo.CountUploaded(taskID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"task_id":         task.TaskID,
		"file_name":       task.FileName,
		"status":          task.Status,
		"total_chunks":    task.TotalChunks,
		"uploaded_chunks": uploadedCount,
		"file_url":        task.FileURL,
	}, nil
}

// bytesReader 将 []byte 包装为满足 io.Reader 和 io.ReaderAt 的读取器。
type bytesReader struct {
	data   []byte
	offset int
}

func newBytesReader(data []byte) *bytesReader { return &bytesReader{data: data} }

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.offset:])
	r.offset += n
	return n, nil
}

func (r *bytesReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= int64(len(r.data)) {
		return 0, io.EOF
	}
	n = copy(p, r.data[off:])
	return n, nil
}
