package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	redisClient "github.com/serhatYilmazz/message-sender/pkg/redis"
	"github.com/sirupsen/logrus"
)

type redisRepository struct {
	client *redisClient.Client
	logger *logrus.Logger
}

func NewRedisRepository(client *redisClient.Client, logger *logrus.Logger) Repository {
	return &redisRepository{
		client: client,
		logger: logger,
	}
}

func (r *redisRepository) StoreWebhookDelivery(ctx context.Context, delivery *WebhookDelivery, ttl time.Duration) error {
	key := delivery.CacheKey()

	data, err := json.Marshal(delivery)
	if err != nil {
		r.logger.WithContext(ctx).WithError(err).Errorf("failed to marshal webhook delivery for message ID: %s", delivery.MessageId)
		return fmt.Errorf("failed to marshal webhook delivery: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		r.logger.WithContext(ctx).WithError(err).Errorf("failed to store webhook delivery in Redis for message ID: %s", delivery.MessageId)
		return fmt.Errorf("failed to store webhook delivery in Redis: %w", err)
	}

	r.logger.WithContext(ctx).
		WithField("message_id", delivery.MessageId).
		WithField("outbox_id", delivery.OutboxMessageId).
		WithField("ttl", ttl).
		Debug("successfully stored webhook delivery in cache")

	return nil
}

func (r *redisRepository) GetWebhookDelivery(ctx context.Context, messageId string) (*WebhookDelivery, error) {
	key := "webhook:delivery:" + messageId

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.WithContext(ctx).Debugf("webhook delivery not found in cache for message ID: %s", messageId)
			return nil, nil
		}
		r.logger.WithContext(ctx).WithError(err).Errorf("failed to get webhook delivery from Redis for message ID: %s", messageId)
		return nil, fmt.Errorf("failed to get webhook delivery from Redis: %w", err)
	}

	var delivery WebhookDelivery
	if err := json.Unmarshal([]byte(data), &delivery); err != nil {
		r.logger.WithContext(ctx).WithError(err).Errorf("failed to unmarshal webhook delivery for message ID: %s", messageId)
		return nil, fmt.Errorf("failed to unmarshal webhook delivery: %w", err)
	}

	r.logger.WithContext(ctx).
		WithField("message_id", messageId).
		Debug("successfully retrieved webhook delivery from cache")

	return &delivery, nil
}
