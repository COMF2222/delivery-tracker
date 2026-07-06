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

	getResult     *domain.Parcel
	getByIDResult *domain.Parcel

	getErr     error
	getByIDErr error
}

type mockParcelPhotoRepo struct {
	getByIdCalled bool

	getResult []domain.ParcelPhoto
	getErr    error
}

type mockParcelHistoryRepo struct {
	getByIdCalled bool

	getResult []domain.ParcelStatusHistory
	getErr    error
}

type mockStatusRepo struct {
	getByIDCalled bool

	getResult int
	getErr    error
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
	return errors.New("create should not be called")
}

func (m *mockParcelHistoryRepo) GetByParcelID(parcelID int) ([]domain.ParcelStatusHistory, error) {
	m.getByIdCalled = true
	return m.getResult, m.getErr
}

func (m *mockParcelHistoryRepo) CreateTx(tx *sqlx.Tx, history *domain.ParcelStatusHistory, oldStatusID int, newStatusID int) error {
	return errors.New("createTx should not be called")
}

func (m *mockStatusRepo) GetStatusID(status domain.Status) (int, error) {
	m.getByIDCalled = true
	return m.getResult, m.getErr
}

func (m *mockTransactionManager) Do(fn func(tx *sqlx.Tx) error) error {
	m.doCalled = true
	return m.doErr
}
