package request

type TransactionRequest struct {
	// Sender UUID
	SenderID string `json:"senderID" example:"8182aadf-e376-4f01-b1d5-98d4e0a619ad"`
	// Sender account UUID
	SenderAccountID string `json:"senderAccountID" example:"5d84ca00-c079-4577-9560-e1014086affe"`
	// Recipient account UUID
	RecipientID string `json:"recipientID" example:"495d45e9-644c-40b8-94e8-103cad128331"`
	// RecipientID account UUID
	RecipientAccountID string `json:"recipientAccountID" example:"8cca0453-8e84-4f3b-aa40-7fc9cd162a34"`
	// Transaction amount
	Amount float64 `json:"amount" example:"17.24"`
	// Transaction type ID, check model.TransactionType
	Type int `json:"type" example:"1"`
} //@name TransactionRequest
