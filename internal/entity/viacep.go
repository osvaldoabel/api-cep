package entity

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/osvaldoabel/api-cep/pkg/utils"
)

type viaCepProvider struct {
	HttpClient *http.Client
	Name       string
	Url        string
}

func NewViaCepProvider(httpClient *http.Client) PostalCodeProvider {
	return &viaCepProvider{
		HttpClient: httpClient,
		Name:       "ViaCep",
		Url:        "",
	}
}

// Struct para o JSON da ViaCEP
type ViaCepAddress struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Service     string `json:"service"`
}

func (va *ViaCepAddress) GetAddressFields() (string, error) {
	result, err := utils.AddressToJson[ViaCepAddress](*va)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *viaCepProvider) GetAddress(ctx context.Context) (Address, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.TimeoutLimit)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.Url, nil)
	if err != nil {
		log.Default().Printf("error while trying to create an http request. %v", err)
		return nil, err
	}

	response := make(chan *http.Response)
	go func() {
		resp, err := p.HttpClient.Do(req)
		if err != nil {
			log.Default().Printf("error while trying to execute request. %v", err)
			return
		}
		defer resp.Body.Close()

		response <- resp
	}()

	select {
	case <-ctx.Done():
		return nil, errors.New("request has been finished or timed out")
	case result := <-response:

		body, err := io.ReadAll(result.Body)
		if err != nil {
			return nil, errors.New("request has been finished or timed out")
		}

		var address ViaCepAddress
		if err = json.Unmarshal(body, &address); err != nil {
			log.Default().Printf("error while trying to unpack http response body. %v", err)
			return nil, err
		}

		return &address, nil
	}

}
