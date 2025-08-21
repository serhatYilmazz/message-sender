package scheduler

import (
	"context"
	"github.com/serhatYilmazz/message-sender/internal/cache"
	"github.com/serhatYilmazz/message-sender/internal/config"
	"github.com/serhatYilmazz/message-sender/internal/outbox"
	"github.com/serhatYilmazz/message-sender/internal/webhook"
	"github.com/sirupsen/logrus"
	"time"
)

type Scheduler interface {
	Start(ctx context.Context) error
	Stop() error
	IsRunning() bool
}

type scheduler struct {
	config        config.SchedulerConfig
	outboxService outbox.Service
	webhookSender webhook.Sender
	cacheService  cache.Service
	logger        *logrus.Logger
	stopChan      chan struct{}
	isRunning     bool
}

func NewScheduler(
	config config.SchedulerConfig,
	outboxService outbox.Service,
	webhookSender webhook.Sender,
	cacheService cache.Service,
	logger *logrus.Logger,
) Scheduler {
	return &scheduler{
		config:        config,
		outboxService: outboxService,
		webhookSender: webhookSender,
		cacheService:  cacheService,
		logger:        logger,
		stopChan:      make(chan struct{}),
		isRunning:     false,
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	if s.isRunning {
		s.logger.WithContext(ctx).Info("[scheduler][Start] scheduler is already running")
		return nil
	}

	s.logger.WithContext(ctx).Info("[scheduler][Start] starting outbox message scheduler")
	s.isRunning = true

	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	s.processOutboxEntries(ctx)

	for {
		select {
		case <-ctx.Done():
			s.logger.WithContext(ctx).Info("[scheduler][Start] context cancelled, stopping scheduler")
			s.isRunning = false
			return ctx.Err()
		case <-s.stopChan:
			s.logger.WithContext(ctx).Info("[scheduler][Start] stop signal received, stopping scheduler")
			s.isRunning = false
			return nil
		case <-ticker.C:
			s.processOutboxEntries(ctx)
		}
	}
}

func (s *scheduler) Stop() error {
	if !s.isRunning {
		s.logger.Info("[scheduler][Stop] scheduler is not running")
		return nil
	}

	s.logger.Info("[scheduler][Stop] stopping scheduler")
	close(s.stopChan)
	s.isRunning = false
	return nil
}

func (s *scheduler) IsRunning() bool {
	return s.isRunning
}

func (s *scheduler) processOutboxEntries(ctx context.Context) {
	processingCtx, cancel := context.WithTimeout(ctx, s.config.SendTimeout)
	defer cancel()

	s.logger.WithContext(processingCtx).Debug("[scheduler][processOutboxEntries] processing outbox entries")

	processedCount, err := s.outboxService.ProcessUnsentEntries(processingCtx, s.config.BatchSize, s.sendMessage)
	if err != nil {
		s.logger.WithContext(processingCtx).WithError(err).Error("failed to process outbox entries")
		return
	}

	if processedCount > 0 {
		s.logger.WithContext(processingCtx).
			WithField("processed_count", processedCount).
			Info("successfully processed outbox entries in scheduler")
	}
}

func (s *scheduler) sendMessage(ctx context.Context, entry outbox.OutboxEntry) error {
	sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	s.logger.WithContext(sendCtx).
		WithField("outbox_id", entry.Id).
		WithField("message_id", entry.MessageId).
		Debug("sending message via webhook")

	response, err := s.webhookSender.SendMessage(sendCtx, entry)
	if err != nil {
		s.logger.WithContext(sendCtx).WithError(err).
			WithField("outbox_id", entry.Id).
			WithField("message_id", entry.MessageId).
			Error("failed to send webhook")
		return err
	}

	if response != nil {
		if cacheErr := s.cacheService.RecordWebhookDelivery(sendCtx, entry, response); cacheErr != nil {
			// Log cache error but don't fail the operation - webhook was successful
			s.logger.WithContext(sendCtx).WithError(cacheErr).
				WithField("outbox_id", entry.Id).
				WithField("message_id", entry.MessageId).
				WithField("webhook_response_message_id", response.MessageId).
				Warn("failed to record webhook delivery in cache, but webhook was successful")
		}
	}

	return nil
}
