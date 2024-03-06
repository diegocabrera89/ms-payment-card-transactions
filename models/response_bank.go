package models

// DataBank structure to define DataBank fields.
type DataBank struct {
	BankID      string  `json:"bankID" dynamodbav:"bankID"`
	AccountID   string  `json:"accountID" dynamodbav:"accountID"`
	RefundID    string  `json:"refundID" dynamodbav:"refundID"`
	TotalAmount float64 `json:"totalAmount" dynamodbav:"totalAmount"`
	CreatedAt   int64   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt   int64   `json:"updatedAt" dynamodbav:"updatedAt"`
	Status      string  `json:"status" dynamodbav:"status"`
}

type ResponseBank struct {
	Status   int      `json:"status"`
	Message  string   `json:"message"`
	DataBank DataBank `json:"data"`
}
