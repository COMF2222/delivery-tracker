package repository

import (
	"database/sql"
	"delivery-tracker/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ParcelRepository struct {
	db *sqlx.DB
}

func NewParcelRepository(db *sqlx.DB) *ParcelRepository {
	return &ParcelRepository{db: db}
}

func (r *ParcelRepository) CreateParcel(parcel *domain.Parcel, statusID int) error {
	query := `INSERT INTO parcels(
		track_number, 
		item_name, 
		recipient_name, 
		recipient_phone, 
		recipient_address, 
		current_status, 
		current_location, 
		is_archived	
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := r.db.QueryRow(
		query,
		parcel.TrackNumber,
		parcel.ItemName,
		parcel.RecipientName,
		parcel.RecipientPhone,
		parcel.RecipientAddress,
		statusID,
		parcel.CurrentLocation,
		parcel.IsArchived,
	).Scan(&parcel.ID)

	if err != nil {
		if isTrackNumberConflict(err) {
			return ErrTrackNumberAlreadyExists
		}

		return fmt.Errorf("failed to create parcel: %w", err)
	}

	return nil
}

func isTrackNumberConflict(err error) bool {
	var pqErr *pq.Error

	return errors.As(err, &pqErr) &&
		pqErr.Code == "23505" &&
		pqErr.Constraint == "parcels_track_number_key"
}

func (r *ParcelRepository) GetByTrackNumber(trackNumber string) (*domain.Parcel, error) {
	parcel := domain.Parcel{}
	query := `SELECT
		p.id,
		p.track_number,
		p.item_name,
		p.recipient_name,
		s.status,
		p.current_location 
	FROM parcels p
    JOIN statuses s ON p.current_status = s.id
    WHERE track_number = $1`

	if err := r.db.QueryRow(query, trackNumber).Scan(&parcel.ID, &parcel.TrackNumber, &parcel.ItemName, &parcel.RecipientName, &parcel.CurrentStatus, &parcel.CurrentLocation); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrParcelNotFound
		}

		return nil, fmt.Errorf("failed to get parcel by track number: %w", err)
	}

	return &parcel, nil
}

func (r *ParcelRepository) GetByID(id int) (*domain.Parcel, error) {
	parcel := domain.Parcel{}
	query := `SELECT 
		p.id, 
		p.track_number, 
		p.item_name, 
		p.recipient_name, 
		p.recipient_phone, 
		p.recipient_address,
		s.status,
		p.current_location,
		p.is_archived,
		p.created_at,
		p.updated_at
	FROM parcels p
	JOIN statuses s ON p.current_status = s.id
	WHERE p.id = $1`

	if err := r.db.QueryRow(query, id).Scan(
		&parcel.ID,
		&parcel.TrackNumber,
		&parcel.ItemName,
		&parcel.RecipientName,
		&parcel.RecipientPhone,
		&parcel.RecipientAddress,
		&parcel.CurrentStatus,
		&parcel.CurrentLocation,
		&parcel.IsArchived,
		&parcel.CreatedAt,
		&parcel.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrParcelNotFound
		}

		return nil, fmt.Errorf("failed to get parcel by ID: %w", err)
	}

	return &parcel, nil
}

func (r *ParcelRepository) UpdateStatus(parcelID, statusID int, location string) error {
	query := `UPDATE parcels 
	SET
	    current_status = $1,
  		current_location = $2,
  		updated_at = NOW()
	WHERE id = $3`

	result, err := r.db.Exec(query, statusID, location, parcelID)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrParcelNotFound
	}

	return nil
}

func (r *ParcelRepository) UpdateStatusTx(tx *sqlx.Tx, parcelID, statusID int, location string) error {
	query := `UPDATE parcels 
	SET
	    current_status = $1,
  		current_location = $2,
  		updated_at = NOW()
	WHERE id = $3`

	result, err := tx.Exec(query, statusID, location, parcelID)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrParcelNotFound
	}

	return nil
}
