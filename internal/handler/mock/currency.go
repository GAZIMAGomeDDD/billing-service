package mock

import "github.com/stretchr/testify/mock"

type MockCurrency struct {
	mock.Mock
}

func (m *MockCurrency) GetCurrencyRate(currency string) (float64, error) {
	args := m.Called(currency)

	arg0 := args.Get(0)

	if arg0 == 0 {
		return 0, args.Error(1)
	}

	return arg0.(float64), args.Error(1)
}
