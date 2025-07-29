# CEP Lookup Service

Este é um serviço simples de consulta de CEP (Código de Endereçamento Postal) brasileiro que utiliza duas APIs diferentes (Brasil API e ViaCEP API) para obter informações de endereço a partir de um CEP.

## Funcionalidades

- Consulta concorrente em duas APIs de CEP diferentes
- Retorna o resultado da primeira API que responder
- Timeout de 1 segundo caso nenhuma API responda a tempo
- Exibe informações detalhadas do endereço

## Requisitos

- Go 1.22 ou superior

## Instalação

1. Clone o repositório:

```bash
git clone https://seu-repositorio/StudyGo.git
cd StudyGo
```

2. Não há dependências externas para instalar, pois o projeto utiliza apenas a biblioteca padrão do Go.

## Como Executar

Para executar a aplicação com o CEP padrão (84350000):

```bash
go run .
```

ou

```bash
go run main.go responseCEP.go
```

### Modificando o CEP

Para consultar um CEP diferente, você precisa modificar a linha 12 no arquivo `main.go`:

```go
cep := "84350000" // Altere para o CEP desejado
```

Após a modificação, execute a aplicação novamente.

## Exemplo de Saída

Quando a consulta é bem-sucedida, você verá uma saída semelhante a esta:

```
Resposta via Brasil API
CEP: 84350-000
Logradouro: 
Complemento: 
Bairro: 
Localidade: Ortigueira
UF: PR
IBGE: 
DDD: 
```

ou

```
Resposta via ViaCEP API
CEP: 84350-000
Logradouro: 
Complemento: 
Bairro: 
Localidade: Ortigueira
UF: PR
IBGE: 4117305
DDD: 42
```

## Como Testar

### Testando com Diferentes CEPs

1. Modifique o CEP no arquivo `main.go` (linha 12)
2. Execute a aplicação para verificar se o novo CEP retorna resultados válidos

### Testando o Timeout

O aplicativo tem um timeout configurado para 1 segundo. Se ambas as APIs demorarem mais que isso para responder, você verá a mensagem:

```
Timeout: Both APIs took too long to respond
```

### Testando com CEP Inválido

Para testar o comportamento com um CEP inválido, modifique o CEP para um valor que não existe (por exemplo, "00000000") e execute a aplicação.

## APIs Utilizadas

### Brasil API

- URL: `https://brasilapi.com.br/api/cep/v1/{cep}`
- Documentação: [Brasil API - CEP](https://brasilapi.com.br/docs)

### ViaCEP API

- URL: `http://viacep.com.br/ws/{cep}/json/`
- Documentação: [ViaCEP](https://viacep.com.br/)

## Estrutura do Projeto

- `main.go`: Contém a lógica principal da aplicação, incluindo as funções para consultar as APIs
- `responseCEP.go`: Define a estrutura de dados para armazenar as respostas das APIs
- `go.mod`: Arquivo de configuração do módulo Go