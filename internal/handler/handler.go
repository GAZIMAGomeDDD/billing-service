package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/GAZIMAGomeDDD/billing-service/internal/model"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/exchangeratesapi"
	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	mux    *chi.Mux
	logger *logrus.Logger
	store  SQLStorage
	cr     Currency
}

type SQLStorage interface {
	GetBalance(string) (*model.User, error)
	MoneyTransfer(string, string, float64) error
	IncreaseOrDecreaseBalance(string, float64) (*model.User, error)
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

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {
	var u model.User
	currency := r.URL.Query().Get("currency")

	json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()

	user, err := h.store.GetBalance(u.ID)
	if err != nil {
		switch err {
		case postgres.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		h.logger.Error(err)
		return
	}

	if currency != "" {
		rate, err := h.cr.GetCurrencyRate(currency)
		if err != nil {
			switch err {
			case exchangeratesapi.ErrWrongCurrency:
				http.Error(w, err.Error(), http.StatusOK)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			h.logger.Error(err)
			return
		}

		user.Balance = math.Round((user.Balance/rate)*100) / 100
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) increaseOrDecreaseBalance(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID    string  `json:"id"`
		Money float64 `json:"money"`
	}

	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	user, err := h.store.IncreaseOrDecreaseBalance(body.ID, body.Money)
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

func (h *Handler) moneyTransfer(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ToID   string  `json:"to_id"`
		FromID string  `json:"from_id"`
		Money  float64 `json:"money"`
	}

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

func (h *Handler) listOfTransactions(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID string `json:"user_id"`
		Limit  int    `json:"limit"`
	}

	pageString := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageString)

	json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()

	transactions, err := h.store.ListOfTransactions(body.UserID, "date", body.Limit, page)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func (h *Handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	var t model.Transaction

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
	h.mux.Post("/increaseOrDecreaseBalance", h.increaseOrDecreaseBalance)
	h.mux.Post("/moneyTransfer", h.moneyTransfer)
	h.mux.Post("/listOfTransactions", h.listOfTransactions)
	h.mux.Post("/getTransaction", h.getTransaction)

	return h.mux
}
