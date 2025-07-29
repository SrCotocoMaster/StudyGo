# Sistema de Cotação do Dólar

Este é um sistema cliente-servidor em Go para consulta de cotação do dólar americano (USD) em relação ao real brasileiro (BRL). O sistema é composto por um servidor HTTP que consome uma API externa de cotação e um cliente que solicita a cotação ao servidor.

## Funcionalidades

- **Servidor HTTP** (`server.go`):
  - Endpoint `/cotacao` na porta 8080
  - Consome a API da AwesomeAPI para obter cotação USD-BRL
  - Armazena cada cotação em banco de dados SQLite
  - Timeout de 200ms para chamada da API externa
  - Timeout de 100ms para persistência no banco de dados
  - Retorna apenas o valor "bid" da cotação em formato JSON

- **Cliente HTTP** (`client.go`):
  - Faz requisição para o servidor local
  - Timeout de 300ms para receber resposta do servidor
  - Salva a cotação recebida em arquivo `cotacao.txt`
  - Formato do arquivo: "Dólar: {valor}"

## Requisitos

- Go 1.23 ou superior
- Conexão com a internet (para API externa)

## Instalação

1. Clone o repositório:

```bash
git clone https://seu-repositorio/StudyGo.git
cd StudyGo
```

2. Instale as dependências:

```bash
go mod tidy
```

## Como Testar

### Passo 1: Executar o Servidor

Abra um terminal e execute o servidor:

```bash
go run server.go
```

Você deve ver a mensagem:
```
Server running on port 8080...
```

O servidor ficará rodando e aguardando requisições. Deixe este terminal aberto.

### Passo 2: Executar o Cliente

Abra um **segundo terminal** (mantendo o servidor rodando no primeiro) e execute o cliente:

```bash
go run client.go
```

### Resultado Esperado

Se tudo funcionar corretamente, você verá no terminal do cliente:

```
Raw server response: {"bid":"5.1234"}
Parsed bid value: 5.1234
Cotação do dólar salva com sucesso!
Dólar: 5.1234
```

E no terminal do servidor, você verá logs similares a:

```
API Response: {"USDBRL":{"code":"USD","codein":"BRL","name":"Dólar Americano/Real Brasileiro","high":"5.1500","low":"5.1000","varBid":"0.0234","pctChange":"0.46","bid":"5.1234","ask":"5.1240","timestamp":"1640995200","create_date":"2021-12-31 18:00:00"}}
Parsed Bid: 5.1234
```

### Verificar Arquivo Gerado

Após a execução do cliente, verifique se foi criado o arquivo `cotacao.txt`:

```bash
type cotacao.txt
```

O conteúdo deve ser algo como:
```
Dólar: 5.1234
```

### Verificar Banco de Dados

O servidor cria automaticamente um banco SQLite (`cotacao.db`) e armazena cada cotação. Para verificar:

```bash
# Se você tiver SQLite instalado
sqlite3 cotacao.db "SELECT * FROM cotacoes;"
```

## Testando Cenários de Erro

### Teste de Timeout do Cliente

Para testar o timeout do cliente (300ms), pare o servidor e execute o cliente:

```bash
go run client.go
```

Você deve ver um erro de conexão recusada.

### Teste de Timeout da API Externa

O servidor tem timeout de 200ms para a API externa. Se a API estiver lenta, você verá logs de erro no servidor.

### Teste de Timeout do Banco de Dados

O servidor tem timeout de 100ms para operações de banco. Em caso de problemas, você verá logs de erro, mas a cotação ainda será retornada ao cliente.

## Estrutura do Projeto

- `server.go`: Servidor HTTP que consome API externa e persiste dados
- `client.go`: Cliente que solicita cotação e salva em arquivo
- `cotacao.db`: Banco de dados SQLite (criado automaticamente)
- `cotacao.txt`: Arquivo com a última cotação (criado pelo cliente)
- `go.mod`: Configuração do módulo Go com dependências

## API Externa Utilizada

- **AwesomeAPI**: `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- Retorna cotação atual do dólar americano em reais brasileiros

## Troubleshooting

### Erro "connection refused"
- Certifique-se de que o servidor está rodando antes de executar o cliente
- Verifique se a porta 8080 não está sendo usada por outro processo

### Erro "context deadline exceeded"
- Verifique sua conexão com a internet
- A API externa pode estar temporariamente indisponível

### Erro de permissão no arquivo
- Certifique-se de ter permissão de escrita no diretório atual
- Verifique se o arquivo `cotacao.txt` não está sendo usado por outro programa