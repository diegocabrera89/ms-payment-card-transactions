package repository

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/models"
	"github.com/diegocabrera89/ms-payment-core/dynamodbcore"
	"github.com/diegocabrera89/ms-payment-core/helpers"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"os"
)

// TokenRepositoryImpl implements the TokenRepository interface of the ms-payment-core package.
type TokenRepositoryImpl struct {
	CoreRepository dynamodbcore.CoreRepository
}

// NewTokenRepository create a new ProcessorRepository instance.
func NewTokenRepository() *TokenRepositoryImpl {
	tokenTable := os.Getenv(constantsmicro.TokenTable)
	region := os.Getenv(constantsmicro.Region)

	coreRepository, _ := dynamodbcore.NewDynamoDBRepository(tokenTable, region)

	return &TokenRepositoryImpl{
		CoreRepository: coreRepository,
	}
}

// MerchantRepositoryImpl implements the MerchantRepository interface of the ms-payment-core package.
type MerchantRepositoryImpl struct {
	CoreRepository dynamodbcore.CoreRepository
}

// NewMerchantRepository create a new MerchantRepository instance.
func NewMerchantRepository() *MerchantRepositoryImpl {
	merchantTable := os.Getenv("ms-payment-merchant-dev-merchant")
	region := os.Getenv(constantsmicro.Region)

	coreRepository, _ := dynamodbcore.NewDynamoDBRepository(merchantTable, region)

	return &MerchantRepositoryImpl{
		CoreRepository: coreRepository,
	}
}

// CreateTokenRepository put item in DynamoDB.
func (r *TokenRepositoryImpl) CreateTokenRepository(ctx context.Context, request events.APIGatewayProxyRequest, processor models.Token) (models.Token, error) {
	logs.LogTrackingInfo("CreateTokenRepository", ctx, request)
	logs.LogTrackingInfoData("CreateTokenRepository", request.Headers["public-id"], ctx, request)
	item, errorMarshallItem := helpers.MarshallItem(processor)
	if errorMarshallItem != nil {
		logs.LogTrackingError("CreateTokenRepository", "MarshallItem", ctx, request, errorMarshallItem)
		return models.Token{}, errorMarshallItem
	}

	errorPutItemCore := r.CoreRepository.PutItemCore(ctx, request, item)
	if errorPutItemCore != nil {
		return models.Token{}, errorPutItemCore
	}
	return processor, nil
}

// GetMerchantByIdRepository get item in DynamoDB.
func (r *MerchantRepositoryImpl) GetMerchantByIdRepository(ctx context.Context, request events.APIGatewayProxyRequest, fieldValueFilterByID string) (models.Token, error) {
	logs.LogTrackingInfo("GetMerchantByIdRepository", ctx, request)
	var pet models.Token
	responseGetPetById, errorGetPetById := r.CoreRepository.GetItemCore(ctx, request, constantsmicro.TokenID, fieldValueFilterByID)
	if errorGetPetById != nil {
		logs.LogTrackingError("GetMerchantByIdRepository", "GetItemCore", ctx, request, errorGetPetById)
		return models.Token{}, errorGetPetById
	}
	errUnmarshalMap := helpers.UnmarshalMapToType(responseGetPetById.Item, &pet)
	if errUnmarshalMap != nil {
		logs.LogTrackingError("GetMerchantByIdRepository", "UnmarshalMap", ctx, request, errUnmarshalMap)
		return models.Token{}, errUnmarshalMap
	}
	return pet, nil
}
