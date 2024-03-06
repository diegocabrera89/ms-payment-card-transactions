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

// RefundCardService represents the service for the RefundCardService entity.
type RefundCardService struct {
	refundCardRepo *repository.RefundCardRepositoryImpl
}

// NewRefundCardService create a new TokenService instance.
func NewRefundCardService() *RefundCardService {
	return &RefundCardService{
		refundCardRepo: repository.NewRefundCardRepository(),
	}
}

// RefundCardService handles the creation of a new token.
func (r *RefundCardService) RefundCardService(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("RefundCardService", ctx, request)
	var chargeCard models.ChargeCard
	err := json.Unmarshal([]byte(request.Body), &chargeCard)
	if err != nil {
		logs.LogTrackingError("RefundCardService", "JSON Unmarshal", ctx, request, err)
		return response.ErrorResponse(http.StatusBadRequest, constantscore.InvalidRequestBody)
	}

	utils.BuildCreateChargeCard(&chargeCard)

	refundCardResponse, errorRefundCardRepository := r.refundCardRepo.RefundCardRepository(ctx, request, chargeCard)
	if errorRefundCardRepository != nil {
		logs.LogTrackingError("RefundCardService", "RefundCardRepository", ctx, request, errorRefundCardRepository)
		return response.ErrorResponse(http.StatusInternalServerError, constantscore.ErrorCreatingItem)
	}

	responseBody, err := json.Marshal(refundCardResponse)
	if err != nil {
		logs.LogTrackingError("RefundCardService", "JSON Marshal", ctx, request, err)
		return response.ErrorResponse(http.StatusInternalServerError, constantscore.InvalidResponseBody)
	}
	return response.SuccessResponse(http.StatusCreated, responseBody, constantscore.ItemCreatedSuccessfully)
}
