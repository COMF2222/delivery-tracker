package service

import (
	"context"
	"delivery-tracker/internal/domain"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type mockParcelCache struct {
	getResult *domain.ParcelDetails
	getErr    error
	setErr    error
	deleteErr error

	getCalled    bool
	setCalled    bool
	deleteCalled bool
}

type mockParcelRepo struct {
	getByTrackCalled bool
	getByIDCalled    bool
	archiveCalled    bool
	updateCalled     bool

	getResult     *domain.Parcel
	getByIDResult *domain.Parcel

	getErr     error
	getByIDErr error
	updateErr  error
	archiveErr error
}

type mockParcelPhotoRepo struct {
	getByIdCalled bool
	createCalled  bool

	getResult []domain.ParcelPhoto
	getErr    error
	createErr error
}

type mockParcelHistoryRepo struct {
	getByIdCalled bool
	createCalled  bool

	getResult []domain.ParcelStatusHistory
	getErr    error
	createErr error
}

type mockStatusRepo struct {
	getStatusIDCalled bool

	getResult int
	getErr    error
}

type mockAuditRepo struct {
	createCalled bool

	createErr error
}

type mockTransactionManager struct {
	doCalled bool
	doErr    error
}

func (m *mockParcelCache) GetByTrack(ctx context.Context, trackNumber string) (*domain.ParcelDetails, error) {
	m.getCalled = true
	return m.getResult, m.getErr
}

func (m *mockParcelCache) SetByTrack(ctx context.Context, trackNumber string, parcel *domain.ParcelDetails, ttl time.Duration) error {
	m.setCalled = true
	return m.setErr
}

func (m *mockParcelCache) DeleteByTrack(ctx context.Context, trackNumber string) error {
	m.deleteCalled = true
	return m.deleteErr
}

func (m *mockParcelRepo) GetByTrackNumber(trackNumber string) (*domain.Parcel, error) {
	m.getByTrackCalled = true
	return m.getResult, m.getErr
}

func (m *mockParcelRepo) GetByID(id int) (*domain.Parcel, error) {
	m.getByIDCalled = true
	return m.getByIDResult, m.getByIDErr
}

func (m *mockParcelPhotoRepo) GetByParcelID(parcelID int) ([]domain.ParcelPhoto, error) {
	m.getByIdCalled = true
	return m.getResult, m.getErr
}

func (m *mockParcelPhotoRepo) Create(photo *domain.ParcelPhoto) error {
	m.createCalled = true
	return m.createErr
}

func (m *mockParcelHistoryRepo) GetByParcelID(parcelID int) ([]domain.ParcelStatusHistory, error) {
	m.getByIdCalled = true
	return m.getResult, m.getErr
}

func (m *mockParcelHistoryRepo) CreateTx(tx *sqlx.Tx, history *domain.ParcelStatusHistory, oldStatusID int, newStatusID int) error {
	m.createCalled = true
	return m.createErr
}

func (m *mockStatusRepo) GetStatusID(status domain.Status) (int, error) {
	m.getStatusIDCalled = true
	return m.getResult, m.getErr
}

func (m *mockAuditRepo) CreateTx(tx *sqlx.Tx, log *domain.AuditLog) error {
	m.createCalled = true
	return m.createErr
}

func (m *mockParcelRepo) UpdateStatusTx(tx *sqlx.Tx, parcelID, statusID int, location string) error {
	m.updateCalled = true
	return m.updateErr
}

func (m *mockParcelRepo) CreateParcel(parcel *domain.Parcel, statusID int) error {
	return errors.New("create parcel should not be called")
}

func (m *mockParcelRepo) ArchiveTx(tx *sqlx.Tx, parcelID int) error {
	m.archiveCalled = true
	return m.archiveErr
}

func (m *mockTransactionManager) Do(fn func(tx *sqlx.Tx) error) error {
	m.doCalled = true

	if m.doErr != nil {
		return m.doErr
	}

	return fn(nil)
}
