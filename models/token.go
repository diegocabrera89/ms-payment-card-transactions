package models

// Token structure to define Token fields.
type Token struct {
	TokenID     string  `json:"tokenID" dynamodbav:"tokenID"`
	Name        string  `json:"name" dynamodbav:"name"`
	Number      string  `json:"number" dynamodbav:"number"`
	ExpiryMonth string  `json:"expiryMonth" dynamodbav:"expiryMonth"`
	ExpiryYear  string  `json:"expiryYear" dynamodbav:"expiryYear"`
	Cvv         string  `json:"cvv" dynamodbav:"cvv"`
	TotalAmount float64 `json:"totalAmount" dynamodbav:"totalAmount"`
	Currency    string  `json:"currency" dynamodbav:"currency"`
	CreatedAt   int64   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt   int64   `json:"updatedAt" dynamodbav:"updatedAt"`
	Status      string  `json:"status" dynamodbav:"status"`
}
