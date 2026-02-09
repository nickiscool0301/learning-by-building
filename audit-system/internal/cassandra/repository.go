package cassandra

import (
	"audit-system/internal/domain"
	"context"

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

}
