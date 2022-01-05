package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GAZIMAGomeDDD/billing-service/internal/handler"
	"github.com/GAZIMAGomeDDD/billing-service/internal/handler/mock"
	"github.com/GAZIMAGomeDDD/billing-service/internal/model"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/exchangeratesapi"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (t *handlerTestSuite) Test_getBalance() {
	s := mock.NewMockStore()
	s.On("GetBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(&model.User{
		ID:      "b91a95a4-078f-4afd-b11c-4850eb65e784",
		Balance: 99.99,
	}, nil)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	u := model.User{}
	json.NewDecoder(resp.Body).Decode(&u)
	t.Equal("b91a95a4-078f-4afd-b11c-4850eb65e784", u.ID)
	t.Equal(99.99, u.Balance)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_getBalanceWithCurrency() {
	s := mock.NewMockStore()
	s.On("GetBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(&model.User{
		ID:      "b91a95a4-078f-4afd-b11c-4850eb65e784",
		Balance: 75.0,
	}, nil)

	cr := new(mock.MockCurrency)
	cr.On("GetCurrencyRate", "USD").Return(75.0, nil)

	h := handler.New(s, cr)

	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getBalance?currency=USD", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	u := model.User{}
	json.NewDecoder(resp.Body).Decode(&u)
	t.Equal("b91a95a4-078f-4afd-b11c-4850eb65e784", u.ID)
	t.Equal(1.0, u.Balance)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_getBalanceErrUserNotFound() {
	s := mock.NewMockStore()
	s.On("GetBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(nil, postgres.ErrUserNotFound)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)

	t.Equal(http.StatusNotFound, resp.StatusCode)
}

func (t *handlerTestSuite) Test_getBalanceSomeServerError() {
	s := mock.NewMockStore()
	s.On("GetBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(nil, fmt.Errorf("some error happened"))

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)

	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (t *handlerTestSuite) Test_getBalanceSomeServerError2() {
	s := mock.NewMockStore()
	s.On("GetBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(&model.User{
		ID:      "b91a95a4-078f-4afd-b11c-4850eb65e784",
		Balance: 75.0,
	}, nil)

	cr := new(mock.MockCurrency)
	cr.On("GetCurrencyRate", "UDS").Return(0, fmt.Errorf("some error happened"))

	h := handler.New(s, cr)

	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getBalance?currency=UDS", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}
func (t *handlerTestSuite) Test_getBalanceErrWrongCurrency() {
	s := mock.NewMockStore()
	s.On("GetBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(&model.User{
		ID:      "b91a95a4-078f-4afd-b11c-4850eb65e784",
		Balance: 75.0,
	}, nil)

	cr := new(mock.MockCurrency)
	cr.On("GetCurrencyRate", "USD").Return(0, exchangeratesapi.ErrWrongCurrency)

	h := handler.New(s, cr)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getBalance?currency=USD", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_increaseBalance() {
	s := mock.NewMockStore()
	s.On("IncreaseOrDecreaseBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784", 175.22).Return(&model.User{
		ID:      "b91a95a4-078f-4afd-b11c-4850eb65e784",
		Balance: 175.22,
	}, nil)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"id":    "b91a95a4-078f-4afd-b11c-4850eb65e784",
			"money": 175.22,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/increaseOrDecreaseBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_DecreaseBalanceErrNotEnoughMoney() {
	s := mock.NewMockStore()
	s.On("IncreaseOrDecreaseBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784", -175000.22).Return(nil, postgres.ErrNotEnoughMoney)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"id":    "b91a95a4-078f-4afd-b11c-4850eb65e784",
			"money": -175000.22,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/increaseOrDecreaseBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_increaseBalanceSomeServerError() {
	s := mock.NewMockStore()
	s.On("IncreaseOrDecreaseBalance", "b91a95a4-078f-4afd-b11c-4850eb65e784", 175.22).Return(nil, fmt.Errorf("some error happened"))

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"id":    "b91a95a4-078f-4afd-b11c-4850eb65e784",
			"money": 175.22,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/increaseOrDecreaseBalance", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (t *handlerTestSuite) Test_moneyTransfer() {
	s := mock.NewMockStore()
	s.On("MoneyTransfer",
		"b91a95a4-078f-4afd-b11c-4850eb65e784",
		"b91a95a4-078f-4afd-b11c-4850eb65e785", 571.32).Return(nil)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "b91a95a4-078f-4afd-b11c-4850eb65e784",
			"from_id": "b91a95a4-078f-4afd-b11c-4850eb65e785",
			"money":   571.32,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/moneyTransfer", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_moneyTransferErrNotEnoughMoney() {
	s := mock.NewMockStore()
	s.On("MoneyTransfer",
		"b91a95a4-078f-4afd-b11c-4850eb65e784",
		"b91a95a4-078f-4afd-b11c-4850eb65e785", 571.32).Return(postgres.ErrNotEnoughMoney)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "b91a95a4-078f-4afd-b11c-4850eb65e784",
			"from_id": "b91a95a4-078f-4afd-b11c-4850eb65e785",
			"money":   571.32,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/moneyTransfer", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_moneyTransferErrUserNotFound() {
	s := mock.NewMockStore()
	s.On("MoneyTransfer",
		"b91a95a4-078f-4afd-b11c-4850eb65e784",
		"b91a95a4-078f-4afd-b11c-4850eb65e785", 571.32).Return(postgres.ErrUserNotFound)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "b91a95a4-078f-4afd-b11c-4850eb65e784",
			"from_id": "b91a95a4-078f-4afd-b11c-4850eb65e785",
			"money":   571.32,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/moneyTransfer", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusNotFound, resp.StatusCode)
}

func (t *handlerTestSuite) Test_moneyTransferSomeServerError() {
	s := mock.NewMockStore()
	s.On("MoneyTransfer",
		"b91a95a4-078f-4afd-b11c-4850eb65e784",
		"b91a95a4-078f-4afd-b11c-4850eb65e785", 571.32).Return(fmt.Errorf("some error happened"))

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(
		map[string]interface{}{
			"to_id":   "b91a95a4-078f-4afd-b11c-4850eb65e784",
			"from_id": "b91a95a4-078f-4afd-b11c-4850eb65e785",
			"money":   571.32,
		},
	)
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/moneyTransfer", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)
	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (t *handlerTestSuite) Test_getTransaction() {
	s := mock.NewMockStore()
	s.On("GetTransaction", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(&model.Transaction{
		ID: "b91a95a4-078f-4afd-b11c-4850eb65e784",
	}, nil)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getTransaction", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)

	t.Equal(http.StatusOK, resp.StatusCode)
}

func (t *handlerTestSuite) Test_getTransactionErrTransactionNotFound() {
	s := mock.NewMockStore()
	s.On("GetTransaction", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(nil, postgres.ErrTransactionNotFound)

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getTransaction", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)

	t.Equal(http.StatusNotFound, resp.StatusCode)
}

func (t *handlerTestSuite) Test_getTransactionSomeServerError() {
	s := mock.NewMockStore()
	s.On("GetTransaction", "b91a95a4-078f-4afd-b11c-4850eb65e784").Return(nil, fmt.Errorf("some error happened"))

	h := handler.New(s, nil)
	testSrv := httptest.NewServer(h.Init())
	c := testSrv.Client()

	body, err := json.Marshal(map[string]string{"id": "b91a95a4-078f-4afd-b11c-4850eb65e784"})
	t.Nil(err)
	req, err := http.NewRequest("POST", testSrv.URL+"/getTransaction", bytes.NewReader(body))
	t.Nil(err)
	resp, err := c.Do(req)
	t.Nil(err)

	t.Equal(http.StatusInternalServerError, resp.StatusCode)
}
