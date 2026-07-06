package service

import (
	"context"
	"delivery-tracker/internal/domain"
	"github.com/jmoiron/sqlx"
	"time"
)

type ParcelReader interface {
	GetByTrackNumber(trackNumber string) (*domain.Parcel, error)
	GetByID(id int) (*domain.Parcel, error)
}

type ParcelWriter interface {
	CreateParcel(parcel *domain.Parcel, statusID int) error
	UpdateStatusTx(tx *sqlx.Tx, parcelID, statusID int, location string) error
	ArchiveTx(tx *sqlx.Tx, parcelID int) error
}

type ParcelLister interface {
	List(limit, offset int) ([]domain.Parcel, error)
	Count() (int, error)
	ListByStatus(status domain.Status, limit, offset int) ([]domain.Parcel, error)
	CountByStatus(status domain.Status) (int, error)
}

type StatusRepository interface {
	GetStatusID(status domain.Status) (int, error)
}

type ParcelPhotoRepository interface {
	GetByParcelID(parcelID int) ([]domain.ParcelPhoto, error)
	Create(photo *domain.ParcelPhoto) error
}

type ParcelStatusHistoryRepository interface {
	GetByParcelID(parcelID int) ([]domain.ParcelStatusHistory, error)
	CreateTx(tx *sqlx.Tx, history *domain.ParcelStatusHistory, oldStatusID int, newStatusID int) error
}

type AuditRepository interface {
	CreateTx(tx *sqlx.Tx, log *domain.AuditLog) error
}

type ParcelCache interface {
	GetByTrack(ctx context.Context, trackNumber string) (*domain.ParcelDetails, error)
	SetByTrack(ctx context.Context, trackNumber string, parcel *domain.ParcelDetails, ttl time.Duration) error
	DeleteByTrack(ctx context.Context, trackNumber string) error
}

type TransactionManager interface {
	Do(fn func(tx *sqlx.Tx) error) error
}
