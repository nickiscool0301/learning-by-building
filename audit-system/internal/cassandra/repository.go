package cassandra

import (
	"audit-system/internal/domain"
	"context"
)

func (r *Repository) Store(ctx context.Context, event domain.AuditEvent) error {

}
