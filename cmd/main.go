package main

import (
	"context"
	"fmt"
	"log"

	entity "github.com/osvaldoabel/api-cep/internal/entity"
)

func main() {
	brasilApiCtx, brasilApiCancel := context.WithCancel(context.Background())
	viaCepCtx, viaCepCancel := context.WithCancel(context.Background())

	ch1 := make(chan string)
	ch2 := make(chan string)

	brasilApi := entity.NewBrasilApiProvider()
	viaCep := entity.NewViaCepProvider()

	go func() {
		brasilApiAddress, err := brasilApi.GetAddress(brasilApiCtx)
		if err != nil {
			log.Default().Printf("error while trying to consult an Address through BrasilAPI. %v", err)
			return
		}

		result, err := brasilApiAddress.GetAddressFields()
		if err != nil {
			log.Default().Printf("error while trying to get address fields. %v", err)
			return
		}

		ch1 <- result
	}()

	go func() {
		viaCepAddress, err := viaCep.GetAddress(viaCepCtx)
		if err != nil {
			log.Default().Printf("error while trying to consult an Address through ViaCep. %v", err)
		}

		result, err := viaCepAddress.GetAddressFields()
		if err != nil {
			log.Default().Printf("error while trying to get address fields. %v", err)
			return
		}

		ch2 <- result
	}()

	select {
	case viaCepResult := <-ch1:
		fmt.Println("ViaCep ==============>")
		fmt.Println(viaCepResult)
		brasilApiCancel()
	case brasilApiResult := <-ch2:
		fmt.Println("BrasilAPI ==============>")
		fmt.Println(brasilApiResult)
		viaCepCancel()
	}

}
