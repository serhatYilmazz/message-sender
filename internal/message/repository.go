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
	MarkAsSent(ctx context.Context, id int64) error
	SaveMessage(ctx context.Context, request model.AddMessageRequest) (*model.MessageDto, error)
}

func closeRows(ctx context.Context, rows *sql.Rows, logger *logrus.Logger) {
	err := rows.Close()
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to close rows: %v", err)
	}
}
