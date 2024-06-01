package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CotacaoOutput struct {
	Bid string `json:"Dolar"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusOK {
		return
	}

	fmt.Println(res.StatusCode)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Falha resposta: %v\n", err)
		return
	}

	var b CotacaoOutput
	if err := json.Unmarshal(body, &b); err != nil {
		fmt.Fprintf(os.Stderr, "Falha resposta: %v\n", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
	}

	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", b.Bid))
	fmt.Println("File created with success!")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing on file: %v\n", err)
	}

	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
