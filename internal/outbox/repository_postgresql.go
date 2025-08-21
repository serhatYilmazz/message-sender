package outbox

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

type PgRepository struct {
	Db     *sql.DB
	Logger *logrus.Logger
}

func (r *PgRepository) SaveOutboxEntry(ctx context.Context, tx *sql.Tx, entry *OutboxEntry) error {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][SaveOutboxEntry] is called for message_id: %s", entry.MessageId)

	query := `INSERT INTO outbox (message_id, payload, sent, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := tx.QueryRowContext(ctx, query,
		entry.MessageId,
		entry.Payload,
		entry.Sent,
		entry.CreatedAt,
		entry.UpdatedAt).Scan(&entry.Id)

	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Errorf("error while saving outbox entry for message_id: %s", entry.MessageId)
		return err
	}

	r.Logger.WithContext(ctx).WithField("outbox_id", entry.Id).Infof("outbox entry saved for message_id: %s", entry.MessageId)
	return nil
}

func (r *PgRepository) GetUnsentEntries(ctx context.Context, limit int) ([]OutboxEntry, error) {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][GetUnsentEntries] is called with limit: %d", limit)

	query := `SELECT id, message_id, payload, sent, created_at, updated_at 
			  FROM outbox 
			  WHERE sent = false 
			  ORDER BY created_at
			  LIMIT $1`

	rows, err := r.Db.QueryContext(ctx, query, limit)
	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Error("error while querying unsent outbox entries")
		return nil, err
	}
	defer closeRows(ctx, rows, r.Logger)

	var entries []OutboxEntry
	for rows.Next() {
		var entry OutboxEntry
		var payload []byte

		err := rows.Scan(&entry.Id, &entry.MessageId, &payload, &entry.Sent, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			r.Logger.WithContext(ctx).WithError(err).Error("error while scanning outbox entry")
			return nil, err
		}

		entry.Payload = payload
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		r.Logger.WithContext(ctx).WithError(err).Error("error during rows iteration")
		return nil, err
	}

	r.Logger.WithContext(ctx).Infof("retrieved %d unsent outbox entries", len(entries))
	return entries, nil
}

func (r *PgRepository) MarkAsSent(ctx context.Context, ids []int64) error {
	r.Logger.WithContext(ctx).Debugf("[PgRepository][MarkAsSent] is called for ids: %v", ids)

	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE outbox SET sent = true, updated_at = $1 WHERE id = ANY($2)`

	_, err := r.Db.ExecContext(ctx, query, time.Now(), pq.Array(ids))
	if err != nil {
		r.Logger.WithContext(ctx).WithError(err).Errorf("error while marking outbox entries as sent for ids: %v", ids)
		return err
	}

	r.Logger.WithContext(ctx).WithField("ids", ids).Infof("marked %d outbox entries as sent", len(ids))
	return nil
}
