package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/donut-gives/backend/config"
	"github.com/donut-gives/backend/models/users"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/paytm/Paytm_Go_Checksum/paytm"
)

var merchantId string
var merchantKey string

func init() {
	merchantId = config.Payment.Paytm.MerchantId
	merchantKey = config.Payment.Paytm.MerchantKey
}

func InitiatePayment(c *gin.Context) {

	amount := c.PostForm("amount")
	custId := c.PostForm("custId")
	timeStamp := time.Now().Format("20060102150405")
	orderId := "ORDR_" + custId + "_" + timeStamp

	//Generating CheckSum
	body := fmt.Sprintf(`{"requestType":"Payment","mid":"%s","orderId":"%s","websiteName":"WEBSTAGING","txnAmount":{"value":"%s","currency":"INR"},"userInfo":{"custId":"%s"},"callbackUrl":"http://localhost:8080/payment/callback"}`, merchantId, orderId, amount, custId)
	paytmCheckSum := PaytmChecksum.GenerateSignatureByString(body, merchantKey)

	//Creating Request
	baseUrl := "https://securegw-stage.paytm.in/theia/api/v1/initiateTransaction?mid=%s&orderId=%s"
	if *config.Env == "prod" {
		baseUrl = "https://securegw.paytm.in/theia/api/v1/initiateTransaction?mid=%s&orderId=%s"
	}

	url := fmt.Sprintf(baseUrl, merchantId, orderId)

	payload := strings.NewReader(fmt.Sprintf(`{"body":{"requestType":"Payment","mid":"%s","orderId":"%s","websiteName":"WEBSTAGING","txnAmount":{"value":"%s","currency":"INR"},"userInfo":{"custId":"%s"},"callbackUrl":"http://localhost:8080/payment/callback"},"head":{"signature":"%s"}}`, merchantId, orderId, amount, custId, paytmCheckSum))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var result interface{}
	if err := json.Unmarshal(resBody, &result); err != nil { // Parse []byte to go struct pointer
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//Reading Response
	txnToken := result.(map[string]interface{})["body"].(map[string]interface{})["txnToken"].(string)
	returnedCheckSum := result.(map[string]interface{})["head"].(map[string]interface{})["signature"].(string)

	//Verifying CheckSum
	resBodyBody := strings.Split(string(resBody), `"body":`)[1]
	resBodyBody = resBodyBody[:len(resBodyBody)-1]
	verifiedCheckSum := PaytmChecksum.VerifySignatureByString(resBodyBody, merchantKey, returnedCheckSum)

	if !verifiedCheckSum {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to initiate payment",
			"error":   "Checksum Verification Failed",
		})
		return
	}

	//Creating Transaction
	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	txn := &users.Transaction{
		MerchantId: merchantId,
		OrderId:    orderId,
		Amount:     amt,
		Type:       "PAYTM",
		Status:     "PENDING",
		Timestamp:  timeStamp,
		Mode:       "ONLINE",
	}

	resPayload := map[string]string{"merchantID": merchantId, "orderID": orderId, "txnToken": txnToken}

	_, err = users.InsertTransaction(custId, txn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction Initiated",
		"data":    resPayload,
	})
}

func VerifyPaymentStatus(c *gin.Context) {

	custId := c.PostForm("custId")
	orderId := c.PostForm("orderId")

	//Generating CheckSum
	body := fmt.Sprintf(`{"mid":"%s","orderId":"%s"}`, merchantId, orderId)
	paytmCheckSum := PaytmChecksum.GenerateSignatureByString(body, merchantKey)

	url := "https://securegw-stage.paytm.in/merchant-status/api/v1/getPaymentStatus"
	if *config.Env == "prod" {
		url = "https://securegw.paytm.in/merchant-status/api/v1/getPaymentStatus"
	}

	payload := strings.NewReader(fmt.Sprintf(`{"body":{"mid":"%s","orderId":"%s"},"head":{"signature":"%s"}}`, merchantId, orderId, paytmCheckSum))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resBody))

	resBodyBody := strings.Split(string(resBody), `"body":`)[1]
	resBodyBody = resBodyBody[:len(resBodyBody)-1]

	var result interface{}
	if err := json.Unmarshal(resBody, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	resultStatus := result.(map[string]interface{})["body"].(map[string]interface{})["resultInfo"].(map[string]interface{})["resultStatus"].(string)
	returnedCheckSum := result.(map[string]interface{})["head"].(map[string]interface{})["signature"].(string)

	verifiedCheckSum := PaytmChecksum.VerifySignatureByString(resBodyBody, merchantKey, returnedCheckSum)

	if !verifiedCheckSum {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Fetch Payment Status",
			"error":   "Checksum Verification Failed",
		})
		return
	}

	resPayload := map[string]string{"merchantId": merchantId, "orderId": orderId, "resultStatus": resultStatus}

	_, err = users.UpdatePaymentStatus(custId, orderId, resultStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction Verification Status",
		"data":    resPayload,
	})
}
