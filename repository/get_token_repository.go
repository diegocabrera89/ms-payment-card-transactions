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

// GetTokenRepositoryImpl implements the TokenRepository interface of the ms-payment-core package.
type GetTokenRepositoryImpl struct {
	CoreRepository dynamodbcore.CoreRepository
}

// NewGetTokenRepository create a new ProcessorRepository instance.
func NewGetTokenRepository() *GetTokenRepositoryImpl {
	tokenTable := os.Getenv(constantsmicro.TokenTable)
	region := os.Getenv(constantsmicro.Region)

	coreRepository, _ := dynamodbcore.NewDynamoDBRepository(tokenTable, region)

	return &GetTokenRepositoryImpl{
		CoreRepository: coreRepository,
	}
}

// GetTokenRepository get item in DynamoDB.
func (r *GetTokenRepositoryImpl) GetTokenRepository(ctx context.Context, request events.APIGatewayProxyRequest, fieldValueFilterByID string) (models.Token, error) {
	logs.LogTrackingInfo("GetTokenRepository", ctx, request)
	var token models.Token
	responseGetToken, errorGetPetById := r.CoreRepository.GetItemCore(ctx, request, constantsmicro.TokenID, fieldValueFilterByID)
	if errorGetPetById != nil {
		logs.LogTrackingError("GetTokenRepository", "GetItemCore", ctx, request, errorGetPetById)
		return models.Token{}, errorGetPetById
	}
	errUnmarshalMap := helpers.UnmarshalMapToType(responseGetToken.Item, &token)
	if errUnmarshalMap != nil {
		logs.LogTrackingError("GetTokenRepository", "UnmarshalMap", ctx, request, errUnmarshalMap)
		return models.Token{}, errUnmarshalMap
	}
	return token, nil
}
