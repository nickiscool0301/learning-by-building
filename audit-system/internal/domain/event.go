package domain

import (
	"errors"
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

func NewAuditEvent(
	entityID string,
	eventType EventType,
	actorID string,
	actorType ActorType,
	description string,
	metadata map[string]string,
	ipAddress string,
) AuditEvent {
	return AuditEvent{
		EventID:     uuid.New(),
		EntityID:    entityID,
		EventType:   eventType,
		ActorID:     actorID,
		ActorType:   actorType,
		Timestamp:   time.Now().UTC(),
		Description: description,
		Metadata:    metadata,
		IPAddress:   ipAddress,
	}
}

func (e AuditEvent) Validate() error {
	var errs []error

	if e.EntityID == "" {
		errs = append(errs, errors.New("entity_id is required"))
	}
	if e.ActorID == "" {
		errs = append(errs, errors.New("actor_id is required"))
	}
	if err := e.EventType.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := e.ActorType.Validate(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
