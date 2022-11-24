package main

import (
	"context"
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	createFile(string(body))

}

func createFile(cotation string) {

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte("Dólar:" + cotation))
	if err != nil {
		panic(err)
	}
}
