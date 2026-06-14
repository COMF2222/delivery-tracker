package repository

import (
	"delivery-tracker/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type ParcelStatusHistoryRepository struct {
	db *sqlx.DB
}

func NewParcelStatusHistoryRepository(db *sqlx.DB) *ParcelStatusHistoryRepository {
	return &ParcelStatusHistoryRepository{db: db}
}

func (r *ParcelStatusHistoryRepository) GetByParcelID(parcelID int) ([]domain.ParcelStatusHistory, error) {
	parcelStatusHistory := make([]domain.ParcelStatusHistory, 0)
	query := `
	SELECT
		p.id AS id,
		p.parcel_id AS parcel_id,
		old_s.status AS old_status,
		new_s.status AS new_status,
		p.location AS location,
		p.changed_by AS changed_by,
		p.created_at AS created_at
	FROM parcel_status_history p
	LEFT JOIN statuses old_s ON p.old_status = old_s.id
	JOIN statuses new_s ON p.new_status = new_s.id
	WHERE p.parcel_id = $1`

	rows, err := r.db.Queryx(query, parcelID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()
	for rows.Next() {
		var history domain.ParcelStatusHistory
		if err = rows.StructScan(&history); err != nil {
			return nil, fmt.Errorf("failed to get parcel status history by parcel ID(%d), %w", parcelID, err)
		}
		parcelStatusHistory = append(parcelStatusHistory, history)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate parcel status history rows: %w", err)
	}
	return parcelStatusHistory, nil
}

func (r *ParcelStatusHistoryRepository) CreateInitialHistory(
	history *domain.ParcelStatusHistory,
	oldStatusID *int,
	newStatusID int) error {
	query := `INSERT INTO parcel_status_history(
		parcel_id, 
		old_status, 
		new_status, 
		location, 
		changed_by
		) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err := r.db.QueryRow(
		query,
		history.ParcelID,
		oldStatusID,
		newStatusID,
		history.Location,
		history.ChangedBy,
	).Scan(&history.ID)

	if err != nil {
		return fmt.Errorf("failed to create initial parcel history: %w", err)
	}

	return nil
}

func (r *ParcelStatusHistoryRepository) Create(
	history *domain.ParcelStatusHistory,
	oldStatusID int,
	newStatusID int) error {
	query := `INSERT INTO parcel_status_history(
		parcel_id, 
		old_status, 
		new_status, 
		location, 
		changed_by
		) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err := r.db.QueryRow(
		query,
		history.ParcelID,
		oldStatusID,
		newStatusID,
		history.Location,
		history.ChangedBy,
	).Scan(&history.ID)

	if err != nil {
		return fmt.Errorf("failed to create parcel history: %w", err)
	}

	return nil
}

func (r *ParcelStatusHistoryRepository) CreateTx(
	tx *sqlx.Tx,
	history *domain.ParcelStatusHistory,
	oldStatusID int,
	newStatusID int) error {
	query := `INSERT INTO parcel_status_history(
		parcel_id, 
		old_status, 
		new_status, 
		location, 
		changed_by
		) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err := tx.QueryRow(
		query,
		history.ParcelID,
		oldStatusID,
		newStatusID,
		history.Location,
		history.ChangedBy,
	).Scan(&history.ID)

	if err != nil {
		return fmt.Errorf("failed to create parcel history: %w", err)
	}

	return nil
}
