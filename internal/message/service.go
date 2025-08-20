package message

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Service interface {
	FindAllMessages(ctx context.Context) ([]Message, error)
	MarkAsSent(ctx context.Context, id int64) error
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

func (s *service) FindAllMessages(ctx context.Context) ([]Message, error) {
	s.Logger.WithContext(ctx).Debug("[message.service][FindAllMessages] is called")
	return s.Repository.FindAllMessages(ctx)
}

func (s *service) MarkAsSent(ctx context.Context, id int64) error {
	s.Logger.WithContext(ctx).WithField("id", id).Debug("[message.service][MarkAsSent] is called")
	return s.Repository.MarkAsSent(ctx, id)
}
