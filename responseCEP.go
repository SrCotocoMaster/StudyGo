package main

// responseCEP represents the combined response from either API
// The struct includes fields from both Brasil API and ViaCEP API
type responseCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"` // cidade in Brasil API
	UF          string `json:"uf"`         // estado in Brasil API
	IBGE        string `json:"ibge"`
	DDD         string `json:"ddd"`
	Source      string // Indicates which API provided the response
}
