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

	// Initialize scheduler
	outboxScheduler := scheduler.NewScheduler(
		cfg.SchedulerConfig,
		outboxService,
		webhookSender,
		logger,
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Start scheduler in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting outbox scheduler...")
		if err := outboxScheduler.Start(ctx); err != nil && err != context.Canceled {
			logger.WithError(err).Error("scheduler stopped with error")
		}
	}()

	// Start API server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting API server...")
		api.NewMessageHandler(messageService, logger)
	}()

	// Wait for shutdown signal
	<-sigChan
	logger.Info("shutdown signal received, stopping services...")

	// Cancel context to stop scheduler
	cancel()

	// Stop scheduler explicitly
	if err := outboxScheduler.Stop(); err != nil {
		logger.WithError(err).Error("error stopping scheduler")
	}

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info("all services stopped gracefully")
}
