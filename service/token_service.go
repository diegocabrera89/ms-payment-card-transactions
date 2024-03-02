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

// TokenService represents the service for the TokenService entity.
type TokenService struct {
	tokenRepo *repository.TokenRepositoryImpl
}

// NewTokenService create a new TokenService instance.
func NewTokenService() *TokenService {
	return &TokenService{
		tokenRepo: repository.NewTokenRepository(),
	}
}

// CreateTokenService handles the creation of a new token.
func (r *TokenService) CreateTokenService(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("CreateTokenService", ctx, request)
	var token models.Token
	err := json.Unmarshal([]byte(request.Body), &token)
	if err != nil {
		logs.LogTrackingError("CreateTokenService", "JSON Unmarshal", ctx, request, err)
		return response.ErrorResponse(http.StatusBadRequest, constantscore.InvalidRequestBody)
	}

	utils.BuildCreateToken(&token)

	createToken, errorTokenRepository := r.tokenRepo.CreateTokenRepository(ctx, request, token)
	if errorTokenRepository != nil {
		logs.LogTrackingError("CreateTokenService", "CreateTokenRepository", ctx, request, errorTokenRepository)
		return response.ErrorResponse(http.StatusInternalServerError, constantscore.ErrorCreatingItem)
	}

	responseBody, err := json.Marshal(createToken)
	if err != nil {
		logs.LogTrackingError("CreateTokenService", "JSON Marshal", ctx, request, err)
		return response.ErrorResponse(http.StatusInternalServerError, constantscore.InvalidResponseBody)
	}
	return response.SuccessResponse(http.StatusCreated, responseBody, constantscore.ItemCreatedSuccessfully)
}
