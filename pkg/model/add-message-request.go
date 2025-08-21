package model

type AddMessageRequest struct {
	Content              string `json:"content" validate:"required,max=20"`
	RecipientPhoneNumber string `json:"recipientPhoneNumber" validate:"required"`
}
