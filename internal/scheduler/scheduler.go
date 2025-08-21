package scheduler

import (
	"context"
	"github.com/serhatYilmazz/message-sender/internal/config"
	"github.com/serhatYilmazz/message-sender/internal/outbox"
	"github.com/serhatYilmazz/message-sender/internal/webhook"
	"github.com/sirupsen/logrus"
	"time"
)

type Scheduler interface {
	Start(ctx context.Context) error
	Stop() error
}

type scheduler struct {
	config        config.SchedulerConfig
	outboxService outbox.Service
	webhookSender webhook.Sender
	logger        *logrus.Logger
	stopChan      chan struct{}
}

func NewScheduler(
	config config.SchedulerConfig,
	outboxService outbox.Service,
	webhookSender webhook.Sender,
	logger *logrus.Logger,
) Scheduler {
	return &scheduler{
		config:        config,
		outboxService: outboxService,
		webhookSender: webhookSender,
		logger:        logger,
		stopChan:      make(chan struct{}),
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	s.logger.WithContext(ctx).Info("[scheduler][Start] starting outbox message scheduler")

	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	s.processOutboxEntries(ctx)
	for {
		select {
		case <-ctx.Done():
			s.logger.WithContext(ctx).Info("[scheduler][Start] context cancelled, stopping scheduler")
			return ctx.Err()
		case <-s.stopChan:
			s.logger.WithContext(ctx).Info("[scheduler][Start] stop signal received, stopping scheduler")
			return nil
		case <-ticker.C:
			s.processOutboxEntries(ctx)
		}
	}
}

func (s *scheduler) Stop() error {
	s.logger.Info("[scheduler][Stop] stopping scheduler")
	close(s.stopChan)
	return nil
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

	_, err := s.webhookSender.SendMessage(sendCtx, entry)

	return err
}
