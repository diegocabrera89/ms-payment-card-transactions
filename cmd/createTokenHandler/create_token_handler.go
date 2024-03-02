package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/service"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"github.com/diegocabrera89/ms-payment-core/middleware/metadata"
	"github.com/diegocabrera89/ms-payment-core/response"
	"net/http"
)

// TokenHandler handles HTTP requests related to the token entity.
type TokenHandler struct {
	tokenService *service.TokenService
}

// NewTokenHandler create a new TokenHandler instance.
func NewTokenHandler() *TokenHandler {
	return &TokenHandler{
		tokenService: service.NewTokenService(),
	}
}

// CreateTokenHandler handler for createTokenHandler new token.
func (h *TokenHandler) CreateTokenHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("CreateTokenHandler", ctx, request)
	createTokenHandler, errorTokenHandler := h.tokenService.CreateTokenService(ctx, request)
	if errorTokenHandler != nil {
		logs.LogTrackingError("CreateTokenHandler", "CreateTokenService", ctx, request, errorTokenHandler)
		return response.ErrorResponse(http.StatusInternalServerError, constantsmicro.ErrorCreatingToken)
	}
	return createTokenHandler, nil
}

func main() {
	// Create an instance of PetHandler in the main function.
	petHandler := NewTokenHandler()

	// Wrap the handler function with logging middleware.
	handlerWithLogging := metadata.MiddlewareMetadata(petHandler.CreateTokenHandler)

	// Start the Lambda handler with the handler function wrapped in the middleware.
	lambda.Start(handlerWithLogging)
}
