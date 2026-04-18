package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aswnrj/task-queue-platform/internal/models"
	"github.com/aswnrj/task-queue-platform/internal/store"
)

type Pool struct {
	taskChan chan *models.Task
	workers  int
	store    store.TaskStore
}

func NewPool(workers int, queueSize int, s store.TaskStore) *Pool {
	return &Pool{
		taskChan: make(chan *models.Task, queueSize),
		workers:  workers,
		store:    s,
	}
}

func (p *Pool) Submit(task *models.Task) {
	p.taskChan <- task
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		go p.work(ctx, i)
	}
	log.Printf("started %d workers", p.workers)
}

func (p *Pool) work(ctx context.Context, id int) {
	for {
		select {
		case task := <-p.taskChan:
			p.process(id, task)
		case <-ctx.Done():
			log.Printf("worker %d: shutting down", id)
			return
		}
	}
}

func (p *Pool) process(workerId int, task *models.Task) {
	log.Printf("worker %d: processing task %s (type=%s)", workerId, task.ID, task.Type)

	p.store.UpdateStatus(task.ID, models.StatusRunning, "")

	time.Sleep(2 * time.Second)

	result := fmt.Sprintf("completed by worker %d", workerId)
	p.store.UpdateStatus(task.ID, models.StatusDone, result)

	log.Printf("worker %d: finished task %s", workerId, task.ID)
}
