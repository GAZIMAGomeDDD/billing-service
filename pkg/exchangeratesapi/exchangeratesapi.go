package exchangeratesapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

func GetCurrencyExchangeRate(currency string) (float64, error) {
	var body struct {
		Rates map[string]float64 `json:"rates"`
	}

	currency = strings.ToUpper(currency)
	url := "http://api.exchangeratesapi.io/v1/latest?access_key=207e0d99dc8df832c4921e5af54e56e4&format=1"

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
		return 0, err
	}

	rates := body.Rates
	if _, ok := rates[currency]; !ok {
		return 0, ErrWrongCurrency
	}

	rate := rates["RUB"] / rates[currency]

	return rate, nil
}
