package service

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-card-transactions/models"
	"github.com/diegocabrera89/ms-payment-card-transactions/repository"
	"github.com/diegocabrera89/ms-payment-card-transactions/utils"
	"github.com/diegocabrera89/ms-payment-core/constantscore"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"github.com/diegocabrera89/ms-payment-core/response"
	"net/http"
)

// ChargeCardService represents the service for the ChargeCardService entity.
type ChargeCardService struct {
	chargeCardRepo *repository.ChargeCardRepositoryImpl
}

// NewChargeCardService create a new TokenService instance.
func NewChargeCardService() *ChargeCardService {
	return &ChargeCardService{
		chargeCardRepo: repository.NewChargeCardRepository(),
	}
}

// ChargeCardService handles the creation of a new token.
func (r *ChargeCardService) ChargeCardService(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("ChargeCardService", ctx, request)
	var chargeCard models.ChargeCard
	err := json.Unmarshal([]byte(request.Body), &chargeCard)
	if err != nil {
		logs.LogTrackingError("ChargeCardService", "JSON Unmarshal", ctx, request, err)
		return response.ErrorResponse(http.StatusBadRequest, constantscore.InvalidRequestBody)
	}

	utils.BuildCreateChargeCard(&chargeCard)

	chargeCardResponse, errorChargeCardRepository := r.chargeCardRepo.ChargeCardRepository(ctx, request, chargeCard)
	if errorChargeCardRepository != nil {
		logs.LogTrackingError("ChargeCardService", "ChargeCardRepository", ctx, request, errorChargeCardRepository)
		return response.ErrorResponse(http.StatusInternalServerError, constantscore.ErrorCreatingItem)
	}

	responseBody, err := json.Marshal(chargeCardResponse)
	if err != nil {
		logs.LogTrackingError("ChargeCardService", "JSON Marshal", ctx, request, err)
		return response.ErrorResponse(http.StatusInternalServerError, constantscore.InvalidResponseBody)
	}
	return response.SuccessResponse(http.StatusCreated, responseBody, constantscore.ItemCreatedSuccessfully)
}
