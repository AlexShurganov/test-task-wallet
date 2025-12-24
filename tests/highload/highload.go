package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sync"
)

var (
	wg sync.WaitGroup
)

func sendRequest(url string, body []byte) {
	defer wg.Done()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error", err, body)
	}
	if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("bad request", string(body))
		errorBody, _ := io.ReadAll(resp.Body)
		fmt.Println("error response", string(errorBody))
	}
	defer resp.Body.Close()
}

func main() {
	wallets := [10]string{
		"f1b98811-f11f-44d4-a702-cd659b9e1c8c",
		"ce89c91b-cc5b-4a37-bd9e-10c89052ef2c",
		"63d99c9f-f215-4da7-a66a-1a45e993d342",
		"5735b475-e6a0-427d-be1b-dca3caccf9a8",
		"2da3ea67-25d1-456e-86e8-714c54b730fa",
		"f249d58e-a6f8-4797-8633-d0146f62bcab",
		"49244334-360b-4d85-9ae0-b0d06afdec65",
		"da9cff91-9a4d-42f2-9944-f2f6a1afb05e",
		"91f5c8d7-38b4-48ea-b0e8-00b2d248b809",
		"edd24d29-ea08-4f05-a75a-5b5143eb9f8b"}

	url := "http://localhost:8888/api/v1/wallet"
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		wallet := fmt.Sprintf("%s", wallets[rand.Intn(len(wallets))])
		op := []string{"DEPOSIT", "WITHDRAW"}[rand.Intn(2)]
		amount := rand.Float64()*1000 + 1
		amount = math.Round(amount*100) / 100
		body, _ := json.Marshal(map[string]interface{}{"walletId": wallet, "operationType": op, "amount": amount})
		go sendRequest(url, body)
	}
	wg.Wait()
}
