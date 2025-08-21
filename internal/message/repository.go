package message

import (
	"context"
	"database/sql"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	// FindAllMessages Limit would be specified
	FindAllMessages(ctx context.Context) ([]model.MessageDto, error)
	SaveMessageWithTx(ctx context.Context, tx *sql.Tx, request model.AddMessageRequest) (*model.MessageDto, error)
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

func closeRows(ctx context.Context, rows *sql.Rows, logger *logrus.Logger) {
	err := rows.Close()
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to close rows: %v", err)
	}
}
