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
	query := "SELECT id, content, phone_number, created_at, updated_at FROM messages FOR UPDATE;"
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Error("error while querying for all messages: ")
		return nil, err
	}
	defer closeRows(ctx, rows, r.Logger)

	var messages = make([]model.MessageDto, 0)
	for rows.Next() {
		var message model.MessageDto
		err := rows.Scan(&message.Id, &message.Content, &message.PhoneNumber, &message.CreatedAt, &message.UpdatedAt)
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

func (r *PgRepository) SaveMessageWithTx(ctx context.Context, tx *sql.Tx, request model.AddMessageRequest) (*model.MessageDto, error) {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][SaveMessageWithTx] is called")
	query := `INSERT INTO messages (id, content, phone_number, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5)`
	id := uuid.New().String()
	now := time.Now()
	_, err := tx.ExecContext(ctx, query, id, request.Content, request.RecipientPhoneNumber, now, now)
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
	}, nil
}

func (r *PgRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][BeginTransaction] is called")
	return r.Db.BeginTx(ctx, nil)
}
