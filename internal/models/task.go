package models

import "time"

type TaskStatus string

const (
	StatusPending TaskStatus = "pending"
	StatusRunning TaskStatus = "running"
	StatusDone    TaskStatus = "done"
	StatusFailed  TaskStatus = "failed"
)

type Task struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Payload   string     `json:"payload"`
	Status    TaskStatus `json:"status"`
	Result    string     `json:"result,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
