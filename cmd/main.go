package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/reinaldosaraiva/desafio-go-multithread/configs"
)

type Address struct {
    Street string `json:"logradouro"`
}

func fetchAddress(ctx context.Context, url string, ch chan<- string) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        ch <- fmt.Sprintf("Erro: %s", err)
        return
    }
    defer resp.Body.Close()

    var responseData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        ch <- fmt.Sprintf("Erro na decodificação: %s", err)
        return
    }

    if erro, ok := responseData["erro"]; ok && erro.(bool) {
        ch <- fmt.Sprintf("CEP não encontrado: %s", url)
        return
    }

    if street, ok := responseData["logradouro"]; ok {
        ch <- fmt.Sprintf("Endereço: %s, API: %s", street, url)
    }
}

func main() {
    startCEP := flag.String("start", "", "Início do intervalo de CEP")
    endCEP := flag.String("end", "", "Fim do intervalo de CEP")
    flag.Parse()

    if *startCEP == "" || *endCEP == "" {
        log.Fatal("Início e fim do intervalo de CEP são parâmetros obrigatórios")
    }

    start, err := strconv.Atoi(*startCEP)
    if err != nil {
        log.Fatalf("Erro ao converter CEP inicial: %s", err)
    }

    end, err := strconv.Atoi(*endCEP)
    if err != nil {
        log.Fatalf("Erro ao converter CEP final: %s", err)
    }

    config, err := configs.LoadConfig(".")
    if err != nil {
        log.Fatalf("Não foi possível carregar a configuração: %s", err)
    }

    for cep := start; cep <= end; cep++ {
        cepStr := strconv.Itoa(cep)
        formattedCEP := cepStr[:5] + "-" + cepStr[5:]

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        ch := make(chan string, 2)
        go fetchAddress(ctx, config.APIUrl1+formattedCEP+".json", ch)
        go fetchAddress(ctx, config.APIUrl2+formattedCEP+"/json/", ch)

        select {
        case result := <-ch:
            fmt.Println(result)
        case <-ctx.Done():
            fmt.Println("Timeout para o CEP:", formattedCEP)
        }
    }
}
