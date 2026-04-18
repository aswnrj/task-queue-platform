package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aswnrj/task-queue-platform/internal/models"
	"github.com/aswnrj/task-queue-platform/internal/store"
	"github.com/aswnrj/task-queue-platform/internal/worker"
	"github.com/google/uuid"
)

type TaskHandler struct {
	store store.TaskStore
	pool  *worker.Pool
}

func NewTaskHandler(s store.TaskStore, p *worker.Pool) *TaskHandler {
	return &TaskHandler{store: s, pool: p}
}

type CreateTaskRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		http.Error(w, `{"error": "type is required"}`, http.StatusBadRequest)
		return
	}
	task := &models.Task{
		ID:        uuid.New().String(),
		Type:      req.Type,
		Payload:   req.Payload,
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := h.store.Create(task); err != nil {
		http.Error(w, `{"error": "failed to create task"}`, http.StatusInternalServerError)
		return
	}
	h.pool.Submit(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.store.Get(id)
	if err != nil {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List()
	if err != nil {
		http.Error(w, `{"error": "failed to list tasks"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
