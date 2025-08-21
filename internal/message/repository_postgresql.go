package message

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
	"time"
)

type PgRepository struct {
	Db     *sql.DB
	Logger *logrus.Logger
}

func (r *PgRepository) FindAllMessages(ctx context.Context) ([]model.MessageDto, error) {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][FindAllMessages] is called")
	query := "SELECT id, content, phone_number, sent, created_at, updated_at FROM messages;"
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Error("error while querying for all messages: ")
		return nil, err
	}
	defer closeRows(ctx, rows, r.Logger)

	var messages = make([]model.MessageDto, 0)
	for rows.Next() {
		var message model.MessageDto
		err := rows.Scan(&message.Id, &message.Content, &message.PhoneNumber, &message.Sent, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *PgRepository) MarkAsSent(ctx context.Context, id int64) error {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][MarkAsSent] is called")
	statement := `
			UPDATE messages
			SET sent = true
			where id = $1
			`
	_, err := r.Db.Exec(statement, id)
	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Error("error while updating message")
		return err
	}

	r.Logger.WithContext(ctx).WithField("id", id).Infof("message is updated as sent = true")
	return nil
}

func (r *PgRepository) SaveMessage(ctx context.Context, request model.AddMessageRequest) (*model.MessageDto, error) {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][SaveMessage] is called")
	query := `INSERT INTO messages (id, content, phone_number, sent, created_at, updated_at) 
				  VALUES ($1, $2, $3, $4, $5, $6)
				 `
	id := uuid.New().String()
	now := time.Now()
	_, err := r.Db.ExecContext(ctx, query, id, request.Content, request.RecipientPhoneNumber, false, now, now)
	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Errorf("error while saving data: %v", request)
		return nil, err
	}

	return &model.MessageDto{
		Id:          id,
		Content:     request.Content,
		PhoneNumber: request.RecipientPhoneNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
		Sent:        false,
	}, nil
}
