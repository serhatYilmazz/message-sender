package outbox

import (
	"encoding/json"
	"time"
)

type OutboxEntry struct {
	Id        int64           `json:"id"`
	MessageId string          `json:"messageId"`
	Payload   json.RawMessage `json:"payload"`
	Sent      bool            `json:"sent"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type MessagePayload struct {
	Id          string `json:"id"`
	Content     string `json:"content"`
	PhoneNumber string `json:"phoneNumber"`
}
