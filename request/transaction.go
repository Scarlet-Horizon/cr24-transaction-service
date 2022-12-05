package request

type TransactionRequest struct {
	// Sender account UUID
	SenderAccountID string `json:"senderAccountID" example:"5d84ca00-c079-4577-9560-e1014086affe"`
	// Recipient UUID
	RecipientID string `json:"recipientID" example:"495d45e9-644c-40b8-94e8-103cad128331"`
	// RecipientID account UUID
	RecipientAccountID string `json:"recipientAccountID" example:"8cca0453-8e84-4f3b-aa40-7fc9cd162a34"`
	// Transaction amount
	Amount float64 `json:"amount" example:"17.24" minimum:"1"`
	// Transaction type ID
	Type int `json:"type" example:"1"`
} //@name TransactionRequest
