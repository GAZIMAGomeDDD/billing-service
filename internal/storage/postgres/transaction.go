package postgres

import (
	"database/sql"
	"fmt"

	"github.com/GAZIMAGomeDDD/billing-service/internal/model"
)

func (s *Store) GetTransaction(tid string) (*model.Transaction, error) {
	var transaction model.Transaction

	err := s.db.Get(&transaction, sqlGetTransaction, tid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	return &transaction, nil
}

func (s *Store) ListOfTransactions(uid, sort string, limit, page int) ([]model.Transaction, error) {
	var sql string

	switch sort {
	case "date_asc":
		sql = fmt.Sprintf(sqlListOfTransactions, "ORDER BY created_at ASC")
	case "money_asc":
		sql = fmt.Sprintf(sqlListOfTransactions, "ORDER BY money ASC")
	case "date_desc":
		sql = fmt.Sprintf(sqlListOfTransactions, "ORDER BY create_at DESC")
	case "money_desc":
		sql = fmt.Sprintf(sqlListOfTransactions, "ORDER BY money DESC")
	default:
		sql = fmt.Sprintf(sqlListOfTransactions, "")
	}

	transactions := []model.Transaction{}

	err := s.db.Select(&transactions, sql, uid, limit, page)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Store) writeTransaction(money float64, to_id, from_id, method *string, tx *sql.Tx) error {
	if _, err := tx.Exec(sqlWriteTransaction, to_id, from_id, money, method); err != nil {
		return err
	}

	return nil
}
