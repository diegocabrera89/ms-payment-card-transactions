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

// RefundCardHandler handles HTTP requests related to the token entity.
type RefundCardHandler struct {
	refundCardService *service.RefundCardService
}

// NewRefundCardHandler create a new RefundCardHandler instance.
func NewRefundCardHandler() *RefundCardHandler {
	return &RefundCardHandler{
		refundCardService: service.NewRefundCardService(),
	}
}

// RefundCardHandler handler for RefundCardHandler new token.
func (h *RefundCardHandler) RefundCardHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logs.LogTrackingInfo("RefundCardHandler", ctx, request)
	refundCardHandler, errorRefundCardHandler := h.refundCardService.RefundCardService(ctx, request)
	if errorRefundCardHandler != nil {
		logs.LogTrackingError("RefundCardHandler", "RefundCardService", ctx, request, errorRefundCardHandler)
		return response.ErrorResponse(http.StatusInternalServerError, constantsmicro.ErrorCreatingToken)
	}
	return refundCardHandler, nil
}

func main() {
	// Create an instance of PetHandler in the main function.
	refundCardHandler := NewRefundCardHandler()

	// Wrap the handler function with logging middleware.
	handlerWithLogging := metadata.MiddlewareMetadata(refundCardHandler.RefundCardHandler)

	// Start the Lambda handler with the handler function wrapped in the middleware.
	lambda.Start(handlerWithLogging)
}
