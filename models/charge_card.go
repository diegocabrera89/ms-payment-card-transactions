package models

// ChargeCard structure to define charge card fields.
type ChargeCard struct {
	ChargeCardID  string  `json:"chargeCardID" dynamodbav:"chargeCardID"`
	Token         string  `json:"token" dynamodbav:"token"`
	Amount        Amount  `json:"amount" dynamodbav:"amount"`
	CurrentAmount float64 `json:"currentAmount" dynamodbav:"currentAmount"`
	CreatedAt     int64   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt     int64   `json:"updatedAt" dynamodbav:"updatedAt"`
	Status        string  `json:"status" dynamodbav:"status"`
}

type Amount struct {
	Currency     string  `json:"currency" dynamodbav:"currency"`
	SubtotalIva  float64 `json:"subtotalIva" dynamodbav:"subtotalIva"`
	SubtotalIva0 float64 `json:"subtotalIva0" dynamodbav:"subtotalIva0"`
	Iva          float64 `json:"iva" dynamodbav:"iva"`
	Ice          float64 `json:"ice" dynamodbav:"ice"`
}
