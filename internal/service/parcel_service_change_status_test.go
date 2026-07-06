package service

import (
	"context"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"errors"
	"testing"
)

func TestChangeStatus_InvalidTransition(t *testing.T) {
	parcelID := 1

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusCreated)}

	mockStatus := &mockStatusRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	svc := ParcelService{
		parcelReader: mockReader,
		statusRepo:   mockStatus,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("invalid transition", func(t *testing.T) {
		err := svc.ChangeStatus(context.Background(), parcelID, domain.StatusInTransit, "loc", 1)
		if !errors.Is(err, ErrInvalidStatusTransition) {
			t.Fatalf("got %v, want ErrInvalidStatusTransition", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if mockStatus.getByIDCalled {
			t.Fatal("StatusRepository.GetStatusID should not be called")
		}

		if mockTxManager.doCalled {
			t.Fatal("TransactionManager.Do should not be called")
		}
	})
}

func TestChangeStatus_ParcelNotFound(t *testing.T) {
	parcelID := 2

	mockReader := &mockParcelRepo{
		getByIDResult: testParcel(domain.StatusCreated),
		getByIDErr:    repository.ErrParcelNotFound,
	}

	mockStatus := &mockStatusRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	svc := ParcelService{
		parcelReader: mockReader,
		statusRepo:   mockStatus,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("parcel not found error", func(t *testing.T) {
		err := svc.ChangeStatus(context.Background(), parcelID, domain.StatusPurchased, "loc", 1)
		if !errors.Is(err, repository.ErrParcelNotFound) {
			t.Fatalf("got %v, want ErrParcelNotFound", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if mockStatus.getByIDCalled {
			t.Fatal("StatusRepository.GetStatusID should not be called")
		}

		if mockTxManager.doCalled {
			t.Fatal("TransactionManager.Do should not be called")
		}
	})
}
