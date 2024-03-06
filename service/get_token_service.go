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

// GetTokenService represents the service for the GetTokenService entity.
type GetTokenService struct {
	getRepo *repository.GetTokenRepositoryImpl
}

// NewGetTokenService create a new GetTokenService instance.
func NewGetTokenService() *GetTokenService {
	return &GetTokenService{
		getRepo: repository.NewGetTokenRepository(),
	}
}

// GetTokenService handles the creation of a new pet.
func (r *GetTokenService) GetTokenService(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("GetTokenService", ctx, request)
	logs.LogTrackingInfoData("GetTokenService request", request, ctx, request)
	var responseBody []byte
	tokenID := request.PathParameters[constantsmicro.TokenID]
	logs.LogTrackingInfoData("GetTokenService tokenID", tokenID, ctx, request)
	if tokenID == "" {
		logs.LogTrackingError("GetTokenService", "PathParameters", ctx, request, nil)
		return response.ErrorResponse(http.StatusBadRequest, constantscore.ErrorCreatingItem)
	}

	getToken, err := r.getRepo.GetTokenRepository(ctx, request, tokenID)
	if err != nil {
		logs.LogTrackingError("GetTokenService", "GetTokenRepository", ctx, request, err)
		return response.ErrorResponse(http.StatusBadRequest, constantscore.ErrorCreatingItem)
	}

	if getToken.TokenID != "" {
		responseBody, err = json.Marshal(getToken)
		if err != nil {
			logs.LogTrackingError("GetTokenService", "JSON Marshal", ctx, request, err)
			return response.ErrorResponse(http.StatusInternalServerError, constantscore.InvalidResponseBody)
		}
		return response.SuccessResponse(http.StatusOK, responseBody, constantscore.ItemSuccessfullyObtained)
	}
	return response.SuccessResponse(http.StatusOK, responseBody, constantscore.DataNotFound)
}
