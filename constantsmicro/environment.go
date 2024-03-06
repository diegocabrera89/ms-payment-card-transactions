package constantsmicro

const (
	// TokenTable constant to define table name.
	TokenTable       = "TOKEN_DYNAMODB"
	TransactionTable = "TRANSACTION_DYNAMODB"
)

const (
	// Region is AWS Region.
	Region = "REGION"
)

const (
	// UrlValidatePublicMerchant endpoint for validate merchant
	UrlValidatePublicMerchant = "https://eji1t2sesj.execute-api.us-east-1.amazonaws.com/dev/v1/merchants/kinds/publicMerchant/"

	// UrlValidatePrivateMerchant endpoint for validate merchant
	UrlValidatePrivateMerchant = "https://eji1t2sesj.execute-api.us-east-1.amazonaws.com/dev/v1/merchants/kinds/privateMerchant/"

	// UrlValidateGetToken endpoint for validate processor
	UrlValidateGetToken = "https://h8bsuqhp9h.execute-api.us-east-1.amazonaws.com/dev/v1/card/token/"

	// UrlGetProcessorByMerchantID endpoint for validate processor
	UrlGetProcessorByMerchantID = "https://3qchm37pya.execute-api.us-east-1.amazonaws.com/dev/v1/processors/kinds/processor/"

	// UrlGetBankID endpoint for validate bank
	UrlGetBankID = "https://3aj9if0cye.execute-api.us-east-1.amazonaws.com/dev/v1/bank/account/"

	// UrlUpdateBankID endpoint for validate bank
	UrlUpdateBankID = "https://3aj9if0cye.execute-api.us-east-1.amazonaws.com/dev/v1/bank/account/"
)
