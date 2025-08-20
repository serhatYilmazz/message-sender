package message

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

type PgRepository struct {
	Db     *sql.DB
	Logger *logrus.Logger
}

func (r *PgRepository) FindAllMessages(ctx context.Context) ([]Message, error) {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][FindAllMessages] is called")
	query := "SELECT id, content, phone_number, sent, created_at FROM messages;"
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Error("error while querying for all messages: ")
		return nil, err
	}
	defer closeRows(ctx, rows, r.Logger)

	var messages = make([]Message, 0)
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.Content, &message.PhoneNumber, &message.Sent, &message.CreatedAt)
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
