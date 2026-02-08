package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditEvent struct {
	EventID     uuid.UUID
	EntityID    string
	EventType   EventType
	ActorID     string
	ActorType   ActorType
	Timestamp   time.Time
	Description string
	Metadata    map[string]string
	IPAddress   string
}
