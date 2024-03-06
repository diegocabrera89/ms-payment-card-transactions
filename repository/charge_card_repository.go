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
	"os"
)

// ChargeCardRepositoryImpl implements the ChargeCardRepository interface of the ms-payment-core package.
type ChargeCardRepositoryImpl struct {
	CoreRepository dynamodbcore.CoreRepository
}

// NewChargeCardRepository create a new ChargeCardRepository instance.
func NewChargeCardRepository() *ChargeCardRepositoryImpl {
	tokenTable := os.Getenv(constantsmicro.TransactionTable)
	region := os.Getenv(constantsmicro.Region)

	coreRepository, _ := dynamodbcore.NewDynamoDBRepository(tokenTable, region)

	return &ChargeCardRepositoryImpl{
		CoreRepository: coreRepository,
	}
}

// ChargeCardRepository put item in DynamoDB.
func (r *ChargeCardRepositoryImpl) ChargeCardRepository(ctx context.Context, request events.APIGatewayProxyRequest, chargeCard models.ChargeCard) (models.ChargeCard, error) {
	logs.LogTrackingInfo("ChargeCardRepository", ctx, request)
	validMerchant, errorMerchantChargeCard := utils.ValidatePublicPrivateMerchantIdRepository(ctx, request, constantsmicro.StatusTokenEnable)
	if errorMerchantChargeCard != nil {
		return models.ChargeCard{}, fmt.Errorf("No merchant found with this key")
	}

	if validMerchant.Data.MerchantID != "" {
		logs.LogTrackingInfoData("ChargeCardRepository ValidatePublicMerchantIdRepository", validMerchant, ctx, request)

		validateToken := utils.ValidateTokenRepository(ctx, request, constantsmicro.StatusTokenEnable)
		logs.LogTrackingInfoData("ChargeCardRepository validateToken", validateToken, ctx, request)
		if !validateToken {
			return models.ChargeCard{}, fmt.Errorf("Token invalid")
		}

		getProcessorByMerchantID, errorGetProcessorByMerchantID := utils.GetProcessorByMerchantIDRepository(ctx, request, constantsmicro.StatusTokenEnable, validMerchant.Data.MerchantID)
		logs.LogTrackingInfoData("ChargeCardRepository getProcessorByMerchantID", getProcessorByMerchantID.DataProcessor, ctx, request)
		if errorGetProcessorByMerchantID != nil {
			logs.LogTrackingError("ChargeCardRepository", "GetProcessorByMerchantIDRepository", ctx, request, errorGetProcessorByMerchantID)
			return models.ChargeCard{}, errorGetProcessorByMerchantID
		}

		getBankID, errorGetBankID := utils.GetBankIDRepository(ctx, request, constantsmicro.StatusCharge, getProcessorByMerchantID.DataProcessor.BankID)
		logs.LogTrackingInfoData("ChargeCardRepository getBankID", getBankID, ctx, request)
		logs.LogTrackingInfoData("ChargeCardRepository getBankID.DataBank", getBankID.DataBank, ctx, request)
		if errorGetBankID != nil {
			logs.LogTrackingError("ChargeCardRepository", "errorGetBankID", ctx, request, errorGetBankID)
			return models.ChargeCard{}, errorGetBankID
		}

		logs.LogTrackingInfoData("ChargeCardRepository CalculateChargeCard getBankID.DataBank.TotalAmount", getBankID.DataBank.TotalAmount, ctx, request)
		logs.LogTrackingInfoData("ChargeCardRepository CalculateChargeCard chargeCard.Amount.SubtotalIva0", chargeCard.Amount.SubtotalIva0, ctx, request)

		currentAmount, errorCalculateChargeCard := utils.CalculateChargeCard(ctx, request, getBankID.DataBank.TotalAmount, chargeCard.Amount.SubtotalIva0)
		logs.LogTrackingInfoData("ChargeCardRepository currentAmount", currentAmount, ctx, request)
		if errorCalculateChargeCard != nil {
			logs.LogTrackingError("ChargeCardRepository", "errorCalculateChargeCard", ctx, request, errorCalculateChargeCard)
			return models.ChargeCard{}, errorCalculateChargeCard
		}

		chargeCard.CurrentAmount = utils.Round(currentAmount, 2)
		item, errorMarshallItem := helpers.MarshallItem(chargeCard)
		logs.LogTrackingInfoData("ChargeCardRepository ValidatePublicMerchantIdRepository item", item, ctx, request)
		if errorMarshallItem != nil {
			logs.LogTrackingError("ChargeCardRepository", "MarshallItem", ctx, request, errorMarshallItem)
			return models.ChargeCard{}, errorMarshallItem
		}

		errorPutItemCore := r.CoreRepository.PutItemCore(ctx, request, item)
		if errorPutItemCore != nil {
			logs.LogTrackingError("ChargeCardRepository", "errorPutItemCore", ctx, request, errorCalculateChargeCard)
			return models.ChargeCard{}, errorPutItemCore
		}
		var bank models.Bank
		bank.BankID = getBankID.DataBank.BankID
		bank.TotalAmount = utils.Round(currentAmount, 2)
		bank.Status = constantsmicro.StatusCharge
		logs.LogTrackingInfoData("ChargeCardRepository bank", bank, ctx, request)

		updateAmountBank, errorUpdateAmountBank := utils.UpdateAmountBank(ctx, request, bank)
		logs.LogTrackingInfoData("ChargeCardRepository updateAmountBank", updateAmountBank, ctx, request)
		if errorUpdateAmountBank != nil {
			logs.LogTrackingError("ChargeCardRepository", "errorUpdateAmountBank", ctx, request, errorUpdateAmountBank)
			return models.ChargeCard{}, errorUpdateAmountBank
		}
		return chargeCard, nil
	}
	return models.ChargeCard{}, fmt.Errorf("Problem with charge")
}
