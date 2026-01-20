# Distributed Word Counter

Build a distributed system where workers count words in parallel, coordinated by a central node.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Coordinator                         │
│  - Holds task queue (pending, in-progress, completed)   │
│  - Tracks worker health via heartbeats                  │
│  - Aggregates word counts from all workers              │
│  - Reassigns tasks if worker dies                       │
└─────────────────────────────────────────────────────────┘
           ▲              ▲              ▲
           │ TCP/gob      │              │
           ▼              ▼              ▼
      ┌─────────┐   ┌─────────┐   ┌─────────┐
      │ Worker1 │   │ Worker2 │   │ Worker3 │
      └─────────┘   └─────────┘   └─────────┘

Each worker:
  1. Sends heartbeat every 5s
  2. Requests task from coordinator
  3. Counts words in the text
  4. Reports result back
```

## Message Flow

1. **Worker starts** → sends `MsgRegister` to coordinator
2. **Every 5 seconds** → worker sends `MsgHeartbeat` (coordinator marks worker alive)
3. **Worker idle** → sends `MsgRequestTask` → gets back a `Task` with text content
4. **Worker finishes** → sends `MsgReportTask` with word counts
5. **Worker crashes** → coordinator notices missing heartbeat after 15s → requeues task

## Failure Handling

- **Worker dies mid-task**: Coordinator detects no heartbeat for 15s, puts task back in queue
- **Task takes too long**: After 10s, coordinator assumes failure and requeues
- **Network hiccup**: Client retries 3 times with exponential backoff (100ms, 200ms, 400ms)

## Success Criteria
System should handle 1 worker crashing mid-task without losing data.