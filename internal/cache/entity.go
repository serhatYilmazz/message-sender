package cache

import (
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
