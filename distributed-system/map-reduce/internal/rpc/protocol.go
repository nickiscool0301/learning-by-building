package rpc

type MessageType uint8

const (
	MsgRegister MessageType = iota + 1
	MsgHeartbeat
	MsgRequestTask
	MsgTaskResponse
	MsgReportTask
	MsgAck
)

type Message struct {
	Type    MessageType
	Payload []byte
}

type Task struct {
	ID       string
	Filename string
	Content  string
	Status   TaskStatus
}

type TaskStatus uint8

const (
	TaskPending TaskStatus = iota
	TaskInProgress
	TaskCompleted
	TaskFailed
)

type WordCountResult struct {
	TaskID string
	Counts map[string]int
}
