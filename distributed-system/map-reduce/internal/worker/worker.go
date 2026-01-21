package worker

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"map-reduce/internal/rpc"
	"strings"
	"time"
	"unicode"
)

type Worker struct {
	id                 string
	coordinatorAddress string
	client             *rpc.Client
}

func New(id, coordinatorAddress string) (*Worker, error) {
	return &Worker{
		id:                 id,
		coordinatorAddress: coordinatorAddress,
		client:             rpc.NewClient(coordinatorAddress),
	}, nil
}

func (w *Worker) Run() {
	// Start heartbeat with goroutine
	go w.sendHeartbeats()

	// Main loop
	for {
		task := w.requestTask()
		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		counts := w.countWords(task.Content)
		w.reportResult(task.ID, counts)
	}
}

func (w *Worker) countWords(content string) map[string]int {
	counts := make(map[string]int)

	words := strings.FieldsFunc(content, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	for _, word := range words {
		word := strings.ToLower(word)
		if word != "" {
			counts[word]++
		}
	}
	return counts
}

func (w *Worker) sendHeartbeats() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		fmt.Printf("[%s] Sending heartbeat\n", w.id)
		w.client.Send(rpc.Message{
			Type:    rpc.MsgHeartbeat,
			Payload: encodeWorkerID(w.id),
		})
	}
}

func encodeWorkerID(id string) []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(id)
	return buf.Bytes()
}

func (w *Worker) requestTask() *rpc.Task {
	fmt.Printf("[%s] Requesting task...\n", w.id)
	resp, err := w.client.Send(rpc.Message{
		Type:    rpc.MsgRequestTask,
		Payload: encodeWorkerID(w.id),
	})
	if err != nil || resp.Type != rpc.MsgTaskResponse || len(resp.Payload) == 0 {
		fmt.Printf("[%s] No task available\n", w.id)
		return nil
	}

	var task rpc.Task
	if err := gob.NewDecoder(bytes.NewReader(resp.Payload)).Decode(&task); err != nil {
		fmt.Printf("[%s] Failed to decode task: %v\n", w.id, err)
		return nil
	}
	fmt.Printf("[%s] Got task: %s\n", w.id, task.ID)
	return &task
}

func (w *Worker) reportResult(taskID string, counts map[string]int) {
	fmt.Printf("[%s] Reporting result for %s: %v\n", w.id, taskID, counts)
	result := rpc.WordCountResult{
		TaskID: taskID,
		Counts: counts,
	}

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(result)

	w.client.Send(rpc.Message{
		Type:    rpc.MsgReportTask,
		Payload: buf.Bytes(),
	})
}
