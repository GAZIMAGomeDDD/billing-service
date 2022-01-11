package postgres

import (
	"database/sql"
	"strings"

	"github.com/GAZIMAGomeDDD/billing-service/internal/model"
)

func (s *Store) createUser(uid string, tx *sql.Tx) error {
	if _, err := tx.Exec(sqlCreateUser, uid); err != nil {
		return err
	}

	return nil
}

func (s *Store) ChangeBalance(uid string, money float64) (*model.User, error) {
	var user model.User

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRow(`SELECT id FROM users WHERE id = $1`, uid).Scan(&user.ID); err != nil {
		switch err {
		case sql.ErrNoRows:
			err = s.createUser(uid, tx)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	row := tx.QueryRow(sqlUpdateBalance, uid, money)
	if err = row.Scan(&user.ID, &user.Balance); err != nil {
		if strings.Contains(err.Error(), "users_balance_check") {
			return nil, ErrNotEnoughMoney
		}

		return nil, err
	}

	var method string

	switch {
	case money > 0:
		method = "Increase balance"

		if err = s.writeTransaction(money, &uid, nil, &method, tx); err != nil {
			return nil, err
		}
	case money < 0:
		method = "Decrease balance"

		if err = s.writeTransaction(money, nil, &uid, &method, tx); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetBalance(uid string) (*model.User, error) {
	var user model.User

	err := s.db.Get(&user, sqlGetBalance, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) MoneyTransfer(to_id, from_id string, money float64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var temp interface{}

	if err = tx.QueryRow(sqlUpdateBalance, to_id, money).Scan(&temp, &temp); err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}

		return err
	}

	if err = tx.QueryRow(sqlUpdateBalance, from_id, -money).Scan(&temp, &temp); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return ErrUserNotFound
		case strings.Contains(err.Error(), "users_balance_check"):
			return ErrNotEnoughMoney
		default:
			return err
		}
	}

	method := "Money transfer"

	if err = s.writeTransaction(money, &to_id, &from_id, &method, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
