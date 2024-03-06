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

// CardTransactionHandler handles HTTP requests related to the token entity.
type CardTransactionHandler struct {
	getCardTransactionService *service.GetCardTransactionService
}

// NewGetTokenHandler create a new CardTransactionHandler instance.
func NewGetTokenHandler() *CardTransactionHandler {
	return &CardTransactionHandler{
		getCardTransactionService: service.NewGetCardTransactionService(),
	}
}

// GetCardTransactionHandler handler for createGetTokenHandler new token.
func (h *CardTransactionHandler) GetCardTransactionHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("GetCardTransactionHandler", ctx, request)
	getCardTransaction, errorGetCardTransaction := h.getCardTransactionService.GetCardTransactionService(ctx, request)
	if errorGetCardTransaction != nil {
		logs.LogTrackingError("GetCardTransactionHandler", "errorGetCardTransaction", ctx, request, errorGetCardTransaction)
		return response.ErrorResponse(http.StatusInternalServerError, constantsmicro.ErrorCreatingToken)
	}
	return getCardTransaction, nil
}

func main() {
	// Create an instance of PetHandler in the main function.
	getCardTransactionHandler := NewGetTokenHandler()

	// Wrap the handler function with logging middleware.
	handlerWithLogging := metadata.MiddlewareMetadata(getCardTransactionHandler.GetCardTransactionHandler)

	// Start the Lambda handler with the handler function wrapped in the middleware.
	lambda.Start(handlerWithLogging)
}
