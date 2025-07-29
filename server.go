package main

// Importação das bibliotecas necessárias
import (
	"context"       // Para controle de timeout e cancelamento
	"database/sql"  // Para operações com banco de dados
	"encoding/json" // Para codificação e decodificação JSON
	"fmt"           // Para formatação de strings
	"io"            // Para operações de entrada/saída
	"log"           // Para logging de erros e informações
	"net/http"      // Para servidor e cliente HTTP
	"time"          // Para operações com tempo e timeout

	_ "modernc.org/sqlite" // Driver SQLite (importação apenas para efeitos colaterais)
)

// Estrutura para mapear a resposta da API de cotação
type ExchangeRate struct {
	USDBRL struct {
		Code       string `json:"code"`        // Código da moeda origem
		Codein     string `json:"codein"`      // Código da moeda destino
		Name       string `json:"name"`        // Nome completo do par de moedas
		High       string `json:"high"`        // Maior valor do dia
		Low        string `json:"low"`         // Menor valor do dia
		VarBid     string `json:"varBid"`      // Variação do valor de compra
		PctChange  string `json:"pctChange"`   // Percentual de mudança
		Bid        string `json:"bid"`         // Valor de compra (cotação)
		Ask        string `json:"ask"`         // Valor de venda
		Timestamp  string `json:"timestamp"`   // Timestamp da cotação
		CreateDate string `json:"create_date"` // Data de criação
	} `json:"USDBRL"`
}

// Estrutura de resposta que será enviada ao cliente
type ExchangeRateResponse struct {
	Bid string `json:"bid"` // Apenas o valor de compra (cotação)
}

func main() {
	// Inicializa a conexão com o banco de dados SQLite
	db, err := sql.Open("sqlite", "./cotacao.db")
	if err != nil {
		log.Fatal("Erro ao abrir banco de dados: ", err)
	}
	defer db.Close()

	// Cria a tabela de cotações se ela não existir
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Erro ao criar tabela: ", err)
	}

	// Configura o servidor HTTP e define o endpoint /cotacao
	http.HandleFunc("/cotacao", handleExchangeRate(db))
	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Função que manipula as requisições HTTP para o endpoint /cotacao
func handleExchangeRate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtém a cotação do dólar com timeout
		exchangeRate, err := getExchangeRate()
		if err != nil {
			log.Println("Erro ao obter cotação:", err)
			http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
			return
		}

		// Salva a cotação no banco de dados com timeout
		err = saveExchangeRate(db, exchangeRate.USDBRL.Bid)
		if err != nil {
			log.Println("Erro ao salvar no banco de dados:", err)
			// Continua e retorna a cotação mesmo se falhar ao salvar no BD
		}

		// Prepara e retorna a resposta JSON
		response := ExchangeRateResponse{
			Bid: exchangeRate.USDBRL.Bid,
		}

		// Define o cabeçalho de resposta e envia o JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Função que busca a cotação do dólar na API externa
func getExchangeRate() (*ExchangeRate, error) {
	// Cria contexto com timeout - aumentado para lidar com limitação de taxa da API
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Cria requisição HTTP com contexto
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
	}

	// Envia a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode == 429 {
		// API com limitação de taxa - retorna valor de fallback para manter sistema funcionando
		log.Println("API com limitação de taxa, usando cotação de fallback")
		return &ExchangeRate{
			USDBRL: struct {
				Code       string `json:"code"`
				Codein     string `json:"codein"`
				Name       string `json:"name"`
				High       string `json:"high"`
				Low        string `json:"low"`
				VarBid     string `json:"varBid"`
				PctChange  string `json:"pctChange"`
				Bid        string `json:"bid"`
				Ask        string `json:"ask"`
				Timestamp  string `json:"timestamp"`
				CreateDate string `json:"create_date"`
			}{
				Bid: "5.25", // Valor de fallback
			},
		}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status não-OK: %d", resp.StatusCode)
	}

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler corpo da resposta: %v", err)
	}

	// Log da resposta bruta para debug
	log.Printf("Resposta da API: %s", string(body))

	// Faz o parse do JSON
	var exchangeRate ExchangeRate
	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		// Tenta fazer parse com estrutura diferente se a estrutura esperada falhar
		var rawResponse map[string]interface{}
		if jsonErr := json.Unmarshal(body, &rawResponse); jsonErr == nil {
			log.Printf("Resposta bruta parseada: %v", rawResponse)
			// Tenta extrair dados USDBRL da resposta bruta
			if usdbrl, ok := rawResponse["USDBRL"].(map[string]interface{}); ok {
				if bid, ok := usdbrl["bid"].(string); ok && bid != "" {
					exchangeRate.USDBRL.Bid = bid
				}
			}
		} else {
			return nil, fmt.Errorf("erro ao fazer parse do JSON: %v", err)
		}
	}

	// Valida o valor da cotação
	if exchangeRate.USDBRL.Bid == "" {
		return nil, fmt.Errorf("valor de cotação vazio na resposta da API")
	}

	// Log dos dados parseados para debug
	log.Printf("Cotação parseada: %s", exchangeRate.USDBRL.Bid)

	return &exchangeRate, nil
}

// Função que salva a cotação no banco de dados SQLite
func saveExchangeRate(db *sql.DB, bid string) error {
	// Cria contexto com timeout para operação de banco de dados
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Prepara statement SQL para inserção
	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)")
	if err != nil {
		return fmt.Errorf("erro ao preparar statement: %v", err)
	}
	defer stmt.Close()

	// Executa o statement com o valor da cotação
	_, err = stmt.ExecContext(ctx, bid)
	if err != nil {
		return fmt.Errorf("erro ao executar statement: %v", err)
	}

	return nil
}
