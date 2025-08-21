package message

import (
	"context"
	"github.com/serhatYilmazz/message-sender/internal/outbox"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
)

type Service interface {
	FindAllMessages(ctx context.Context) ([]model.MessageDto, error)
	SaveMessage(ctx context.Context, request model.AddMessageRequest) (*model.MessageDto, error)
}

type service struct {
	Repository    Repository
	OutboxService outbox.Service
	Logger        *logrus.Logger
}

func NewMessageService(repository Repository, outboxService outbox.Service, logger *logrus.Logger) Service {
	return &service{
		Repository:    repository,
		OutboxService: outboxService,
		Logger:        logger,
	}
}

func (s *service) FindAllMessages(ctx context.Context) ([]model.MessageDto, error) {
	s.Logger.WithContext(ctx).Debug("[message.service][FindAllMessages] is called")
	return s.Repository.FindAllMessages(ctx)
}

func (s *service) SaveMessage(ctx context.Context, request model.AddMessageRequest) (*model.MessageDto, error) {
	s.Logger.WithContext(ctx).Debugf("[message.service][SaveMessage] is called with %v", request)

	tx, err := s.Repository.BeginTransaction(ctx)
	if err != nil {
		s.Logger.WithContext(ctx).WithError(err).Error("failed to begin transaction")
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.Logger.WithContext(ctx).WithError(rollbackErr).Error("failed to rollback transaction")
			}
		}
	}()

	savedMessage, err := s.Repository.SaveMessageWithTx(ctx, tx, request)
	if err != nil {
		s.Logger.WithContext(ctx).WithError(err).Error("failed to save message")
		return nil, err
	}

	err = s.OutboxService.CreateEntryForMessage(ctx, tx, savedMessage.Id, savedMessage.Content, savedMessage.PhoneNumber)
	if err != nil {
		s.Logger.WithContext(ctx).WithError(err).Error("failed to create outbox entry")
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		s.Logger.WithContext(ctx).WithError(err).Error("failed to commit transaction")
		return nil, err
	}

	s.Logger.WithContext(ctx).WithField("message_id", savedMessage.Id).Info("message and outbox entry saved successfully")
	return savedMessage, nil
}
