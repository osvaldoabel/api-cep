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

type brasilApiProvider struct {
	HttpClient *http.Client
	Name       string
	Url        string
}

func NewBrasilApiProvider() PostalCodeProvider {
	return &brasilApiProvider{
		Name: "BrasilAPI",
		Url:  "",
	}
}

// Struct para o JSON da BrasilAPI
type brasilApiAddress struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func (ba *brasilApiAddress) GetAddressFields() (string, error) {
	return utils.AddressToJson[brasilApiAddress](*ba)
}

func (p *brasilApiProvider) GetAddress(ctx context.Context) (Address, error) {
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

		var address brasilApiAddress
		if err = json.Unmarshal(body, &address); err != nil {
			log.Default().Printf("error while trying to unpack http response body. %v", err)
			return nil, err
		}

		return &address, nil
	}

}
