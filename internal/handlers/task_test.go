package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aswnrj/task-queue-platform/internal/models"
	"github.com/aswnrj/task-queue-platform/internal/store"
	"github.com/aswnrj/task-queue-platform/internal/worker"
)

func setupHandler() *TaskHandler {
	s := store.NewMemoryStore()
	p := worker.NewPool(1, 10, s)
	return NewTaskHandler(s, p)
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"valid task", `{"type":"email","payload":"hello"}`, http.StatusCreated},
		{"missing type", `{"payload":"hello"}`, http.StatusBadRequest},
		{"empty body", `{}`, http.StatusBadRequest},
		{"invalid json", `not json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/task", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := setupHandler()

			h.Create(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, w.Code)
			}
			if w.Code == http.StatusCreated {
				var task models.Task
				json.NewDecoder(w.Body).Decode(&task)
				if task.ID == "" {
					t.Errorf("expected task id to be set")
				}
				if task.Status != models.StatusPending {
					t.Errorf("expected pending, got %s", task.Status)
				}
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"type":"email","payload":"test"}`))
	w := httptest.NewRecorder()

	h.Create(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("error in event creation")
	}
	var created models.Task
	json.NewDecoder(w.Body).Decode(&created)

	getReq := httptest.NewRequest("GET", "/tasks/"+created.ID, nil)
	getReq.SetPathValue("id", created.ID)
	getW := httptest.NewRecorder()
	h.Get(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", getW.Code)
	}

	var fetched models.Task
	json.NewDecoder(getW.Body).Decode(&fetched)
	if fetched.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, fetched.ID)
	}
}

func TestTaskNotFound(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest("GET", "/tasks/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	h.Get(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestListTasks(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	h.List(w, req)
	var tasks []models.Task
	json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}
