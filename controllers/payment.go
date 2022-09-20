package controllers

import (
	"donutBackend/config"
	"donutBackend/models/users"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/paytm/Paytm_Go_Checksum/paytm"
)

var merchantId string;
var merchantKey string;

func init(){
	merchantId = config.Paytm.MerchantID
	merchantKey = config.Paytm.MerchantKey
}

func InitiatePayment(w http.ResponseWriter, r *http.Request) {

	amount:=r.FormValue("amount")
	custId:=r.FormValue("custID")
	timeStamp:=time.Now().Format("20060102150405")
	orderId := "ORDR_"+custId+"_"+timeStamp
	
	//Generating CheckSum
	body:=fmt.Sprintf(`{"requestType":"Payment","mid":"%s","orderId":"%s","websiteName":"WEBSTAGING","txnAmount":{"value":"%s","currency":"INR"},"userInfo":{"custId":"%s"},"callbackUrl":"http://localhost:8080/payment/callback"}`,merchantId,orderId,amount,custId)
	paytmCheckSum:=PaytmChecksum.GenerateSignatureByString(body,merchantKey)  
  
	//Creating Request
	url := fmt.Sprintf("https://securegw-stage.paytm.in/theia/api/v1/initiateTransaction?mid=%s&orderId=%s",merchantId,orderId)
	payload := strings.NewReader(fmt.Sprintf(`{"body":{"requestType":"Payment","mid":"%s","orderId":"%s","websiteName":"WEBSTAGING","txnAmount":{"value":"%s","currency":"INR"},"userInfo":{"custId":"%s"},"callbackUrl":"http://localhost:8080/payment/callback"},"head":{"signature":"%s"}}`,merchantId,orderId,amount,custId,paytmCheckSum))
  
	client := &http.Client {
	}
	req, err := http.NewRequest("POST", url, payload)
  
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
  
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	defer res.Body.Close()
  
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
  
	var result interface{}
	if err := json.Unmarshal(resBody, &result); err != nil {   // Parse []byte to go struct pointer
		fmt.Fprint(w, err.Error())
		return
	}
  
	//Reading Response
	txnToken:=result.(map[string]interface{})["body"].(map[string]interface{})["txnToken"].(string)
	returnedCheckSum:=result.(map[string]interface{})["head"].(map[string]interface{})["signature"].(string)
  
	//Verifying CheckSum
	resBodyBody:=strings.Split(string(resBody),`"body":`)[1]
	resBodyBody=resBodyBody[:len(resBodyBody)-1]
	verifiedCheckSum:=PaytmChecksum.VerifySignatureByString(resBodyBody,merchantKey,returnedCheckSum)
  
	if(!verifiedCheckSum){
	  resPayload := map[string]string{"error": "Checksum Verification Failed"}
	  response, err := json.Marshal(resPayload)
	  if err != nil {
		fmt.Fprint(w, err.Error())
		return
	  }
	  w.Header().Set("Content-Type", "application/json")
	  w.WriteHeader(http.StatusCreated)
	  w.Write(response)
	  return
	}
  
	//Creating Transaction
	amt, err := strconv.ParseFloat(amount, 64);
	if  err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	txn:=&users.Transaction{
	  MerchantId: merchantId,
	  OrderId: orderId,
	  Amount: amt,
	  Type: "PAYTM",
	  Status: "PENDING",
	  Timestamp: timeStamp,
	  Mode: "ONLINE",
	}
  
	resPayload := map[string]string{"merchantID":merchantId,"orderID":orderId,"txnToken": txnToken}
	response, err := json.Marshal(resPayload)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	_,err=users.InsertTransaction(custId,txn)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
  
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return
}

func VerifyPaymentStatus(w http.ResponseWriter, r *http.Request) {

	custId:=r.FormValue("custId")
	orderId:=r.FormValue("orderId")

	//Generating CheckSum
	body:=fmt.Sprintf(`{"mid":"%s","orderId":"%s"}`,merchantId,orderId)
	paytmCheckSum:=PaytmChecksum.GenerateSignatureByString(body,merchantKey) 

	url := "https://securegw-stage.paytm.in/merchant-status/api/v1/getPaymentStatus"
	payload := strings.NewReader(fmt.Sprintf(`{"body":{"mid":"%s","orderId":"%s"},"head":{"signature":"%s"}}`,merchantId,orderId,paytmCheckSum))

	client := &http.Client {
	}
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

	resBodyBody:=strings.Split(string(resBody),`"body":`)[1]
	resBodyBody=resBodyBody[:len(resBodyBody)-1]

	var result interface{}
	if err := json.Unmarshal(resBody, &result); err != nil {   // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	resultStatus:=result.(map[string]interface{})["body"].(map[string]interface{})["resultInfo"].(map[string]interface{})["resultStatus"].(string)
	returnedCheckSum:=result.(map[string]interface{})["head"].(map[string]interface{})["signature"].(string)

	verifiedCheckSum:=PaytmChecksum.VerifySignatureByString(resBodyBody,merchantKey,returnedCheckSum)

	if(!verifiedCheckSum){
		resPayload := map[string]string{"error": "Checksum Verification Failed"}
		response, err := json.Marshal(resPayload)
		if err != nil {
		  fmt.Fprint(w, err.Error())
		  return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
		return
	}	


	resPayload := map[string]string{"merchantId":merchantId,"orderId":orderId,"resultStatus": resultStatus}
	response, err := json.Marshal(resPayload)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	_,err=users.UpdatePaymentStatus(custId,orderId,resultStatus)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
  
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return
}