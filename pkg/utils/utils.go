package utils

import (
	"encoding/json"
	"time"
)

const (
	TimeoutLimit = time.Second * 1
	ViaCepUrl    = "https://viacep.com.br/ws/01153000/json"
	BrasilAPIUrl = "https://brasilapi.com.br/api/cep/v1/01153000"
)

func AddressToJson[T any](address T) (string, error) {
	addressSTr, err := json.Marshal(*&address)
	if err != nil {
		return "", err
	}

	return string(addressSTr), nil
}
