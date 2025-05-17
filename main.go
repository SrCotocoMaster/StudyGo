package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	cep := "84350000" // Using a valid CEP for testing

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

// displayAddress prints the address information and the source API
func displayAddress(resp responseCEP) {
	fmt.Println("Response from:", resp.Source)
	fmt.Println("CEP:", resp.Cep)
	fmt.Println("Logradouro:", resp.Logradouro)
	fmt.Println("Complemento:", resp.Complemento)
	fmt.Println("Bairro:", resp.Bairro)
	fmt.Println("Localidade:", resp.Localidade)
	fmt.Println("UF:", resp.UF)
	fmt.Println("IBGE:", resp.IBGE)
	fmt.Println("DDD:", resp.DDD)
}

// brasilApiResponse represents the response structure from Brasil API
type brasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

// brasilApi makes a request to the Brasil API
func brasilApi(c chan responseCEP, cep string) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 900 * time.Millisecond, // Setting timeout slightly below 1 second
	}

	// Make the request
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making request to Brasil API:", err)
		return // Don't send anything through the channel if there's an error
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Brasil API returned non-OK status: %d\n", resp.StatusCode)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Brasil API response:", err)
		return
	}

	// Check if the response contains an error message (Brasil API returns a different structure for errors)
	var errorResp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
		fmt.Println("Brasil API returned error:", errorResp.Message)
		return
	}

	// Parse the JSON response using the Brasil API specific struct
	var brasilResp brasilApiResponse
	if err := json.Unmarshal(body, &brasilResp); err != nil {
		fmt.Println("Error parsing Brasil API response:", err)
		return
	}

	// Check if we got a valid CEP
	if brasilResp.Cep == "" {
		fmt.Println("Brasil API returned empty CEP")
		return
	}

	// Map the Brasil API response to our common responseCEP struct
	cepResp := responseCEP{
		Cep:        brasilResp.Cep,
		Logradouro: brasilResp.Street,
		Bairro:     brasilResp.Neighborhood,
		Localidade: brasilResp.City,
		UF:         brasilResp.State,
		Source:     "Brasil API",
	}

	// Send through the channel
	c <- cepResp
}

// viaCep makes a request to the ViaCEP API
func viaCep(c chan responseCEP, cep string) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 900 * time.Millisecond, // Setting timeout slightly below 1 second
	}

	// Make the request
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error making request to ViaCEP API:", err)
		return // Don't send anything through the channel if there's an error
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ViaCEP API returned non-OK status: %d\n", resp.StatusCode)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading ViaCEP API response:", err)
		return
	}

	// Check for ViaCEP error response (it returns a JSON with "erro": true if CEP is invalid)
	if string(body) == `{"erro": true}` || string(body) == `{"erro":true}` {
		fmt.Println("ViaCEP API returned error: CEP not found")
		return
	}

	// Parse the JSON response
	var cepResp responseCEP
	if err := json.Unmarshal(body, &cepResp); err != nil {
		fmt.Println("Error parsing ViaCEP API response:", err)
		return
	}

	// Set the source and send through the channel
	cepResp.Source = "ViaCEP API"
	c <- cepResp
}
