package message

import "time"

type Message struct {
	Id          int64
	Content     string
	PhoneNumber string
	Sent        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
