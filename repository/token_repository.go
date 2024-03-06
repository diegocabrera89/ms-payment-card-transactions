package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/models"
	"github.com/diegocabrera89/ms-payment-core/dynamodbcore"
	"github.com/diegocabrera89/ms-payment-core/helpers"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"io/ioutil"
	"net/http"
	"os"
)

// TokenRepositoryImpl implements the TokenRepository interface of the ms-payment-core package.
type TokenRepositoryImpl struct {
	CoreRepository dynamodbcore.CoreRepository
}

// NewTokenRepository create a new ProcessorRepository instance.
func NewTokenRepository() *TokenRepositoryImpl {
	tokenTable := os.Getenv(constantsmicro.TokenTable)
	region := os.Getenv(constantsmicro.Region)

	coreRepository, _ := dynamodbcore.NewDynamoDBRepository(tokenTable, region)

	return &TokenRepositoryImpl{
		CoreRepository: coreRepository,
	}
}

// CreateTokenRepository put item in DynamoDB.
func (r *TokenRepositoryImpl) CreateTokenRepository(ctx context.Context, request events.APIGatewayProxyRequest, processor models.Token) (models.Token, error) {
	logs.LogTrackingInfo("CreateTokenRepository", ctx, request)
	validMerchant, errorMerchantCreateToken := r.ValidatePublicPrivateMerchantIdRepository(ctx, request, constantsmicro.StatusTokenEnable)
	if errorMerchantCreateToken != nil {
		return models.Token{}, fmt.Errorf("No merchant found with this key")
	}

	if validMerchant.Data.MerchantID != "" {
		logs.LogTrackingInfoData("CreateTokenRepository validMerchant", validMerchant, ctx, request)
		item, errorMarshallItem := helpers.MarshallItem(processor)
		if errorMarshallItem != nil {
			logs.LogTrackingError("CreateTokenRepository", "MarshallItem", ctx, request, errorMarshallItem)
			return models.Token{}, errorMarshallItem
		}

		errorPutItemCore := r.CoreRepository.PutItemCore(ctx, request, item)
		if errorPutItemCore != nil {
			return models.Token{}, errorPutItemCore
		}
		return processor, nil
	}
	return models.Token{}, fmt.Errorf("No merchant found")
}

// ValidatePublicPrivateMerchantIdRepository get item in DynamoDB.
func (r *TokenRepositoryImpl) ValidatePublicPrivateMerchantIdRepository(ctx context.Context, request events.APIGatewayProxyRequest, statusTokenEnable string) (models.Response, error) {
	logs.LogTrackingInfo("ValidatePublicPrivateMerchantIdRepository", ctx, request)
	publicID := request.Headers[constantsmicro.PublicIDHeader]
	logs.LogTrackingInfoData("ValidatePublicPrivateMerchantIdRepository publicID", publicID, ctx, request)
	privateID := request.Headers[constantsmicro.PrivateIDHeader]
	logs.LogTrackingInfoData("ValidatePublicPrivateMerchantIdRepository privateID", privateID, ctx, request)
	var url string
	if publicID != "" {
		url = constantsmicro.UrlValidatePublicMerchant + publicID
	}

	if privateID != "" {
		url = constantsmicro.UrlValidatePrivateMerchant + privateID
	}
	logs.LogTrackingInfoData("ValidatePublicPrivateMerchantIdRepository url", url, ctx, request)
	resp, err := http.Get(url)
	if err != nil {
		logs.LogTrackingError("ValidatePublicPrivateMerchantIdRepository", "Error when making request", ctx, request, err)
		return models.Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.LogTrackingError("ValidatePublicPrivateMerchantIdRepository", "Error reading response", ctx, request, err)
		return models.Response{}, err
	}

	if resp.StatusCode != http.StatusOK {
		logs.LogTrackingError("ValidatePublicPrivateMerchantIdRepository", "Error", ctx, request, err)
		return models.Response{}, err
	}

	var response models.Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		logs.LogTrackingError("ValidatePublicPrivateMerchantIdRepository", "JSON Unmarshal", ctx, request, err)
		return models.Response{}, err
	}
	logs.LogTrackingInfoData("ValidatePublicPrivateMerchantIdRepository response", response, ctx, request)
	logs.LogTrackingInfoData("ValidatePublicPrivateMerchantIdRepository response.Data.", response.Data, ctx, request)
	if response.Data.StatusMerchant != statusTokenEnable {
		return models.Response{}, fmt.Errorf("Data not found")
	}
	return response, nil
}

// ValidateTokenRepository get item in DynamoDB.
func (r *TokenRepositoryImpl) ValidateTokenRepository(ctx context.Context, request events.APIGatewayProxyRequest, statusTokenEnable string) bool {
	logs.LogTrackingInfo("ValidateTokenRepository", ctx, request)
	var chargeCard models.ChargeCard
	err := json.Unmarshal([]byte(request.Body), &chargeCard)

	tokenID := chargeCard.Token

	logs.LogTrackingInfoData("ValidateTokenRepository tokenID", tokenID, ctx, request)

	url := constantsmicro.UrlValidateGetToken + tokenID

	logs.LogTrackingInfoData("ValidateTokenRepository url", url, ctx, request)
	resp, err := http.Get(url)
	if err != nil {
		logs.LogTrackingError("ValidateTokenRepository", "Error when making request", ctx, request, err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.LogTrackingError("ValidateTokenRepository", "Error reading response", ctx, request, err)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		logs.LogTrackingError("ValidateTokenRepository", "Error", ctx, request, err)
		return false
	}

	var response models.ResponseToken
	err = json.Unmarshal(body, &response)
	if err != nil {
		logs.LogTrackingError("ValidateTokenRepository", "JSON Unmarshal", ctx, request, err)
		return false
	}

	if response.Token.Status != statusTokenEnable {
		return false
	}
	return true
}
