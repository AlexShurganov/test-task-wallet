package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var (
	postURL = "http://localhost:8888/api/v1/wallet"
	getURL  = "http://localhost:8888/api/v1/wallets/"
)

func sendRequest(method, url string, body []byte) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error creating request", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error making request", err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("Http code: %d\n", resp.StatusCode)
	log.Printf("Response: %s\n", string(responseBody))
}

func main() {
	testUUID := "f1b98811-f11f-44d4-a702-cd659b9e1c8c"

	log.Println("DEPOSIT 1000")
	body, _ := json.Marshal(map[string]interface{}{
		"walletId": testUUID, "operationType": "DEPOSIT", "amount": 1000})
	sendRequest("POST", postURL, body)

	log.Println("WTHDRAW 500")
	body, _ = json.Marshal(map[string]interface{}{
		"walletId": testUUID, "operationType": "WITHDRAW", "amount": 500})
	sendRequest("POST", postURL, body)

	log.Println("WTHDRAW 2000")
	body, _ = json.Marshal(map[string]interface{}{
		"walletId": testUUID, "operationType": "WITHDRAW", "amount": 2000})
	sendRequest("POST", postURL, body)

	log.Println("Getting wallet info")
	sendRequest("GET", getURL+testUUID, nil)

}
