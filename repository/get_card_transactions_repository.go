package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/models"
	"github.com/diegocabrera89/ms-payment-core/dynamodbcore"
	"github.com/diegocabrera89/ms-payment-core/helpers"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"os"
)

// GetCardTransactionRepositoryImpl implements the TokenRepository interface of the ms-payment-core package.
type GetCardTransactionRepositoryImpl struct {
	CoreRepository dynamodbcore.CoreRepository
}

// NewGetCardTransactionRepository create a new ProcessorRepository instance.
func NewGetCardTransactionRepository() *GetCardTransactionRepositoryImpl {
	tokenTable := os.Getenv(constantsmicro.TransactionTable)
	region := os.Getenv(constantsmicro.Region)

	coreRepository, _ := dynamodbcore.NewDynamoDBRepository(tokenTable, region)

	return &GetCardTransactionRepositoryImpl{
		CoreRepository: coreRepository,
	}
}

// GetCardTransactionRepository put item in DynamoDB.
func (r *GetCardTransactionRepositoryImpl) GetCardTransactionRepository(ctx context.Context, request events.APIGatewayProxyRequest, tokenID string) ([]models.ChargeCard, error) {
	logs.LogTrackingInfo("GetCardTransactionRepositoryImpl", ctx, request)
	getCardTransaction, errorCardTransaction := r.CoreRepository.GetItemByFieldCore(ctx, request, constantsmicro.Token, tokenID, constantsmicro.TokenIDIndex, "", "")
	if errorCardTransaction != nil {
		logs.LogTrackingError("GetCardTransactionRepositoryImpl", "GetItemByFieldCore", ctx, request, errorCardTransaction)
		return []models.ChargeCard{}, errorCardTransaction
	}
	logs.LogTrackingInfoData("GetCardTransactionRepositoryImpl", getCardTransaction, ctx, request)

	// Check if there is at least one element in the response
	if getCardTransaction.Count == 0 {
		return []models.ChargeCard{}, fmt.Errorf("No items found")
	}

	// Create a slice of models.ChargeCard instances
	merchants := make([]models.ChargeCard, len(getCardTransaction.Items))
	logs.LogTrackingInfoData("GetCardTransactionRepositoryImpl merchants", merchants, ctx, request)

	// Deserialize maps into models.ChargeCard instances
	for i, item := range getCardTransaction.Items {
		var m models.ChargeCard
		err := helpers.UnmarshalMapToType(item, &m)
		if err != nil {
			logs.LogTrackingError("GetCardTransactionRepositoryImpl", "UnmarshalMapToType", ctx, request, err)
			return []models.ChargeCard{}, err
		}
		merchants[i] = m
		logs.LogTrackingInfoData("GetCardTransactionRepositoryImpl processors[i]", merchants[i], ctx, request)
	}

	logs.LogTrackingInfoData("GetCardTransactionRepositoryImpl processors", merchants, ctx, request)

	return merchants, nil
}
