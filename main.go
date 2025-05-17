package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	cep := "84350000"

	c1 := make(chan responseCEP)
	c2 := make(chan responseCEP)

	go brasilApi(c1, cep)
	go viaCep(c2, cep)

	select {
	case resp := <-c1:
		displayAddress(resp)

	case resp := <-c2:
		displayAddress(resp)

	case <-time.After(time.Second * 1):
		fmt.Println("Timeout: Both APIs took too long to respond")
	}
}

func displayAddress(resp responseCEP) {
	fmt.Println("Resposta via ", resp.Source)
	fmt.Println("CEP:", resp.Cep)
	fmt.Println("Logradouro:", resp.Logradouro)
	fmt.Println("Complemento:", resp.Complemento)
	fmt.Println("Bairro:", resp.Bairro)
	fmt.Println("Localidade:", resp.Localidade)
	fmt.Println("UF:", resp.UF)
	fmt.Println("IBGE:", resp.IBGE)
	fmt.Println("DDD:", resp.DDD)
}

type brasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func brasilApi(c chan responseCEP, cep string) {
	client := &http.Client{
		Timeout: 900 * time.Millisecond,
	}

	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error ao fazer a request para Brasil API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Brasil API returned non-OK status: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler Brasil API response:", err)
		return
	}

	var errorResp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
		fmt.Println("Brasil API returned error:", errorResp.Message)
		return
	}

	var brasilResp brasilApiResponse
	if err := json.Unmarshal(body, &brasilResp); err != nil {
		fmt.Println("Error parsing Brasil API response:", err)
		return
	}

	if brasilResp.Cep == "" {
		fmt.Println("Brasil API returned empty CEP")
		return
	}

	cepResp := responseCEP{
		Cep:        brasilResp.Cep,
		Logradouro: brasilResp.Street,
		Bairro:     brasilResp.Neighborhood,
		Localidade: brasilResp.City,
		UF:         brasilResp.State,
		Source:     "Brasil API",
	}

	c <- cepResp
}

func viaCep(c chan responseCEP, cep string) {
	client := &http.Client{
		Timeout: 900 * time.Millisecond,
	}

	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error ao fazer a request para ViaCEP API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ViaCEP API returned non-OK status: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error ao ler ViaCEP API response:", err)
		return
	}

	if string(body) == `{"erro": true}` || string(body) == `{"erro":true}` {
		fmt.Println("ViaCEP API retorn error: CEP nao encontrado")
		return
	}

	var cepResp responseCEP
	if err := json.Unmarshal(body, &cepResp); err != nil {
		fmt.Println("Error parsing ViaCEP API response:", err)
		return
	}

	cepResp.Source = "ViaCEP API"
	c <- cepResp
}
