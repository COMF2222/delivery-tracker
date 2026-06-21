package repository

import (
	"delivery-tracker/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuditRepository struct {
	db *sqlx.DB
}

func NewAuditRepository(db *sqlx.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log *domain.AuditLog) error {
	query := `INSERT INTO audit_logs (
		user_id, 
		action, 
		old_value, 
		new_value, 
		entity_type, 
		entity_id
		) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := r.db.QueryRow(
		query,
		log.UserID,
		log.Action,
		log.OldValue,
		log.NewValue,
		log.EntityType,
		log.EntityID,
	).Scan(&log.ID)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

func (r *AuditRepository) CreateTx(tx *sqlx.Tx, log *domain.AuditLog) error {
	query := `INSERT INTO audit_logs (
		user_id, 
		action, 
		old_value, 
		new_value, 
		entity_type, 
		entity_id
		) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := tx.QueryRow(
		query,
		log.UserID,
		log.Action,
		log.OldValue,
		log.NewValue,
		log.EntityType,
		log.EntityID,
	).Scan(&log.ID)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}
