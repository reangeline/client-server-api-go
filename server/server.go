package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Dolar struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", GetApi)
	fmt.Println("Server started on port:8080")
	http.ListenAndServe(":8080", nil)
}

func GetApi(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var d Dolar
	error := json.Unmarshal(body, &d)
	if error != nil {
		panic(error)
	}

	json.NewEncoder(w).Encode(&d.Usdbrl.Bid)

}
