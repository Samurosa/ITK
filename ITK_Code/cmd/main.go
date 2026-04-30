package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

type PaymentProcessor interface {
	ProcessPayment(amount float64) error
}

type Sber struct {
	APIKey string
}

type Tbank struct {
	APIKey string
}

type Alfabank struct {
	APIKey string
}

var (
	ErrInvalidAmount = errors.New("некорректная сумма платежа")

	ErrProviderUnavailable = errors.New("провайдер недоступен")

	ErrInvalidAPIkey = errors.New("неизвестный провайдер")
)

func (b Sber) ProcessPayment(a float64) error {
	return process(a, b.APIKey)
}

func (b Alfabank) ProcessPayment(a float64) error {
	return process(a, b.APIKey)
}

func (b Tbank) ProcessPayment(a float64) error {
	return process(a, b.APIKey)
}

func process(a float64, api string) error {
	if api == "" {
		return ErrInvalidAPIkey
	}

	if a < 0 {
		return ErrInvalidAmount
	}

	err := checkProviderUnavailable()
	if err != nil {
		return err
	}

	return nil
}

func checkProviderUnavailable() error {
	r := rand.Float64()
	if r < 0.2 {
		return ErrProviderUnavailable
	}
	return nil
}
func main() {

	processor := []PaymentProcessor{
		Tbank{APIKey: "ASDASDqw123sad1"},
		Alfabank{APIKey: "123213asdasd1213"},
		Sber{APIKey: "qweqwesadasd12321"},
		Sber{APIKey: ""},
	}

	for _, p := range processor {
		err := p.ProcessPayment(111)
		if err != nil {
			switch t := p.(type) {
			case Sber:
				fmt.Println("Sber error:", t.APIKey, err)
			case Alfabank:
				fmt.Println("Alfabank error:", t.APIKey, err)
			case Tbank:
				fmt.Println("Tbank error:", t.APIKey, err)
			}
		}
	}

}
