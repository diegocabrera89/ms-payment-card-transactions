package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/models"
	"github.com/diegocabrera89/ms-payment-card-transactions/utils"
	"github.com/diegocabrera89/ms-payment-core/dynamodbcore"
	"github.com/diegocabrera89/ms-payment-core/helpers"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"github.com/google/uuid"
	"os"
)

// RefundCardRepositoryImpl implements the RefundCardRepository interface of the ms-payment-core package.
type RefundCardRepositoryImpl struct {
	CoreRepository dynamodbcore.CoreRepository
}

// NewRefundCardRepository create a new RefundCardRepository instance.
func NewRefundCardRepository() *RefundCardRepositoryImpl {
	transactionTable := os.Getenv(constantsmicro.TransactionTable)
	region := os.Getenv(constantsmicro.Region)

	coreRepository, _ := dynamodbcore.NewDynamoDBRepository(transactionTable, region)

	return &RefundCardRepositoryImpl{
		CoreRepository: coreRepository,
	}
}

// RefundCardRepository put item in DynamoDB.
func (r *RefundCardRepositoryImpl) RefundCardRepository(ctx context.Context, request events.APIGatewayProxyRequest, chargeCard models.ChargeCard) (models.ChargeCard, error) {
	logs.LogTrackingInfo("RefundCardRepository", ctx, request)
	validMerchant, errorMerchantChargeCard := utils.ValidatePublicPrivateMerchantIdRepository(ctx, request, constantsmicro.StatusTokenEnable)
	if errorMerchantChargeCard != nil {
		return models.ChargeCard{}, fmt.Errorf("No merchant found with this key")
	}

	if validMerchant.Data.MerchantID != "" {
		logs.LogTrackingInfoData("RefundCardRepository ValidatePublicMerchantIdRepository", validMerchant, ctx, request)

		validateToken := utils.ValidateTokenRepository(ctx, request, constantsmicro.StatusTokenEnable)
		logs.LogTrackingInfoData("RefundCardRepository validateToken", validateToken, ctx, request)
		if !validateToken {
			return models.ChargeCard{}, fmt.Errorf("Token invalid")
		}

		getProcessorByMerchantID, errorGetProcessorByMerchantID := utils.GetProcessorByMerchantIDRepository(ctx, request, constantsmicro.StatusTokenEnable, validMerchant.Data.MerchantID)
		logs.LogTrackingInfoData("RefundCardRepository getProcessorByMerchantID", getProcessorByMerchantID.DataProcessor, ctx, request)
		if errorGetProcessorByMerchantID != nil {
			logs.LogTrackingError("RefundCardRepository", "GetProcessorByMerchantIDRepository", ctx, request, errorGetProcessorByMerchantID)
			return models.ChargeCard{}, errorGetProcessorByMerchantID
		}

		getBankID, errorGetBankID := utils.GetBankIDRepository(ctx, request, constantsmicro.StatusRefund, getProcessorByMerchantID.DataProcessor.BankID)
		logs.LogTrackingInfoData("RefundCardRepository getBankID", getBankID, ctx, request)
		logs.LogTrackingInfoData("RefundCardRepository getBankID.DataBank", getBankID.DataBank, ctx, request)
		if errorGetBankID != nil {
			logs.LogTrackingError("RefundCardRepository", "errorGetBankID", ctx, request, errorGetBankID)
			return models.ChargeCard{}, errorGetBankID
		}

		logs.LogTrackingInfoData("RefundCardRepository CalculateChargeCard getBankID.DataBank.TotalAmount", getBankID.DataBank.TotalAmount, ctx, request)
		logs.LogTrackingInfoData("RefundCardRepository CalculateChargeCard chargeCard.Amount.SubtotalIva0", chargeCard.Amount.SubtotalIva0, ctx, request)

		currentAmountRefund, errorCalculateChargeCard := utils.CalculateRefundCard(ctx, request, getBankID.DataBank.TotalAmount, chargeCard.Amount.SubtotalIva0)
		logs.LogTrackingInfoData("RefundCardRepository currentAmountRefund", currentAmountRefund, ctx, request)
		if errorCalculateChargeCard != nil {
			logs.LogTrackingError("RefundCardRepository", "errorCalculateChargeCard", ctx, request, errorCalculateChargeCard)
			return models.ChargeCard{}, errorCalculateChargeCard
		}

		chargeCard.CurrentAmount = utils.Round(currentAmountRefund, 2)
		chargeCard.Status = constantsmicro.StatusRefund
		item, errorMarshallItem := helpers.MarshallItem(chargeCard)
		logs.LogTrackingInfoData("RefundCardRepository ValidatePublicMerchantIdRepository item", item, ctx, request)
		if errorMarshallItem != nil {
			logs.LogTrackingError("RefundCardRepository", "MarshallItem", ctx, request, errorMarshallItem)
			return models.ChargeCard{}, errorMarshallItem
		}

		errorPutItemCore := r.CoreRepository.PutItemCore(ctx, request, item)
		if errorPutItemCore != nil {
			logs.LogTrackingError("RefundCardRepository", "errorPutItemCore", ctx, request, errorCalculateChargeCard)
			return models.ChargeCard{}, errorPutItemCore
		}
		var bank models.Bank
		bank.BankID = getBankID.DataBank.BankID
		bank.RefundID = uuid.New().String() // Generate a unique ID for the bank
		bank.Status = constantsmicro.StatusRefund
		bank.TotalAmount = utils.Round(currentAmountRefund, 2)
		logs.LogTrackingInfoData("RefundCardRepository bank", bank, ctx, request)

		updateAmountBank, errorUpdateAmountBank := utils.UpdateAmountBank(ctx, request, bank)
		logs.LogTrackingInfoData("RefundCardRepository updateAmountBank", updateAmountBank, ctx, request)
		if errorUpdateAmountBank != nil {
			logs.LogTrackingError("RefundCardRepository", "errorUpdateAmountBank", ctx, request, errorUpdateAmountBank)
			return models.ChargeCard{}, errorUpdateAmountBank
		}
		return chargeCard, nil
	}
	return models.ChargeCard{}, fmt.Errorf("Problem with charge")
}
