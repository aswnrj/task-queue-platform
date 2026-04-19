package store

import (
	"context"
	"fmt"
	"time"

	"github.com/aswnrj/task-queue-platform/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) Create(task *models.Task) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO tasks (id, type, payload, status, result, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		task.ID, task.Type, task.Payload, task.Status, task.Result, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func (s *PostgresStore) Get(id string) (*models.Task, error) {
	row := s.pool.QueryRow(context.Background(),
		`SELECT id, type, payload, status, result, created_at, updated_at 
		FROM tasks WHERE id = $1`, id,
	)
	var task models.Task
	if err := row.Scan(&task.ID, &task.Type, &task.Payload, &task.Status, &task.Result, &task.CreatedAt, &task.UpdatedAt); err != nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return &task, nil
}

func (s *PostgresStore) List() ([]*models.Task, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT id, type, payload, status, result, created_at, updated_at
		FROM tasks ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Type, &t.Payload, &t.Status, &t.Result, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func (s *PostgresStore) UpdateStatus(id string, status models.TaskStatus, result string) error {
	_, err := s.pool.Exec(context.Background(),
		`UPDATE tasks SET status = $1, result = $2, updated_at = $3 WHERE id = $4`,
		status, result, time.Now().UTC(), id,
	)
	return err
}
