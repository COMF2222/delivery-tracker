package service

import (
	"context"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"errors"
	"testing"
)

func TestArchive_ParcelNotFound(t *testing.T) {
	parcelID := 2

	mockReader := &mockParcelRepo{
		getByIDErr: repository.ErrParcelNotFound,
	}

	mockWriter := &mockParcelRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	mockAudit := &mockAuditRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		parcelCache:  mockCache,
		auditRepo:    mockAudit,
		txManager:    mockTxManager,
	}

	t.Run("parcel not found error", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if !errors.Is(err, repository.ErrParcelNotFound) {
			t.Fatalf("got %v, want ErrParcelNotFound", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if mockWriter.archiveCalled {
			t.Fatal("ParcelWriter.ArchiveTx should not be called")
		}

		if mockTxManager.doCalled {
			t.Fatal("TransactionManager.Do should not be called")
		}

		if mockAudit.createCalled {
			t.Fatal("AuditRepository.CreateTx should not be called")
		}
	})
}

func TestArchive_ParcelAlreadyArchived(t *testing.T) {
	parcelID := 1

	mockReader := &mockParcelRepo{
		getByIDResult: testParcel(domain.StatusCreated, true),
	}

	mockWriter := &mockParcelRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	mockAudit := &mockAuditRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		parcelCache:  mockCache,
		auditRepo:    mockAudit,
		txManager:    mockTxManager,
	}

	t.Run("parcel already archived", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if !errors.Is(err, ErrParcelAlreadyArchived) {
			t.Fatalf("got %v, want ErrParcelAlreadyArchived", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if mockWriter.archiveCalled {
			t.Fatal("ParcelWriter.ArchiveTx should not be called")
		}

		if mockTxManager.doCalled {
			t.Fatal("TransactionManager.Do should not be called")
		}

		if mockAudit.createCalled {
			t.Fatal("AuditRepository.CreateTx should not be called")
		}
	})
}

func TestArchive_ParcelNotDelivered(t *testing.T) {
	parcelID := 1

	mockReader := &mockParcelRepo{
		getByIDResult: testParcel(domain.StatusCreated, false),
	}

	mockWriter := &mockParcelRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	mockAudit := &mockAuditRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		parcelCache:  mockCache,
		auditRepo:    mockAudit,
		txManager:    mockTxManager,
	}

	t.Run("parcel not delivered", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if !errors.Is(err, ErrParcelNotDelivered) {
			t.Fatalf("got %v, want ErrParcelNotDelivered", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if mockWriter.archiveCalled {
			t.Fatal("ParcelWriter.ArchiveTx should not be called")
		}

		if mockTxManager.doCalled {
			t.Fatal("TransactionManager.Do should not be called")
		}

		if mockAudit.createCalled {
			t.Fatal("AuditRepository.CreateTx should not be called")
		}
	})
}

func TestArchive_TransactionError(t *testing.T) {
	parcelID := 1

	transactionErr := errors.New("transaction err")

	mockReader := &mockParcelRepo{
		getByIDResult: testParcel(domain.StatusDelivered, false),
	}

	mockWriter := &mockParcelRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{doErr: transactionErr}

	mockAudit := &mockAuditRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		parcelCache:  mockCache,
		auditRepo:    mockAudit,
		txManager:    mockTxManager,
	}

	t.Run("transaction error", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if !errors.Is(err, transactionErr) {
			t.Fatalf("got %v, want transactionErr", err)
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if mockWriter.archiveCalled {
			t.Fatal("ParcelWriter.ArchiveTx should not be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}

		if mockAudit.createCalled {
			t.Fatal("AuditRepository.CreateTx should not be called")
		}
	})
}

func TestArchive_ArchiveTxError(t *testing.T) {
	parcelID := 1

	archiveErr := errors.New("archive error")

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusDelivered, false)}

	mockWriter := &mockParcelRepo{archiveErr: archiveErr}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	mockAudit := &mockAuditRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		auditRepo:    mockAudit,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("archive error", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if !errors.Is(err, archiveErr) {
			t.Fatalf("got %v, want archiveErr", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}

		if !mockWriter.archiveCalled {
			t.Fatal("expected ParcelWriter.ArchiveTx to be called")
		}

		if mockAudit.createCalled {
			t.Fatal("AuditRepository.CreateTx should not be called")
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

	})
}

func TestArchive_AuditError(t *testing.T) {
	parcelID := 1

	auditErr := errors.New("audit error")

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusDelivered, false)}

	mockWriter := &mockParcelRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	mockAudit := &mockAuditRepo{createErr: auditErr}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		auditRepo:    mockAudit,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("audit create error", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if !errors.Is(err, auditErr) {
			t.Fatalf("got %v, want auditErr", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}

		if !mockWriter.archiveCalled {
			t.Fatal("expected ParcelWriter.ArchiveTx to be called")
		}

		if !mockAudit.createCalled {
			t.Fatal("expected AuditRepository.CreateTx to be called")
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}

	})
}

func TestArchive_CacheDeleteError(t *testing.T) {
	parcelID := 1

	deleteCacheError := errors.New("delete cache by track error")

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusDelivered, false)}

	mockWriter := &mockParcelRepo{}

	mockCache := &mockParcelCache{deleteErr: deleteCacheError}

	mockTxManager := &mockTransactionManager{}

	mockAudit := &mockAuditRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		auditRepo:    mockAudit,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("cache delete error", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}

		if !mockWriter.archiveCalled {
			t.Fatal("expected ParcelWriter.ArchiveTx to be called")
		}

		if !mockAudit.createCalled {
			t.Fatal("expected AuditRepository.CreateTx to be called")
		}

		if !mockCache.deleteCalled {
			t.Fatal("expected ParcelCache.DeleteByTrack to be called")
		}

	})
}

func TestArchive_Success(t *testing.T) {
	parcelID := 1

	mockReader := &mockParcelRepo{getByIDResult: testParcel(domain.StatusDelivered, false)}

	mockWriter := &mockParcelRepo{}

	mockCache := &mockParcelCache{}

	mockTxManager := &mockTransactionManager{}

	mockAudit := &mockAuditRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		parcelWriter: mockWriter,
		auditRepo:    mockAudit,
		parcelCache:  mockCache,
		txManager:    mockTxManager,
	}

	t.Run("success", func(t *testing.T) {
		err := svc.Archive(context.Background(), parcelID, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockTxManager.doCalled {
			t.Fatal("expected TransactionManager.Do to be called")
		}

		if !mockWriter.archiveCalled {
			t.Fatal("expected ParcelWriter.ArchiveTx to be called")
		}

		if !mockAudit.createCalled {
			t.Fatal("expected AuditRepository.CreateTx to be called")
		}

		if !mockCache.deleteCalled {
			t.Fatal("expected ParcelCache.DeleteByTrack to be called")
		}
	})
}
