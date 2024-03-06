package models

// Bank structure to define Bank fields.
type Bank struct {
	BankID      string  `json:"bankID" dynamodbav:"bankID"`
	AccountID   string  `json:"accountID" dynamodbav:"accountID"`
	RefundID    string  `json:"refundID" dynamodbav:"refundID"`
	TotalAmount float64 `json:"totalAmount" dynamodbav:"totalAmount"`
	CreatedAt   int64   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt   int64   `json:"updatedAt" dynamodbav:"updatedAt"`
	Status      string  `json:"status" dynamodbav:"status"`
}
