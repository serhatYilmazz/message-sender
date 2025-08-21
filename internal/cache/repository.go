package cache

import (
	"context"
	"time"
)

type Repository interface {
	StoreWebhookDelivery(ctx context.Context, delivery *WebhookDelivery, ttl time.Duration) error
	GetWebhookDelivery(ctx context.Context, messageId string) (*WebhookDelivery, error)
}
