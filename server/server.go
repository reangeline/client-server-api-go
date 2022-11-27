package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	createDB()
	http.HandleFunc("/cotacao", GetApi)
	fmt.Println("Server started on port:8080")
	http.ListenAndServe(":8080", nil)
}
func createDB() {
	db := conn()
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS cotations(cotation TEXT)")
	if err != nil {
		panic(err)
	}
	db.Close()
}

func conn() *sql.DB {
	db, err := sql.Open("sqlite3", "db/dolar.db")
	if err != nil {
		panic(err)
	}
	return db
}

func GetApi(w http.ResponseWriter, r *http.Request) {

	c := http.Client{Timeout: time.Millisecond * 200}

	resp, err := c.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
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

	db, err := sql.Open("sqlite3", "dolar.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cotation := NewCotation(string(body))
	err = insert(cotation)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(&d.Usdbrl.Bid)

}

func insert(cotation *Cotation) error {

	db := conn()
	defer db.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO cotations(cotation) values(?)", cotation.cotation)
	if err != nil {
		fmt.Println("erro")
		return err
	}

	return nil
}
