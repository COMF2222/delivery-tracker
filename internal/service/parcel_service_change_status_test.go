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

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusCreated, false)}

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

		if mockStatus.getStatusIDCalled {
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
		getByIDResult: testParcel(domain.StatusCreated, false),
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

		if mockStatus.getStatusIDCalled {
			t.Fatal("StatusRepository.GetStatusID should not be called")
		}

		if mockTxManager.doCalled {
			t.Fatal("TransactionManager.Do should not be called")
		}
	})
}

func TestChangeStatus_StatusRepoError(t *testing.T) {
	parcelID := 1

	statusErr := errors.New("status err")

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusCreated, false)}

	mockStatus := &mockStatusRepo{getErr: statusErr}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	svc := ParcelService{
		parcelReader: mockReader,
		statusRepo:   mockStatus,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("status repo error", func(t *testing.T) {
		err := svc.ChangeStatus(context.Background(), parcelID, domain.StatusPurchased, "loc", 1)
		if !errors.Is(err, statusErr) {
			t.Fatalf("got %v, want statusErr", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if mockTxManager.doCalled {
			t.Fatal("TransactionManager.Do should not be called")
		}
	})
}

func TestChangeStatus_TransactionError(t *testing.T) {
	parcelID := 1

	transactionErr := errors.New("transaction err")

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusCreated, false)}

	mockStatus := &mockStatusRepo{getResult: 2}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{doErr: transactionErr}

	svc := ParcelService{
		parcelReader: mockReader,
		statusRepo:   mockStatus,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("transaction error", func(t *testing.T) {
		err := svc.ChangeStatus(context.Background(), parcelID, domain.StatusPurchased, "loc", 1)
		if !errors.Is(err, transactionErr) {
			t.Fatalf("got %v, want transactionErr", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}
	})
}

func TestChangeStatus_TransactionSuccess(t *testing.T) {
	parcelID := 1

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusCreated, false)}

	mockStatus := &mockStatusRepo{getResult: 2}

	mockCache := &mockParcelCache{}

	mockAudit := &mockAuditRepo{}

	mockHistory := &mockParcelHistoryRepo{}

	mockTxManager := &mockTransactionManager{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockReader,
		statusRepo:   mockStatus,
		historyRepo:  mockHistory,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
		auditRepo:    mockAudit,
	}

	t.Run("success", func(t *testing.T) {
		err := svc.ChangeStatus(context.Background(), parcelID, domain.StatusPurchased, "loc", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockCache.deleteCalled {
			t.Fatal("expected ParcelCache.DeleteByTrack to be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockReader.updateCalled {
			t.Fatal("expected ParcelWriter.UpdateStatusTx to be called")
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if !mockAudit.createCalled {
			t.Fatal("expected AuditRepo.CreateTx to be called")
		}

		if !mockHistory.createCalled {
			t.Fatal("expected HistoryRepository.CreateTx to be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}
	})
}

func TestChangeStatus_CacheDeleteError(t *testing.T) {
	parcelID := 1

	deleteCacheError := errors.New("delete cache by track error")

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusCreated, false)}

	mockStatus := &mockStatusRepo{getResult: 2}

	mockCache := &mockParcelCache{deleteErr: deleteCacheError}

	mockAudit := &mockAuditRepo{}

	mockHistory := &mockParcelHistoryRepo{}

	mockTxManager := &mockTransactionManager{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockReader,
		statusRepo:   mockStatus,
		historyRepo:  mockHistory,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
		auditRepo:    mockAudit,
	}

	t.Run("cache delete error does not fail status change", func(t *testing.T) {
		err := svc.ChangeStatus(context.Background(), parcelID, domain.StatusPurchased, "loc", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockCache.deleteCalled {
			t.Fatal("expected ParcelCache.DeleteByTrack to be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockReader.updateCalled {
			t.Fatal("expected ParcelWriter.UpdateStatusTx to be called")
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if !mockAudit.createCalled {
			t.Fatal("expected AuditRepo.CreateTx to be called")
		}

		if !mockHistory.createCalled {
			t.Fatal("expected HistoryRepository.CreateTx to be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}
	})
}
