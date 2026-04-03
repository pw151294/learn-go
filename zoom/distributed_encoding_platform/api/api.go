// zoom/distributed_encoding_platform/api/api.go
package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/dlq"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/model"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/producer"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/store"
)

// Register 注册所有 HTTP 路由
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/task", handleTask)
	mux.HandleFunc("/tasks", handleTasks)
	mux.HandleFunc("/dlq", handleDLQ)
}

// handleTask 处理 POST /task（提交）和 GET /task/{id}（查询）
func handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		submitTask(w, r)
	case http.MethodGet:
		// 支持 /task?id=xxx 或 /task/xxx
		id := r.URL.Query().Get("id")
		if id == "" {
			id = strings.TrimPrefix(r.URL.Path, "/task/")
		}
		getTask(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// submitTask 接收任务请求，写入 Kafka
func submitTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		File     string `json:"file"`
		Format   string `json:"format"`
		Priority string `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Priority == "" {
		req.Priority = "normal"
	}
	// 提前验证 priority 合法性，避免写入脏数据到 store
	switch req.Priority {
	case "high", "normal", "low":
	default:
		http.Error(w, "invalid priority: must be high, normal, or low", http.StatusBadRequest)
		return
	}

	task := &model.Task{
		ID:        uuid.New().String(),
		File:      req.File,
		Format:    req.Format,
		Priority:  req.Priority,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	store.Save(task)

	if err := producer.Send(r.Context(), task); err != nil {
		http.Error(w, "failed to enqueue task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

// getTask 按 ID 查询任务状态
func getTask(w http.ResponseWriter, _ *http.Request, id string) {
	task, ok := store.Get(id)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// handleTasks 返回所有任务列表
func handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, store.List())
}

// handleDLQ 返回死信队列内容
func handleDLQ(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, dlq.List())
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
