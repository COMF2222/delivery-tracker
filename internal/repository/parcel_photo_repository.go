package repository

import (
	"delivery-tracker/internal/domain"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type ParcelPhotoRepository struct {
	db *sqlx.DB
}

func NewParcelPhotoRepository(db *sqlx.DB) *ParcelPhotoRepository {
	return &ParcelPhotoRepository{db: db}
}

func (r *ParcelPhotoRepository) GetByParcelID(parcelID int) ([]domain.ParcelPhoto, error) {
	parcelPhotos := make([]domain.ParcelPhoto, 0)

	query := `SELECT id, parcel_id, file_path, created_at FROM parcel_photos WHERE parcel_id = $1`

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
		var photo domain.ParcelPhoto
		if err = rows.StructScan(&photo); err != nil {
			return nil, fmt.Errorf("failed to get parcel photo by parcel id(%d): %w", parcelID, err)
		}
		parcelPhotos = append(parcelPhotos, photo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate parcel photo rows: %w", err)
	}

	return parcelPhotos, nil
}

func (r *ParcelPhotoRepository) Create(photo *domain.ParcelPhoto) error {
	query := `INSERT INTO parcel_photos(parcel_id, file_path) VALUES ($1, $2) RETURNING id`

	if err := r.db.QueryRow(query, photo.ParcelID, photo.FilePath).Scan(&photo.ID); err != nil {
		return fmt.Errorf("failed to create parcel photo: %w", err)
	}

	return nil
}
