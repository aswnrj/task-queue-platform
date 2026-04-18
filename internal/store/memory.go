package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/aswnrj/task-queue-platform/internal/models"
)

type MemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]*models.Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]*models.Task),
	}
}

func (s *MemoryStore) Create(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
	return nil
}

func (s *MemoryStore) Get(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return task, nil
}

func (s *MemoryStore) List() ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *MemoryStore) UpdateStatus(id string, status models.TaskStatus, result string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found: %s", id)
	}
	task.Status = status
	task.Result = result
	task.UpdatedAt = time.Now()
	return nil
}
