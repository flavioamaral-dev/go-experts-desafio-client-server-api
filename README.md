<h1 align="center">
   <a href="#"> Client Server API </a>
</h1>

<h3 align="center">
    Desafio curso FullCycle - Client Server Api em Go lang
</h3>

## Funcionamento

O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
 
O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.
 
Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
 
O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
 
Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
 
O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
 
O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.
---



#### Running Project

```bash

# Clonar
$ git clone https://github.com/flavioamaral-dev/go-experts-desafio-client-server-api.git


# Acessar o diretorio server
$ cd server

# executar run main.go
$ go run main.go

# Servidor rodando
Server is running on :8080


# Acessar o diretório client
$ cd client 

# executar run main.go
$ go run main.go



# os dados são gravados em
cotacao.txt


