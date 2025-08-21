package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/serhatYilmazz/message-sender/internal/config"
	"github.com/serhatYilmazz/message-sender/internal/outbox"
	"github.com/serhatYilmazz/message-sender/internal/webhook"
	"github.com/sirupsen/logrus"
)

type Service interface {
	RecordWebhookDelivery(ctx context.Context, outboxEntry outbox.OutboxEntry, webhookResponse *webhook.Response) error
	GetDeliveryRecord(ctx context.Context, messageId string) (*WebhookDelivery, error)
	IsMessageProcessed(ctx context.Context, messageId string) (bool, error)
}

type service struct {
	repository Repository
	config     config.RedisConfig
	logger     *logrus.Logger
}

func NewService(repository Repository, config config.RedisConfig, logger *logrus.Logger) Service {
	return &service{
		repository: repository,
		config:     config,
		logger:     logger,
	}
}

func (s *service) RecordWebhookDelivery(ctx context.Context, outboxEntry outbox.OutboxEntry, webhookResponse *webhook.Response) error {
	s.logger.WithContext(ctx).
		WithField("outbox_id", outboxEntry.Id).
		WithField("message_id", outboxEntry.MessageId).
		WithField("webhook_response_message_id", webhookResponse.MessageId).
		Debug("recording webhook delivery in cache")

	delivery := &WebhookDelivery{
		MessageId:       webhookResponse.MessageId,
		OutboxMessageId: outboxEntry.MessageId,
		DeliveredAt:     time.Now(),
		Response:        webhookResponse.Message,
	}

	if err := s.repository.StoreWebhookDelivery(ctx, delivery, s.config.TTL); err != nil {
		s.logger.WithContext(ctx).WithError(err).Errorf("failed to store webhook delivery record for message ID: %s", delivery.MessageId)
		return fmt.Errorf("failed to store webhook delivery record: %w", err)
	}

	s.logger.WithContext(ctx).
		WithField("message_id", delivery.MessageId).
		WithField("outbox_id", outboxEntry.Id).
		Info("successfully recorded webhook delivery in cache")

	return nil
}

func (s *service) GetDeliveryRecord(ctx context.Context, messageId string) (*WebhookDelivery, error) {
	s.logger.WithContext(ctx).
		WithField("message_id", messageId).
		Debug("retrieving webhook delivery record from cache")

	delivery, err := s.repository.GetWebhookDelivery(ctx, messageId)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Errorf("failed to retrieve webhook delivery record for message ID: %s", messageId)
		return nil, fmt.Errorf("failed to retrieve webhook delivery record: %w", err)
	}

	if delivery == nil {
		s.logger.WithContext(ctx).
			WithField("message_id", messageId).
			Debug("webhook delivery record not found in cache")
		return nil, nil
	}

	s.logger.WithContext(ctx).
		WithField("message_id", messageId).
		Debug("successfully retrieved webhook delivery record from cache")

	return delivery, nil
}

func (s *service) IsMessageProcessed(ctx context.Context, messageId string) (bool, error) {
	s.logger.WithContext(ctx).
		WithField("message_id", messageId).
		Debug("checking if message has been processed")

	delivery, err := s.GetDeliveryRecord(ctx, messageId)
	if err != nil {
		return false, err
	}

	isProcessed := delivery != nil
	s.logger.WithContext(ctx).
		WithField("message_id", messageId).
		WithField("is_processed", isProcessed).
		Debug("checked message processing status")

	return isProcessed, nil
}
