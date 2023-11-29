package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Address struct {
    Street string `json:"logradouro"`
}

func fetchAddress(ctx context.Context, cep string, ch chan<- string) {
    urls := []string{
        fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep),
        fmt.Sprintf("https://api.postmon.com.br/v1/cep/%s", cep),
    }

    for _, url := range urls {
        req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            ch <- fmt.Sprintf("Erro: %s", err)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            bodyBytes, _ := ioutil.ReadAll(resp.Body)
            bodyString := string(bodyBytes)
            ch <- fmt.Sprintf("Erro na requisição da API: %d, %s", resp.StatusCode, bodyString)
            return
        }

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
            return
        }
    }
}
func main() {
    startCEP := flag.String("start", "", "Início do intervalo de CEP")
    endCEP := flag.String("end", "", "Fim do intervalo de CEP")
    flag.Parse()

    start, err := strconv.Atoi(*startCEP)
    if err != nil {
        log.Fatalf("Erro ao converter CEP inicial: %s", err)
    }

    end, err := strconv.Atoi(*endCEP)
    if err != nil {
        log.Fatalf("Erro ao converter CEP final: %s", err)
    }

    for cep := start; cep <= end; cep++ {
        cepStr := fmt.Sprintf("%08d", cep)
        formattedCEP := cepStr[:5] + "-" + cepStr[5:]

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        ch := make(chan string)
        go fetchAddress(ctx, formattedCEP, ch)

        select {
        case res := <-ch:
            fmt.Println(res)
        case <-ctx.Done():
            fmt.Println("Tempo esgotado")
        }
    }
}