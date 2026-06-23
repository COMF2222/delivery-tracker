package repository

import (
	"delivery-tracker/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type StatusRepository struct {
	db *sqlx.DB
}

func NewStatusRepository(db *sqlx.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) GetStatusID(status domain.Status) (int, error) {
	var id int
	query := `SELECT id FROM statuses WHERE status = $1`

	if err := r.db.QueryRow(query, status).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to get status id: %w", err)
	}

	return id, nil
}
