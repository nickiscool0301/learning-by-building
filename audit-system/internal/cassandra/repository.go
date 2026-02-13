package cassandra

import (
	"audit-system/internal/domain"
	"audit-system/pkg/timebucket"
	"context"
	"fmt"

	"github.com/gocql/gocql"
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
		event.ActorID, string(event.ActorType), event.Description, event.Metadata, event.IPAddress
	)

	batch.Query(insertByActor,
		event.ActorID, timeBucket, timeUUID, event.EntityID, event.EntityID, string(event.EventType),
		event.Description, event.Metadata, event.IPAddress
	)

	batch.Query(insertByType,
		string(event.EventType), timeBucket, timeUUID, event.EventID, event.EntityID, event.ActorID, event.Description, event.Metadata
	)

	if err := r.session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("failed to store audit event %s: %w", event.EventID, err)
	}
	return nil
}
