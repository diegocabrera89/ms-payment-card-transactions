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

// ChargeCardHandler handles HTTP requests related to the token entity.
type ChargeCardHandler struct {
	chargeCardService *service.ChargeCardService
}

// NewChargeCardHandler create a new ChargeCardHandler instance.
func NewChargeCardHandler() *ChargeCardHandler {
	return &ChargeCardHandler{
		chargeCardService: service.NewChargeCardService(),
	}
}

// CreateChargeCardHandler handler for createChargeCardHandler new token.
func (h *ChargeCardHandler) CreateChargeCardHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("CreateChargeCardHandler", ctx, request)
	createChargeCardHandler, errorChargeCardHandler := h.chargeCardService.ChargeCardService(ctx, request)
	if errorChargeCardHandler != nil {
		logs.LogTrackingError("CreateChargeCardHandler", "CreateTokenService", ctx, request, errorChargeCardHandler)
		return response.ErrorResponse(http.StatusInternalServerError, constantsmicro.ErrorCreatingToken)
	}
	return createChargeCardHandler, nil
}

func main() {
	// Create an instance of PetHandler in the main function.
	ChargeCardHandler := NewChargeCardHandler()

	// Wrap the handler function with logging middleware.
	handlerWithLogging := metadata.MiddlewareMetadata(ChargeCardHandler.CreateChargeCardHandler)

	// Start the Lambda handler with the handler function wrapped in the middleware.
	lambda.Start(handlerWithLogging)
}
