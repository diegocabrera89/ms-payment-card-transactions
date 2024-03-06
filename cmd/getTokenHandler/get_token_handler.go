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

// GetTokenHandler handles HTTP requests related to the token entity.
type GetTokenHandler struct {
	geTokenService *service.GetTokenService
}

// NewGetTokenHandler create a new GetTokenHandler instance.
func NewGetTokenHandler() *GetTokenHandler {
	return &GetTokenHandler{
		geTokenService: service.NewGetTokenService(),
	}
}

// GetTokenHandler handler for createGetTokenHandler new token.
func (h *GetTokenHandler) GetTokenHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("GetTokenHandler", ctx, request)
	createGetTokenHandler, errorGetTokenHandler := h.geTokenService.GetTokenService(ctx, request)
	if errorGetTokenHandler != nil {
		logs.LogTrackingError("GetTokenHandler", "GetTokenService", ctx, request, errorGetTokenHandler)
		return response.ErrorResponse(http.StatusInternalServerError, constantsmicro.ErrorCreatingToken)
	}
	return createGetTokenHandler, nil
}

func main() {
	// Create an instance of PetHandler in the main function.
	GetTokenHandler := NewGetTokenHandler()

	// Wrap the handler function with logging middleware.
	handlerWithLogging := metadata.MiddlewareMetadata(GetTokenHandler.GetTokenHandler)

	// Start the Lambda handler with the handler function wrapped in the middleware.
	lambda.Start(handlerWithLogging)
}
