package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/serhatYilmazz/message-sender/internal/config"
	"github.com/serhatYilmazz/message-sender/internal/outbox"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Sender interface {
	SendMessage(ctx context.Context, entry outbox.OutboxEntry) (*Response, error)
}

type sender struct {
	config     config.WebhookConfig
	httpClient *http.Client
	logger     *logrus.Logger
}

func NewSender(config config.WebhookConfig, logger *logrus.Logger) Sender {
	return &sender{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger,
	}
}

func (s *sender) SendMessage(ctx context.Context, entry outbox.OutboxEntry) (*Response, error) {
	s.logger.WithContext(ctx).Debugf("[webhook.sender][SendMessage] sending message with ID: %d", entry.Id)

	var messagePayload outbox.MessagePayload
	err := json.Unmarshal(entry.Payload, &messagePayload)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Errorf("failed to unmarshal payload for outbox entry ID: %d", entry.Id)
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	webhookPayload := map[string]interface{}{
		"id":          messagePayload.Id,
		"content":     messagePayload.Content,
		"phoneNumber": messagePayload.PhoneNumber,
		"timestamp":   entry.CreatedAt.Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(webhookPayload)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Errorf("failed to marshal webhook payload for outbox entry ID: %d", entry.Id)
		return nil, fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Errorf("failed to create HTTP request for outbox entry ID: %d", entry.Id)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Errorf("failed to send webhook for outbox entry ID: %d", entry.Id)
		return nil, fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.logger.WithContext(ctx).Errorf("webhook returned non-success status %d for outbox entry ID: %d", resp.StatusCode, entry.Id)
		return nil, fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		s.logger.WithContext(ctx).Errorf("webhook response could not be decoded")
		return nil, err
	}

	s.logger.WithContext(ctx).WithField("outbox_id", entry.Id).
		WithField("response", response).
		WithField("message_id", messagePayload.Id).
		Info("webhook sent successfully")

	return &response, nil
}
