package main

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Cotation struct {
	cotation string
}

func NewCotation(cotation string) *Cotation {
	return &Cotation{
		cotation: cotation,
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	cotation := NewCotation(string(body))
	err = insert(db, cotation)
	if err != nil {
		panic(err)
	}

	createFile(cotation)

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

func createFile(cotation *Cotation) {

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte("Dólar:" + cotation.cotation))
	if err != nil {
		panic(err)
	}
}
