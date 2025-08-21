package outbox

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	SaveOutboxEntry(ctx context.Context, tx *sql.Tx, entry *OutboxEntry) error
	GetUnsentEntries(ctx context.Context, limit int) ([]OutboxEntry, error)
	MarkAsSent(ctx context.Context, ids []int64) error
}

func closeRows(ctx context.Context, rows *sql.Rows, logger *logrus.Logger) {
	err := rows.Close()
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to close rows: %v", err)
	}
}
