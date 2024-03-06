package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/models"
	"github.com/diegocabrera89/ms-payment-core/logs"
	"github.com/google/uuid"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// BuildCreateToken build processor object.
func BuildCreateToken(token *models.Token) {
	token.TokenID = uuid.New().String()       // Generate a unique ID for the token
	token.CreatedAt = time.Now().UTC().Unix() //Date in UTC
	token.Status = constantsmicro.StatusTokenEnable
}

// BuildCreateChargeCard build processor object.
func BuildCreateChargeCard(chargeCard *models.ChargeCard) {
	chargeCard.ChargeCardID = uuid.New().String()  // Generate a unique ID for the charge
	chargeCard.CreatedAt = time.Now().UTC().Unix() //Date in UTC
	chargeCard.Status = constantsmicro.StatusCharge
}

// ValidatePublicPrivateMerchantIdRepository get item in DynamoDB.
func ValidatePublicPrivateMerchantIdRepository(ctx context.Context, request events.APIGatewayProxyRequest, statusTokenEnable string) (models.Response, error) {
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
func ValidateTokenRepository(ctx context.Context, request events.APIGatewayProxyRequest, statusTokenEnable string) bool {
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

// GetProcessorByMerchantIDRepository get item in DynamoDB.
func GetProcessorByMerchantIDRepository(ctx context.Context, request events.APIGatewayProxyRequest, statusTokenEnable string, merchantID string) (models.ResponseProcessor, error) {
	logs.LogTrackingInfo("GetProcessorByMerchantIDRepository", ctx, request)
	logs.LogTrackingInfoData("GetProcessorByMerchantIDRepository merchantID", merchantID, ctx, request)

	url := constantsmicro.UrlGetProcessorByMerchantID + merchantID

	logs.LogTrackingInfoData("GetProcessorByMerchantIDRepository url", url, ctx, request)
	resp, err := http.Get(url)
	if err != nil {
		logs.LogTrackingError("GetProcessorByMerchantIDRepository", "Error when making request", ctx, request, err)
		return models.ResponseProcessor{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.LogTrackingError("GetProcessorByMerchantIDRepository", "Error reading response", ctx, request, err)
		return models.ResponseProcessor{}, err
	}

	if resp.StatusCode != http.StatusOK {
		logs.LogTrackingError("GetProcessorByMerchantIDRepository", "Error", ctx, request, err)
		return models.ResponseProcessor{}, err
	}

	var response models.ResponseProcessor

	err = json.Unmarshal(body, &response)
	if err != nil {
		logs.LogTrackingError("GetProcessorByMerchantIDRepository", "JSON Unmarshal", ctx, request, err)
		return models.ResponseProcessor{}, err
	}
	logs.LogTrackingInfoData("GetProcessorByMerchantIDRepository response", response, ctx, request)
	logs.LogTrackingInfoData("GetProcessorByMerchantIDRepository response.Data", response.DataProcessor, ctx, request)
	if response.DataProcessor.Status != statusTokenEnable {
		return models.ResponseProcessor{}, fmt.Errorf("Data not found")
	}
	return response, nil
}

// GetBankIDRepository get item in DynamoDB.
func GetBankIDRepository(ctx context.Context, request events.APIGatewayProxyRequest, statusTokenEnable string, bankID string) (models.ResponseBank, error) {
	logs.LogTrackingInfo("GetBankIDRepository", ctx, request)
	logs.LogTrackingInfoData("GetBankIDRepository merchantID", bankID, ctx, request)

	url := constantsmicro.UrlGetBankID + bankID

	logs.LogTrackingInfoData("GetBankIDRepository url", url, ctx, request)
	resp, err := http.Get(url)
	if err != nil {
		logs.LogTrackingError("GetBankIDRepository", "Error when making request", ctx, request, err)
		return models.ResponseBank{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.LogTrackingError("GetBankIDRepository", "Error reading response", ctx, request, err)
		return models.ResponseBank{}, err
	}

	if resp.StatusCode != http.StatusOK {
		logs.LogTrackingError("GetBankIDRepository", "Error", ctx, request, err)
		return models.ResponseBank{}, err
	}

	var response models.ResponseBank

	err = json.Unmarshal(body, &response)
	if err != nil {
		logs.LogTrackingError("GetBankIDRepository", "JSON Unmarshal", ctx, request, err)
		return models.ResponseBank{}, err
	}
	logs.LogTrackingInfoData("GetBankIDRepository response", response, ctx, request)
	logs.LogTrackingInfoData("GetBankIDRepository response.Data.DataBank", response.DataBank, ctx, request)
	logs.LogTrackingInfoData("GetBankIDRepository statusTokenEnable", statusTokenEnable, ctx, request)
	if response.DataBank.Status == constantsmicro.StatusDisable {
		return models.ResponseBank{}, fmt.Errorf("Data not found or disable")
	}
	return response, nil
}

// CalculateChargeCard calculate amount charge
func CalculateChargeCard(ctx context.Context, request events.APIGatewayProxyRequest, chargeValue float64, currentValue float64) (float64, error) {
	logs.LogTrackingInfoData("CalculateChargeCard chargeValue", chargeValue, ctx, request)
	logs.LogTrackingInfoData("CalculateChargeCard currentValue", currentValue, ctx, request)
	logs.LogTrackingInfoData("CalculateChargeCard chargeValue-currentValue", chargeValue-currentValue, ctx, request)
	if currentValue > chargeValue {
		return 0, fmt.Errorf("The charge amount cannot be greater than the account amount")
	}

	return chargeValue - currentValue, nil
}

// CalculateRefundCard calculate amount refund
func CalculateRefundCard(ctx context.Context, request events.APIGatewayProxyRequest, chargeValue float64, currentValue float64) (float64, error) {
	logs.LogTrackingInfoData("CalculateChargeCard chargeValue", chargeValue, ctx, request)
	logs.LogTrackingInfoData("CalculateChargeCard currentValue", currentValue, ctx, request)
	logs.LogTrackingInfoData("CalculateChargeCard chargeValue-currentValue", chargeValue-currentValue, ctx, request)
	if currentValue > chargeValue {
		return 0, fmt.Errorf("The charge amount cannot be greater than the account amount")
	}

	return chargeValue + currentValue, nil
}

// UpdateAmountBank get item in DynamoDB.
func UpdateAmountBank(ctx context.Context, request events.APIGatewayProxyRequest, bank models.Bank) (models.ResponseBank, error) {
	logs.LogTrackingInfo("UpdateAmountBank", ctx, request)
	logs.LogTrackingInfoData("UpdateAmountBank bank", bank, ctx, request)

	url := constantsmicro.UrlUpdateBankID

	// Convert data to JSON format
	jsonData, errMarshal := json.Marshal(bank)
	if errMarshal != nil {
		fmt.Println("Error converting data to JSON errMarshal", errMarshal)
		return models.ResponseBank{}, errMarshal
	}

	// Create a PUT request with JSON data
	req, errNewRequest := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if errNewRequest != nil {
		fmt.Println("Error creating HTTP request errNewRequest", errNewRequest)
		return models.ResponseBank{}, errNewRequest
	}

	// Set the header to indicate that the body is JSON
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		fmt.Println("Error when making HTTP request errDo", errDo)
		return models.ResponseBank{}, errDo
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("The request was not successful. Status code", resp.StatusCode)
		return models.ResponseBank{}, fmt.Errorf("The request was not successful. Status code")
	}

	// Read server response
	body, errReadAll := ioutil.ReadAll(resp.Body)
	if errReadAll != nil {
		fmt.Println("Error reading response errReadAll", errReadAll)
		return models.ResponseBank{}, fmt.Errorf("Error reading response")
	}

	// Print server response
	fmt.Println("Server response", string(body))

	logs.LogTrackingInfoData("UpdateAmountBank body", body, ctx, request)
	var respBank models.ResponseBank
	err := json.Unmarshal(body, &respBank)
	if err != nil {
		logs.LogTrackingError("UpdateAmountBank", "JSON Unmarshal", ctx, request, err)
		return models.ResponseBank{}, err
	}
	respBank.DataBank = models.DataBank(bank)
	logs.LogTrackingInfoData("UpdateAmountBank respBank", respBank, ctx, request)
	logs.LogTrackingInfoData("UpdateAmountBank respBank.DataBank", respBank.DataBank, ctx, request)

	return respBank, nil
}

// Round decimal number
func Round(num float64, decimals int) float64 {
	expo := math.Pow(10, float64(decimals))
	return math.Round(num*expo) / expo
}
