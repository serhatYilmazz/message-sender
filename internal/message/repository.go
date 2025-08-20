package message

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	// FindAllMessages Limit would be specified
	FindAllMessages(ctx context.Context) ([]Message, error)
	MarkAsSent(ctx context.Context, id int64) error
}

func closeRows(ctx context.Context, rows *sql.Rows, logger *logrus.Logger) {
	err := rows.Close()
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to close rows: %v", err)
	}
}
