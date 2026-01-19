package coordinator

import (
	"sync"
	"time"
)

type Coordinator struct {
	mu           sync.Mutex
	tasks        map[string]*Task       //taskId -> Task
	pendingQueue []string               //taskId waiting to be assigned
	workers      map[string]*WorkerInfo //workerId -> WorkerInfo
	results      map[string]int         //aggregated word counts
	taskTimeout  time.Duration
}

type WorkerInfo struct {
	ID            string
	LastHeartbeat time.Time
	CurrentTaskID string
}

func New() *Coordinator {
	c := &Coordinator{
		tasks:       make(map[string]*Task),
		workers:     make(map[string]*WorkerInfo),
		results:     make(map[string]int),
		taskTimeout: 10 * time.Second,
	}
	go c.checkWorkerHealth()
	return c
}

func (c *Coordinator) checkWorkerHealth() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		for _, worker := range c.workers {
			// No heartbeat within timeout 15 seconds -> dead
			if now.Sub(worker.LastHeartbeat) > 15*time.Second {
				c.handleWorkerFailure(worker.ID)
			}
		}

		// Check for time-out tasks
		for _, task := range c.tasks {
			if task.Status == TaskInProgress &&
				now.Sub(task.AssignedAt) > c.taskTimeout {
				task.Status = TaskPending
				c.pendingQueue = append(c.pendingQueue, task.ID)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Coordinator) handleWorkerFailure(workerID string) {
	worker := c.workers[workerID]
	if worker.CurrentTaskID != "" {
		task := c.tasks[worker.CurrentTaskID]
		if task.Status == TaskInProgress {
			task.Status = TaskPending
			c.pendingQueue = append(c.pendingQueue, task.ID)
		}
	}
	delete(c.workers, workerID)
}

func (c *Coordinator) AssignTask(workerID string) *Task {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.pendingQueue) == 0 {
		return nil
	}

	taskID := c.pendingQueue[0]
	c.pendingQueue = c.pendingQueue[1:]

	task := c.tasks[taskID]
	task.Status = TaskInProgress
	task.AssignedTo = workerID
	task.AssignedAt = time.Now()

	if worker, ok := c.workers[workerID]; ok {
		worker.CurrentTaskID = taskID
	}

	return task
}

func (c *Coordinator) HandleResult(taskID string, counts map[string]int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	task, ok := c.tasks[taskID]
	if !ok || task.Status == TaskCompleted {
		return
	}

	task.Status = TaskCompleted

	// Aggregate results
	for word, count := range counts {
		c.results[word] += count
	}
}
