package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	_ "github.com/GAZIMAGomeDDD/billing-service/docs"

	"github.com/GAZIMAGomeDDD/billing-service/internal/model"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/exchangeratesapi"
	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/neilotoole/errgroup"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Billing Service API
// @version 1.0
type Handler struct {
	mux    *chi.Mux
	logger *logrus.Logger
	store  SQLStorage
	cr     Currency
}

type SQLStorage interface {
	GetBalance(string) (*model.User, error)
	MoneyTransfer(string, string, float64) error
	ChangeBalance(string, float64) (*model.User, error)
	GetTransaction(string) (*model.Transaction, error)
	ListOfTransactions(string, string, int, int) ([]model.Transaction, error)
}

type Currency interface {
	GetCurrencyRate(string) (float64, error)
}

func New(s SQLStorage, cr Currency) *Handler {
	return &Handler{
		mux:    chi.NewRouter(),
		store:  s,
		logger: logrus.New(),
		cr:     cr,
	}
}

// @Produce json
// @Param getBalance body model.GetBalanceQuery true "--"
// @Param currency query string false "currency"
// @Success 200 {object} model.GetBalanceResponse
// @Failure 404 {string} http.Error "user not found"
// @Failure 500 {string} http.Error "Internal server error"
// @Router /getBalance [post]
func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {
	var body model.GetBalanceQuery
	currency := r.URL.Query().Get("currency")

	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	eg, _ := errgroup.WithContext(r.Context())

	var user *model.User
	var err error
	eg.Go(func() error {
		user, err = h.store.GetBalance(body.UserID)
		if err != nil {
			return err
		}

		return nil
	})

	var rate float64
	if currency != "" {
		eg.Go(func() error {
			rate, err = h.cr.GetCurrencyRate(currency)
			if err != nil {
				return err
			}

			return nil
		})

	}

	if err = eg.Wait(); err != nil {
		switch err {
		case postgres.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case exchangeratesapi.ErrWrongCurrency:
			http.Error(w, err.Error(), http.StatusOK)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		h.logger.Error(err)
		return
	}

	if currency != "" {
		user.Balance = math.Round((user.Balance/rate)*100) / 100
	}

	json.NewEncoder(w).Encode(user)
}

// @Produce json
// @Param changeBalance body model.ChangeBalanceQuery true "--"
// @Success 200 {object} model.ChangeBalanceResponse
// @Failure 200 {string} http.Error "not enough money"
// @Failure 500  {string} http.Error "Internal server error"
// @Router /changeBalance [post]
func (h *Handler) changeBalance(w http.ResponseWriter, r *http.Request) {
	var body model.ChangeBalanceQuery

	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	user, err := h.store.ChangeBalance(body.UserID, body.Money)
	if err != nil {
		switch err {
		case postgres.ErrNotEnoughMoney:
			http.Error(w, err.Error(), http.StatusOK)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		h.logger.Error(err)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// @Produce json
// @Param moneyTransfer body model.MoneyTransfer true "--"
// @Failure 404 {string} http.Error "user not found"
// @Failure 200 {string} http.Error "not enough money"
// @Failure 500  {string} http.Error "Internal server error"
// @Router /moneyTransfer [post]
func (h *Handler) moneyTransfer(w http.ResponseWriter, r *http.Request) {
	var body model.MoneyTransfer

	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	if err := h.store.MoneyTransfer(body.ToID, body.FromID, body.Money); err != nil {
		switch err {
		case postgres.ErrNotEnoughMoney:
			http.Error(w, err.Error(), http.StatusOK)
		case postgres.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		h.logger.Error(err)
		return
	}
}

// @Produce json
// @Param listOfTransactions body model.ListOfTransactionsQuery true "--"
// @Param page query int true "page"
// @Param sort query string false "sort"
// @Success 200 {array} model.Transaction
// @Router /listOfTransactions [post]
func (h *Handler) listOfTransactions(w http.ResponseWriter, r *http.Request) {
	var body model.ListOfTransactionsQuery

	pageString := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageString)

	sort := r.URL.Query().Get("sort")

	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	transactions, err := h.store.ListOfTransactions(body.UserID, sort, body.Limit, page)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

// @Produce json
// @Param getTransaction body model.GetTransaction true "--"
// @Success 200 {object} model.Transaction
// @Router /getTransaction [post]
func (h *Handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	var t model.GetTransaction

	json.NewDecoder(r.Body).Decode(&t)
	defer r.Body.Close()

	transaction, err := h.store.GetTransaction(t.ID)
	if err != nil {
		switch err {
		case postgres.ErrTransactionNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		h.logger.Error(err)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) Init() *chi.Mux {
	h.mux.Use(logger.Logger("router", h.logger))
	h.mux.Post("/getBalance", h.getBalance)
	h.mux.Post("/changeBalance", h.changeBalance)
	h.mux.Post("/moneyTransfer", h.moneyTransfer)
	h.mux.Post("/listOfTransactions", h.listOfTransactions)
	h.mux.Post("/getTransaction", h.getTransaction)

	h.mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json")))

	return h.mux
}
