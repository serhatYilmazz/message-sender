package scheduler

import (
	"context"
	"github.com/serhatYilmazz/message-sender/pkg/model"
	"github.com/sirupsen/logrus"
)

type ControlService interface {
	ProcessMessageSender(ctx context.Context, request model.MessageSenderRequest) error
	GetSchedulerStatus(ctx context.Context) bool
}

type controlService struct {
	manager Manager
	logger  *logrus.Logger
}

func NewControlService(manager Manager, logger *logrus.Logger) ControlService {
	return &controlService{
		manager: manager,
		logger:  logger,
	}
}

func (c *controlService) ProcessMessageSender(ctx context.Context, request model.MessageSenderRequest) error {
	c.logger.WithContext(ctx).Debugf("[scheduler.control][ProcessMessageSender] request: %+v", request)

	switch request.IsMessageSenderEnabled {
	case true:
		c.logger.WithContext(ctx).Info("[scheduler.control][ProcessMessageSender] enabling message sender")
		return c.manager.StartScheduler(ctx)
	case false:
		c.logger.WithContext(ctx).Info("[scheduler.control][ProcessMessageSender] disabling message sender")
		return c.manager.StopScheduler()
	}

	return nil
}

func (c *controlService) GetSchedulerStatus(ctx context.Context) bool {
	c.logger.WithContext(ctx).Debug("[scheduler.control][GetSchedulerStatus] checking scheduler status")
	return c.manager.IsSchedulerRunning()
}
