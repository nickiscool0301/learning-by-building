package cassandra

import (
	"audit-system/internal/domain"
	"audit-system/pkg/timebucket"
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

const (
	insertByEntity = `INSERT INTO audit_by_entity (entity_id, time_bucket, event_id, actor_id, actor_type, description, metadata, ip_address, event_type, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	insertByActor = `INSERT INTO audit_by_actor (actor_id, time_bucket, event_id, entity_id, event_type, description, metadata, ip_address)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	insertByType = `INSERT INTO audit_by_type (event_type, time_bucket, event_id, entity_id, actor_id, description, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	selectByEntity = `SELECT event_id, entity_id, actor_id, actor_type, timestamp, description, metadata, ip_address, event_type
		FROM audit_by_entity WHERE entity_id = ? AND time_bucket = ? AND event_id >= ? AND event_id <= ?`
)

type Repository struct {
	session *gocql.Session
}

func NewRepository(session *gocql.Session) *Repository {
	return &Repository{
		session: session,
	}

}
func (r *Repository) Store(ctx context.Context, event domain.AuditEvent) error {
	timeBucket := timebucket.Compute(event.Timestamp, "daily")
	timeUUID := gocql.UUIDFromTime(event.Timestamp)

	// Why UnloggedBatch?
	// Cassandra logged batches are for atomicity across partitions and carry significant overhead.
	// Unlogged batches are used for performance optimization when atomicity is not required.
	// If need atomicity, use LoggedBatch.
	batch := r.session.NewBatch(gocql.UnloggedBatch)

	batch.Query(insertByEntity,
		event.EntityID, timeBucket, timeUUID, event.EntityID, string(event.EventType),
		event.ActorID, string(event.ActorType), event.Description, event.Metadata, event.IPAddress,
	)

	batch.Query(insertByActor,
		event.ActorID, timeBucket, timeUUID, event.EntityID, event.EventID, string(event.EventType),
		event.Description, event.Metadata, event.IPAddress,
	)

	batch.Query(insertByType,
		string(event.EventType), timeBucket, timeUUID, event.EventID, event.EntityID, event.ActorID, event.Description, event.Metadata,
	)

	if err := r.session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("failed to store audit event %s: %w", event.EventID, err)
	}
	return nil
}

func (r *Repository) FindByEntity(ctx context.Context, entityID string, start, end time.Time) ([]domain.AuditEvent, error) {
	buckets := timebucket.Range(start, end, "daily")
	var event domain.AuditEvent
	var events []domain.AuditEvent

	for _, bucket := range buckets {
		iter := r.session.Query(selectByEntity, entityID, bucket, gocql.MinTimeUUID(start), gocql.MaxTimeUUID(end)).WithContext(ctx).Iter()

		var actorType, eventType string
		for iter.Scan(&event.EventID, &event.EntityID, &event.ActorID, &actorType, &event.Timestamp,
			&event.Description, &event.Metadata, &event.IPAddress, &eventType) {
			event.ActorType = domain.ActorType(actorType)
			event.EventType = domain.EventType(eventType)
			events = append(events, event)
		}

		if err := iter.Close(); err != nil {
			return nil, fmt.Errorf("failed to query event for entity %s: %w", entityID, err)
		}
	}
	return events, nil
}
