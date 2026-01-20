package coordinator

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"fmt"
	"map-reduce/internal/rpc"
	"sync"
	"time"
)

type Coordinator struct {
	mu           sync.Mutex
	tasks        map[string]*rpc.Task   //taskId -> Task
	pendingQueue list.List              //taskId waiting to be assigned
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
		tasks:       make(map[string]*rpc.Task),
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
			if task.Status == rpc.TaskInProgress &&
				now.Sub(task.AssignedAt) > c.taskTimeout {
				task.Status = rpc.TaskPending
				c.pendingQueue.PushBack(task.ID)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Coordinator) handleWorkerFailure(workerID string) {
	worker := c.workers[workerID]
	if worker.CurrentTaskID != "" {
		task := c.tasks[worker.CurrentTaskID]
		if task.Status == rpc.TaskInProgress {
			task.Status = rpc.TaskPending
			c.pendingQueue.PushBack(task.ID)
		}
	}
	delete(c.workers, workerID)
}

func (c *Coordinator) AssignTask(workerID string) *rpc.Task {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pendingQueue.Len() == 0 {
		return nil
	}

	taskID := c.pendingQueue.Front().Value.(string)
	c.pendingQueue.Remove(c.pendingQueue.Front())

	task := c.tasks[taskID]
	task.Status = rpc.TaskInProgress
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
	if !ok || task.Status == rpc.TaskCompleted {
		return
	}

	task.Status = rpc.TaskCompleted

	// Aggregate results
	for word, count := range counts {
		c.results[word] += count
	}
}

func (c *Coordinator) HandleHeartbeat(payload []byte) (rpc.MessageType, []byte, error) {
	// Decode worker ID from payload
	var workerID string
	gob.NewDecoder(bytes.NewReader(payload)).Decode(&workerID)

	c.mu.Lock()
	defer c.mu.Unlock()

	if worker, ok := c.workers[workerID]; ok {
		worker.LastHeartbeat = time.Now()
	} else {
		// new worker, need to register
		c.workers[workerID] = &WorkerInfo{
			ID:            workerID,
			LastHeartbeat: time.Now(),
		}
	}
	return rpc.MsgAck, nil, nil

}

func (c *Coordinator) HandleRequestTask(payload []byte) (rpc.MessageType, []byte, error) {
	// Decode worker ID from payload
	var workerID string
	gob.NewDecoder(bytes.NewReader(payload)).Decode(&workerID)

	// Assign task
	task := c.AssignTask(workerID)
	if task == nil {
		return rpc.MsgTaskResponse, nil, nil
	}

	// Encode task to payload
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(task); err != nil {
		return rpc.MsgAck, nil, err
	}
	return rpc.MsgTaskResponse, buf.Bytes(), nil
}

func (c *Coordinator) HandleReportTask(payload []byte) (rpc.MessageType, []byte, error) {
	// decode word count result from payload
	var result rpc.WordCountResult
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(&result); err != nil {
		return rpc.MsgAck, nil, err
	}

	// handle result
	c.HandleResult(result.TaskID, result.Counts)
	return rpc.MsgAck, nil, nil
}

func (c *Coordinator) RegisterHandlers(server *rpc.Server) {
	server.Register(rpc.MsgHeartbeat, c.HandleHeartbeat)
	server.Register(rpc.MsgRequestTask, c.HandleRequestTask)
	server.Register(rpc.MsgReportTask, c.HandleReportTask)
}

func (c *Coordinator) AddTask(filename, content string) string {
	taskID := generateTaskID()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.tasks[taskID] = &rpc.Task{
		ID:       taskID,
		Filename: filename,
		Status:   rpc.TaskPending,
		Content:  content,
	}
	c.pendingQueue.PushBack(taskID)
	return taskID
}

func generateTaskID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}

func (c *Coordinator) GetResults() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()

	copy := make(map[string]int, len(c.results))
	for k, v := range c.results {
		copy[k] = v
	}
	return copy
}

func (c *Coordinator) AllTasksCompleted() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.tasks) == 0 {
		return false
	}

	for _, task := range c.tasks {
		if task.Status != rpc.TaskCompleted {
			return false
		}
	}
	return true
}
