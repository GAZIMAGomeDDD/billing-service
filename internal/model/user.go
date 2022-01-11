package model

type User struct {
	ID      string  `json:"id" db:"id"`
	Balance float64 `json:"balance" db:"balance"`
}

type GetBalanceQuery struct {
	UserID string `json:"id" example:"b91a95a4-078f-4afd-b11c-4850eb65e784"`
}

type GetBalanceResponse struct {
	UserID  string  `json:"id" example:"b91a95a4-078f-4afd-b11c-4850eb65e784"`
	Balance float64 `json:"balance" example:"99.99"`
}

type ChangeBalanceQuery struct {
	UserID string  `json:"id" example:"b91a95a4-078f-4afd-b11c-4850eb65e784"`
	Money  float64 `json:"money" example:"99.99"`
}

type ChangeBalanceResponse struct {
	UserID  string  `json:"id" example:"b91a95a4-078f-4afd-b11c-4850eb65e784"`
	Balance float64 `json:"balance" example:"99.99"`
}

type MoneyTransfer struct {
	ToID   string  `json:"to_id" example:"b91a95a4-078f-4afd-b11c-4850eb65e784"`
	FromID string  `json:"from_id" example:"b81a95a4-078f-5dfd-b11c-4850eb35e785"`
	Money  float64 `json:"money" example:"99.99"`
}
