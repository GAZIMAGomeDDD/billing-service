package postgres

import (
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*Store, error) {
	s := new(Store)
	s.db = db

	if err := s.initSchema(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) initSchema() error {
	if _, err := s.db.Exec(schema); err != nil {
		return err
	}

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}
