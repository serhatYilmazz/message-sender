package scheduler

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type Manager interface {
	StartScheduler(ctx context.Context) error
	StopScheduler() error
	IsSchedulerRunning() bool
}

type manager struct {
	scheduler Scheduler
	logger    *logrus.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
}

func NewManager(scheduler Scheduler, logger *logrus.Logger) Manager {
	return &manager{
		scheduler: scheduler,
		logger:    logger,
	}
}

func (m *manager) StartScheduler(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.scheduler.IsRunning() {
		m.logger.Info("[scheduler.manager][StartScheduler] scheduler is already running")
		return nil
	}

	m.logger.Info("[scheduler.manager][StartScheduler] starting scheduler")

	m.ctx, m.cancel = context.WithCancel(ctx)

	go func() {
		if err := m.scheduler.Start(m.ctx); err != nil && !errors.Is(err, context.Canceled) {
			m.logger.WithError(err).Error("[scheduler.manager] scheduler stopped with error")
		}
	}()

	m.logger.Info("[scheduler.manager][StartScheduler] scheduler started successfully")
	return nil
}

func (m *manager) StopScheduler() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.scheduler.IsRunning() {
		m.logger.Info("[scheduler.manager][StopScheduler] scheduler is not running")
		return nil
	}

	m.logger.Info("[scheduler.manager][StopScheduler] stopping scheduler")

	if m.cancel != nil {
		m.cancel()
	}

	if err := m.scheduler.Stop(); err != nil {
		m.logger.WithError(err).Error("[scheduler.manager][StopScheduler] error stopping scheduler")
		return err
	}

	m.logger.Info("[scheduler.manager][StopScheduler] scheduler stopped successfully")
	return nil
}

func (m *manager) IsSchedulerRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.scheduler.IsRunning()
}
