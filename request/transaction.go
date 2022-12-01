package request

type TransactionRequest struct {
	// Sender account UUID
	SenderID string `json:"senderID" example:"5d84ca00-c079-4577-9560-e1014086affe"`
	// Recipient account UUID
	RecipientID string `json:"recipientID" example:"8cca0453-8e84-4f3b-aa40-7fc9cd162a34"`
	// Transaction amount
	Amount float64 `json:"amount" example:"17.24"`
	// Transaction type ID, check model.TransactionType
	Type int `json:"type" example:"1"`
} //@name TransactionRequest
