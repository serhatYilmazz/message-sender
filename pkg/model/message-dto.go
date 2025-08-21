package model

import "time"

type MessageDto struct {
	Id          string    `json:"id"`
	Content     string    `json:"content"`
	PhoneNumber string    `json:"phoneNumber"`
	Sent        bool      `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
