package model

import "time"

type Transaction struct {
	ID       string    `json:"id" db:"id"`
	ToID     *string   `json:"to_id,omitempty" db:"to_id"`
	FromID   *string   `json:"from_id,omitempty" db:"from_id"`
	Money    float64   `json:"money" db:"money"`
	Method   string    `json:"method" db:"method"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}
