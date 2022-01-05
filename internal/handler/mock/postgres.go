package mock

import (
	"github.com/GAZIMAGomeDDD/billing-service/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func NewMockStore() *MockStore {
	store := new(MockStore)

	return store
}

func (m *MockStore) GetBalance(uid string) (*model.User, error) {
	args := m.Called(uid)

	arg0 := args.Get(0)

	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*model.User), args.Error(1)
}

func (m *MockStore) MoneyTransfer(to_id, from_id string, money float64) error {
	args := m.Called(to_id, from_id, money)

	return args.Error(0)
}

func (m *MockStore) IncreaseOrDecreaseBalance(uid string, money float64) (*model.User, error) {
	args := m.Called(uid, money)

	arg0 := args.Get(0)

	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*model.User), args.Error(1)
}

func (m *MockStore) GetTransaction(tid string) (*model.Transaction, error) {
	args := m.Called(tid)

	arg0 := args.Get(0)

	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*model.Transaction), args.Error(1)
}

func (m *MockStore) ListOfTransactions(uid, sort string, limit, page int) ([]model.Transaction, error) {
	args := m.Called(uid, sort, limit, page)

	arg0 := args.Get(0)

	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.([]model.Transaction), args.Error(1)
}
