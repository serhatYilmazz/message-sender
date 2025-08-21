// Message Sender API
//
// This is a sample Message Sender API server.
//
// @title Message Sender API

package main

import (
	"context"
	"github.com/serhatYilmazz/message-sender/api"
	"github.com/serhatYilmazz/message-sender/internal/config"
	"github.com/serhatYilmazz/message-sender/internal/message"
	"github.com/serhatYilmazz/message-sender/internal/outbox"
	"github.com/serhatYilmazz/message-sender/internal/scheduler"
	"github.com/serhatYilmazz/message-sender/internal/webhook"
	"github.com/serhatYilmazz/message-sender/pkg/db"
	"github.com/serhatYilmazz/message-sender/pkg/log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/serhatYilmazz/message-sender/docs"
)

//go:generate swag init
func main() {
	env := os.Getenv("ENVIRONMENT")
	logger := log.NewLogger(env)

	cfg, err := config.Load(logger, "./configs")
	if err != nil {
		return
	}

	postgresDb, err := db.NewPostgresDb(cfg.DbConfig)
	if err != nil {
		logger.Fatal("db connection is failed:", err)
	}

	// Initialize repositories
	pgMessageRepository := &message.PgRepository{
		Db:     postgresDb,
		Logger: logger,
	}

	pgOutboxRepository := &outbox.PgRepository{
		Db:     postgresDb,
		Logger: logger,
	}

	// Initialize services
	outboxService := outbox.NewService(pgOutboxRepository, logger)

	webhookSender := webhook.NewSender(cfg.WebhookConfig, logger)
	messageService := message.NewMessageService(pgMessageRepository, outboxService, logger)

	// Initialize scheduler components
	outboxScheduler := scheduler.NewScheduler(
		cfg.SchedulerConfig,
		outboxService,
		webhookSender,
		logger,
	)

	schedulerManager := scheduler.NewManager(outboxScheduler, logger)
	if cfg.SchedulerConfig.Enabled {
		err := schedulerManager.StartScheduler(context.Background())
		if err != nil {
			logger.WithError(err).Error("error while starting scheduler on startup")
		}
	}
	schedulerControlService := scheduler.NewControlService(schedulerManager, logger)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Start API server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting API server...")
		api.NewMessageHandler(messageService, schedulerControlService, logger)
	}()

	logger.Info("application started successfully. Use /api/messages/process-message-sender to control the scheduler")

	// Wait for shutdown signal
	<-sigChan
	logger.Info("shutdown signal received, stopping services...")

	// Stop scheduler if running
	if err := schedulerManager.StopScheduler(); err != nil {
		logger.WithError(err).Error("error stopping scheduler")
	}

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info("all services stopped gracefully")
}
