# Simulador de Busca de CEP com Go

## Descrição do Projeto
Este projeto em Go realiza consultas simultâneas a diferentes APIs de CEP para buscar endereços com base em um intervalo de CEPs fornecido. Utiliza goroutines para fazer requisições paralelas e tratar respostas de acordo com diferentes cenários, incluindo limites de requisições e verificações de erros.

## Requisitos
- Go (versão 1.8)
- Acesso à internet para realizar as consultas às APIs


## Uso
Navegue até a pasta do projeto e execute o programa com os seguintes comandos:

Para consultar um intervalo de CEPs:
```bash
go run main.go -start=64014050 -end=64014100
```

Substitua `64014050` e `64014100` pelos CEPs inicial e final do intervalo desejado.

## Estrutura do Projeto
- `main.go`: Arquivo principal que executa as consultas aos CEPs.
- `configs/config.go`: Carrega configurações a partir de um arquivo `.env`.

