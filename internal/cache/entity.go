package cache

import (
	"fmt"
	"time"
)

type WebhookDelivery struct {
	MessageId       string    `json:"messageId"`
	OutboxMessageId string    `json:"outboxMessageId"`
	DeliveredAt     time.Time `json:"deliveredAt"`
	Response        string    `json:"response,omitempty"`
}

func (wd *WebhookDelivery) CacheKey() string {
	return "webhook:delivery:" + wd.OutboxMessageId
}

func OutboxCacheKey(outboxId int64) string {
	return fmt.Sprintf("webhook:outbox:%d", outboxId)
}
