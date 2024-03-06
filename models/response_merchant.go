package models

// Data structure to define Data fields.
type Data struct {
	MerchantID     string `json:"merchantID" dynamodbav:"merchantID"`
	Name           string `json:"name" dynamodbav:"name"`
	Identification string `json:"identification" dynamodbav:"identification"`
	SocialReason   string `json:"socialReason" dynamodbav:"socialReason"`
	Url            string `json:"url" dynamodbav:"url"`
	PublicID       string `json:"publicID" dynamodbav:"publicID"`
	PrivateID      string `json:"privateID" dynamodbav:"privateID"`
	Country        string `json:"country" dynamodbav:"country"`
	CreatedAt      int64  `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt      int64  `json:"updatedAt" dynamodbav:"updatedAt"`
	StatusMerchant string `json:"statusMerchant" dynamodbav:"statusMerchant"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}
