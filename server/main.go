package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRLCotacaoResponse struct {
	USDBRlCotacao Cotacao `json:"USDBRL"`
}

type Cotacao struct {
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	High      string `json:"high"`
	Low       string `json:"low"`
}

type CotacaoOutput struct {
	Bid string `json:"Dolar"`
}

var db *sql.DB

const awesomeApiURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {

	var err error
	db, err = sql.Open("sqlite3", "./bid.sqlite")
	if err != nil {
		log.Fatal("could open database connection", err)
	}
	defer db.Close()

	sqlQuery := `create table IF NOT EXISTS quotation (bid text);`
	_, err = db.Exec(sqlQuery)
	if err != nil {
		log.Fatal("Falha ao criar tabela", err)
	}

	http.HandleFunc("/cotacao", handler)
	fmt.Println("Server rodando :8080")
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", awesomeApiURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Não foi possível fazer a solicitação: %v", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Não foi possível obter a cotação: %v", err)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Falha na leitura: %v", err)
		return
	}

	var q USDBRLCotacaoResponse
	if err := json.Unmarshal(body, &q); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Falha na leitura: %v", err)
		return
	}

	bidOutput := CotacaoOutput{
		Bid: q.USDBRlCotacao.Bid,
	}

	ctx, cancel = context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Falha ao abrir transacao: %v", err)
		return
	}

	_, err = tx.ExecContext(ctx, "insert into quotation(bid) values(?)", bidOutput.Bid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Erro ao inserir %v", err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Falha no commit: %v", err)
		tx.Rollback()
		return
	}

	o, err := json.Marshal(bidOutput)
	if err != nil {
		fmt.Println("Error:", err)
	}

	w.Write(o)
}
