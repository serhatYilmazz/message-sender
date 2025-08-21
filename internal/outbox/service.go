package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

type Service interface {
	CreateEntryForMessage(ctx context.Context, tx *sql.Tx, messageId, content, phoneNumber string) error
	ProcessUnsentEntries(ctx context.Context, limit int, processor func(ctx context.Context, entry OutboxEntry) error) (int, error)
	MarkEntriesAsSent(ctx context.Context, ids []int64) error
}

type service struct {
	repository Repository
	logger     *logrus.Logger
}

func NewService(repository Repository, logger *logrus.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) CreateEntryForMessage(ctx context.Context, tx *sql.Tx, messageId, content, phoneNumber string) error {
	s.logger.WithContext(ctx).Debugf("[outbox.service][CreateEntryForMessage] creating outbox entry for message: %s", messageId)

	payload := MessagePayload{
		Id:          messageId,
		Content:     content,
		PhoneNumber: phoneNumber,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("failed to marshal outbox payload")
		return err
	}

	outboxEntry := &OutboxEntry{
		MessageId: messageId,
		Payload:   payloadBytes,
		Sent:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repository.SaveOutboxEntry(ctx, tx, outboxEntry)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("failed to save outbox entry")
		return err
	}

	s.logger.WithContext(ctx).WithField("message_id", messageId).WithField("outbox_id", outboxEntry.Id).Info("outbox entry created successfully")
	return nil
}

func (s *service) ProcessUnsentEntries(ctx context.Context, limit int, processor func(ctx context.Context, entry OutboxEntry) error) (int, error) {
	s.logger.WithContext(ctx).Debugf("[outbox.service][ProcessUnsentEntries] processing unsent entries with limit: %d", limit)

	entries, err := s.repository.GetUnsentEntries(ctx, limit)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("failed to get unsent outbox entries")
		return 0, err
	}

	if len(entries) == 0 {
		s.logger.WithContext(ctx).Debug("no unsent outbox entries found")
		return 0, nil
	}

	s.logger.WithContext(ctx).Infof("processing %d unsent outbox entries", len(entries))

	var processedCount int
	var successfulIds []int64

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			s.logger.WithContext(ctx).Warn("processing cancelled due to context cancellation")
			break
		default:
			if err := processor(ctx, entry); err != nil {
				s.logger.WithContext(ctx).WithError(err).
					WithField("outbox_id", entry.Id).
					WithField("message_id", entry.MessageId).
					Error("failed to process outbox entry")
				continue
			}

			successfulIds = append(successfulIds, entry.Id)
			processedCount++
		}
	}

	if len(successfulIds) > 0 {
		if err := s.repository.MarkAsSent(ctx, successfulIds); err != nil {
			s.logger.WithContext(ctx).WithError(err).
				WithField("successful_ids", successfulIds).
				Error("failed to mark outbox entries as sent")
			return processedCount, err
		}

		s.logger.WithContext(ctx).
			WithField("processed_count", processedCount).
			WithField("total_entries", len(entries)).
			Info("outbox entries processed successfully")
	}

	return processedCount, nil
}

func (s *service) MarkEntriesAsSent(ctx context.Context, ids []int64) error {
	s.logger.WithContext(ctx).Debugf("[outbox.service][MarkEntriesAsSent] marking %d entries as sent", len(ids))

	if len(ids) == 0 {
		return nil
	}

	err := s.repository.MarkAsSent(ctx, ids)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).WithField("ids", ids).Error("failed to mark entries as sent")
		return err
	}

	s.logger.WithContext(ctx).WithField("count", len(ids)).Info("entries marked as sent successfully")
	return nil
}
