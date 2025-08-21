package message

import (
	"context"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
)

type Service interface {
	FindAllMessages(ctx context.Context) ([]model.MessageDto, error)
	MarkAsSent(ctx context.Context, id int64) error
	ProcessMessageSender(ctx context.Context, request model.MessageSenderRequest) error
	SaveMessage(ctx context.Context, request model.AddMessageRequest) (*model.MessageDto, error)
}

type service struct {
	Repository Repository
	Logger     *logrus.Logger
}

func NewMessageService(repository Repository, logger *logrus.Logger) Service {
	return &service{
		Repository: repository,
		Logger:     logger,
	}
}

func (s *service) FindAllMessages(ctx context.Context) ([]model.MessageDto, error) {
	s.Logger.WithContext(ctx).Debug("[message.service][FindAllMessages] is called")
	return s.Repository.FindAllMessages(ctx)
}

func (s *service) MarkAsSent(ctx context.Context, id int64) error {
	s.Logger.WithContext(ctx).WithField("id", id).Debug("[message.service][MarkAsSent] is called")
	return s.Repository.MarkAsSent(ctx, id)
}

func (s *service) ProcessMessageSender(ctx context.Context, request model.MessageSenderRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) SaveMessage(ctx context.Context, request model.AddMessageRequest) (*model.MessageDto, error) {
	s.Logger.WithContext(ctx).Debugf("[message.service][SaveMessage] is called with %v", request)
	return s.Repository.SaveMessage(ctx, request)
}
