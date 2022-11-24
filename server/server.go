package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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

type Cotation struct {
	cotation string
}

func NewCotation(cotation string) *Cotation {
	return &Cotation{
		cotation: cotation,
	}
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

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cotation := NewCotation(string(body))
	err = insert(db, cotation)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(&d.Usdbrl.Bid)

}

func insert(db *sql.DB, cotation *Cotation) error {

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS cotations (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, cotation TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("insert into cotations(cotation) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cotation.cotation)
	if err != nil {
		return err
	}
	return nil
}
