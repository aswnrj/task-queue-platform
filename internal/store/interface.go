package store

import "github.com/aswnrj/task-queue-platform/internal/models"

type TaskStore interface {
	Create(task *models.Task) error
	Get(id string) (*models.Task, error)
	List() ([]*models.Task, error)
	UpdateStatus(id string, status models.TaskStatus, result string) error
}
