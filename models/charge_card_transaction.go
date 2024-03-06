package models

// ChargeCardTransaction structure to define charge card fields.
type ChargeCardTransaction struct {
	ChargeCardTransactionID string  `json:"chargeCardTransactionID" dynamodbav:"chargeCardTransactionID"`
	Token                   string  `json:"token" dynamodbav:"token"`
	AmountCharge            Amount  `json:"amountCharge" dynamodbav:"amountCharge"`
	CurrentAmount           float64 `json:"currentAmount" dynamodbav:"currentAmount"`
	CreatedAt               int64   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt               int64   `json:"updatedAt" dynamodbav:"updatedAt"`
	Status                  string  `json:"status" dynamodbav:"status"`
}
