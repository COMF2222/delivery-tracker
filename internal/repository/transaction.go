package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) Do(fn func(tx *sqlx.Tx) error) error {
	tx, err := tm.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	fnErr := fn(tx)

	if fnErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("rollback transaction: %w, fnErr: %w", rollbackErr, fnErr)
		}
		return fnErr
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
