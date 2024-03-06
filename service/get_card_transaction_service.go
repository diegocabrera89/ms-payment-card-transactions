package service

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/repository"
	"github.com/diegocabrera89/ms-payment-core/constantscore"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"github.com/diegocabrera89/ms-payment-core/response"
	"net/http"
)

// GetCardTransactionService represents the service for the GetCardTransactionService entity.
type GetCardTransactionService struct {
	getCardTransactionRepo *repository.GetCardTransactionRepositoryImpl
}

// NewGetCardTransactionService create a new GetCardTransactionService instance.
func NewGetCardTransactionService() *GetCardTransactionService {
	return &GetCardTransactionService{
		getCardTransactionRepo: repository.NewGetCardTransactionRepository(),
	}
}

// GetCardTransactionService handles the creation of a new pet.
func (r *GetCardTransactionService) GetCardTransactionService(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("GetCardTransactionService", ctx, request)
	logs.LogTrackingInfoData("GetCardTransactionService request", request, ctx, request)
	var responseBody []byte
	tokenID := request.PathParameters[constantsmicro.TokenID]
	logs.LogTrackingInfoData("GetCardTransactionService tokenID", tokenID, ctx, request)
	if tokenID == "" {
		logs.LogTrackingError("GetCardTransactionService", "PathParameters", ctx, request, nil)
		return response.ErrorResponse(http.StatusBadRequest, constantscore.ErrorCreatingItem)
	}

	getToken, err := r.getCardTransactionRepo.GetCardTransactionRepository(ctx, request, tokenID)
	logs.LogTrackingInfoData("GetCardTransactionService info getToken", getToken, ctx, request)
	if err != nil {
		logs.LogTrackingError("GetCardTransactionService", "GetCardTransactionRepository", ctx, request, err)
		return response.ErrorResponse(http.StatusBadRequest, constantscore.ErrorCreatingItem)
	}

	logs.LogTrackingInfoData("GetCardTransactionService info getToken.ChargeCardID", getToken[0].ChargeCardID, ctx, request)
	if getToken[0].ChargeCardID != "" {
		responseBody, err = json.Marshal(getToken)
		if err != nil {
			logs.LogTrackingError("GetCardTransactionService", "JSON Marshal", ctx, request, err)
			return response.ErrorResponse(http.StatusInternalServerError, constantscore.InvalidResponseBody)
		}
		return response.SuccessResponse(http.StatusOK, responseBody, constantscore.ItemSuccessfullyObtained)
	}
	return response.SuccessResponse(http.StatusOK, responseBody, constantscore.DataNotFound)
}
