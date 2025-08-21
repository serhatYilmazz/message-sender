package message

import "time"

type Message struct {
	Id          string    `json:"id"`
	Content     string    `json:"content"`
	PhoneNumber string    `json:"phoneNumber"`
	Sent        bool      `json:"sent"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
