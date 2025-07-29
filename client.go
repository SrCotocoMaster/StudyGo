package main

// Importação das bibliotecas necessárias para o cliente
import (
	"context"       // Para controle de timeout e cancelamento
	"encoding/json" // Para decodificação de JSON
	"fmt"           // Para formatação de strings e saída
	"io"            // Para operações de entrada/saída
	"log"           // Para logging de erros
	"net/http"      // Para cliente HTTP
	"os"            // Para operações com arquivos
	"time"          // Para operações com tempo e timeout
)

// Estrutura para mapear a resposta do servidor
type CotacaoResponse struct {
	Bid string `json:"bid"` // Valor da cotação do dólar
}

func main() {
	// Cria contexto com timeout de 300ms para a requisição
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Cria requisição HTTP com contexto para o servidor local
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar requisição: ", err)
	}

	// Envia a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Erro ao fazer requisição: ", err)
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Servidor retornou status não-OK: %d", resp.StatusCode)
	}

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erro ao ler corpo da resposta: ", err)
	}

	// Imprime resposta bruta para debug
	fmt.Println("Resposta bruta do servidor:", string(body))

	// Faz o parse do JSON
	var cotacao CotacaoResponse
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		// Tenta fazer parse com estrutura diferente se a estrutura esperada falhar
		var rawResponse map[string]interface{}
		if jsonErr := json.Unmarshal(body, &rawResponse); jsonErr == nil {
			fmt.Println("Resposta bruta parseada:", rawResponse)
			// Tenta extrair o valor bid da resposta bruta
			if bid, ok := rawResponse["bid"].(string); ok && bid != "" {
				cotacao.Bid = bid
			}
		} else {
			log.Fatal("Erro ao fazer parse do JSON: ", err)
		}
	}

	// Valida o valor da cotação
	if cotacao.Bid == "" {
		log.Fatal("Valor de cotação vazio na resposta do servidor")
	}

	// Imprime valor da cotação parseada para debug
	fmt.Println("Valor da cotação parseada:", cotacao.Bid)

	// Salva a cotação no arquivo
	err = saveToFile(cotacao.Bid)
	if err != nil {
		log.Fatal("Erro ao salvar no arquivo: ", err)
	}

	fmt.Println("Cotação do dólar salva com sucesso!")
	fmt.Printf("Dólar: %s\n", cotacao.Bid)
}

// Função que salva a cotação do dólar em arquivo de texto
func saveToFile(bid string) error {
	// Cria o arquivo cotacao.txt
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	// Escreve a cotação no arquivo no formato "Dólar: {valor}"
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", bid))
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo: %v", err)
	}

	return nil
}
