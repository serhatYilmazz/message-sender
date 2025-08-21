package model

type AddMessageResponse struct {
	Id                   string `json:"id"`
	Content              string `json:"content"`
	RecipientPhoneNumber string `json:"recipientPhoneNumber"`
}
