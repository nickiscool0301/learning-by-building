package domain

import (
	"context"
	"time"
)

type EventRepository interface {
	Store(ctx context.Context, event AuditEvent) error
	FindByEntity(ctx context.Context, entityID string, start, end time.Time) ([]AuditEvent, error)
	FindByActor(ctx context.Context, actorID string, start, end time.Time) ([]AuditEvent, error)
	FindByType(ctx context.Context, eventType EventType, start, end time.Time) ([]AuditEvent, error)
}
