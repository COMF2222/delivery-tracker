package service

import (
	"context"
	"delivery-tracker/internal/cache"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"errors"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

type mockParcelCache struct {
	getResult *domain.ParcelDetails
	getErr    error
	setErr    error

	getCalled bool
	setCalled bool
}

type mockParcelRepo struct {
	getByTrackCalled bool

	getResult *domain.Parcel
	getErr    error
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

func (m *mockParcelCache) GetByTrack(ctx context.Context, trackNumber string) (*domain.ParcelDetails, error) {
	m.getCalled = true
	return m.getResult, m.getErr
}

func (m *mockParcelCache) SetByTrack(ctx context.Context, trackNumber string, parcel *domain.ParcelDetails, ttl time.Duration) error {
	m.setCalled = true
	return m.setErr
}

func (m *mockParcelCache) DeleteByTrack(ctx context.Context, trackNumber string) error {
	return nil
}

func (m *mockParcelRepo) GetByTrackNumber(trackNumber string) (*domain.Parcel, error) {
	m.getByTrackCalled = true
	return m.getResult, m.getErr
}

func (m *mockParcelRepo) GetByID(id int) (*domain.Parcel, error) {
	return nil, errors.New("GetByID should not be called")
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

func TestParcelService_GetByTrackNumber_CacheHit(t *testing.T) {
	trackNumber := "Q4P405SHH8EG"

	mockCache := &mockParcelCache{
		getResult: &domain.ParcelDetails{
			TrackNumber: trackNumber,
			ItemName:    "Reina Petite Ring",
		},
	}

	mockReader := &mockParcelRepo{}

	svc := &ParcelService{parcelReader: mockReader, parcelCache: mockCache}

	t.Run("success", func(t *testing.T) {
		got, err := svc.GetByTrackNumber(context.Background(), trackNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.TrackNumber != trackNumber {
			t.Fatalf("expected track number %s, got %s", trackNumber, got.TrackNumber)
		}

		if !mockCache.getCalled {
			t.Fatal("expected ParcelCache.GetByTrack to be called")
		}

		if mockCache.setCalled {
			t.Fatal("ParcelCache.SetByTrack should not be called on cache hit")
		}

		if mockReader.getByTrackCalled {
			t.Fatal("ParcelReader.GetByTrackNumber should not be called on cache hit")
		}
	})
}

func TestParcelService_GetByTrackNumber_CacheMiss(t *testing.T) {
	trackNumber := "Q4P405SHH8EG"

	mockCache := &mockParcelCache{
		getErr: cache.ErrCacheMiss,
	}

	mockReader := &mockParcelRepo{
		getResult: &domain.Parcel{
			ID:            1,
			TrackNumber:   "Q4P405SHH8EG",
			ItemName:      "Reina Petite Ring",
			RecipientName: "Иван",
			CurrentStatus: domain.StatusCreated,
		},
	}

	mockPhoto := &mockParcelPhotoRepo{
		getResult: []domain.ParcelPhoto{
			{ID: 1, ParcelID: 1, FilePath: "/uploads/__1231132.jpg"},
		},
	}

	mockHistory := &mockParcelHistoryRepo{
		getResult: []domain.ParcelStatusHistory{
			{
				ID:        1,
				ParcelID:  1,
				OldStatus: nil,
				NewStatus: domain.StatusPurchased,
				Location:  "loc",
				ChangedBy: 1,
			},
		},
	}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		historyRepo:  mockHistory,
		parcelCache:  mockCache,
	}

	t.Run("success", func(t *testing.T) {
		got, err := svc.GetByTrackNumber(context.Background(), trackNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.TrackNumber != trackNumber {
			t.Fatalf("expected track number %s, got %s", trackNumber, got.TrackNumber)
		}

		if !mockCache.getCalled {
			t.Fatal("expected ParcelCache.GetByTrack to be called")
		}

		if !mockCache.setCalled {
			t.Fatal("expected ParcelCache.SetByTrack to be called after cache miss")
		}

		if !mockReader.getByTrackCalled {
			t.Fatal("expected ParcelReader.GetByTrackNumber to be called")
		}

		if !mockPhoto.getByIdCalled {
			t.Fatal("expected ParcelPhotoRepository.GetByParcelID to be called")
		}

		if !mockHistory.getByIdCalled {
			t.Fatal("expected ParcelStatusHistoryRepository.GetByParcelID to be called")
		}
	})
}

func TestParcelService_GetByTrackNumber_CacheSetError(t *testing.T) {
	trackNumber := "Q4P405SHH8EG"

	setErr := errors.New("set cache error")

	mockCache := &mockParcelCache{
		getErr: cache.ErrCacheMiss,
		setErr: setErr,
	}

	mockReader := &mockParcelRepo{
		getResult: &domain.Parcel{
			ID:            1,
			TrackNumber:   "Q4P405SHH8EG",
			ItemName:      "Reina Petite Ring",
			RecipientName: "Иван",
			CurrentStatus: domain.StatusCreated,
		},
	}

	mockPhoto := &mockParcelPhotoRepo{
		getResult: []domain.ParcelPhoto{
			{ID: 1, ParcelID: 1, FilePath: "/uploads/__1231132.jpg"},
		},
	}

	mockHistory := &mockParcelHistoryRepo{
		getResult: []domain.ParcelStatusHistory{
			{
				ID:        1,
				ParcelID:  1,
				OldStatus: nil,
				NewStatus: domain.StatusPurchased,
				Location:  "loc",
				ChangedBy: 1,
			},
		},
	}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		historyRepo:  mockHistory,
		parcelCache:  mockCache,
	}

	t.Run("set cache error", func(t *testing.T) {
		_, err := svc.GetByTrackNumber(context.Background(), trackNumber)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockCache.getCalled {
			t.Fatal("expected ParcelCache.GetByTrack to be called")
		}

		if !mockCache.setCalled {
			t.Fatal("expected cache ParcelCache.SetByTrack to be called")
		}

		if !mockReader.getByTrackCalled {
			t.Fatal("expected ParcelReader.GetByTrackNumber to be called")
		}

		if !mockPhoto.getByIdCalled {
			t.Fatal("expected ParcelPhotoRepository.GetByParcelID to be called")
		}

		if !mockHistory.getByIdCalled {
			t.Fatal("expected ParcelStatusHistoryRepository.GetByParcelID to be called")
		}
	})
}

func TestParcelService_GetByTrackNumber_RepoError(t *testing.T) {
	trackNumber := "Q4P405SHH8EG"

	mockCache := &mockParcelCache{
		getErr: cache.ErrCacheMiss,
	}

	mockReader := &mockParcelRepo{
		getErr: repository.ErrParcelNotFound,
	}

	mockPhoto := &mockParcelPhotoRepo{}

	mockHistory := &mockParcelHistoryRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		historyRepo:  mockHistory,
		parcelCache:  mockCache,
	}

	t.Run("parcel repo error", func(t *testing.T) {
		_, err := svc.GetByTrackNumber(context.Background(), trackNumber)
		if !errors.Is(err, repository.ErrParcelNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}

		if !mockCache.getCalled {
			t.Fatal("expected ParcelCache.GetByTrack to be called")
		}

		if mockCache.setCalled {
			t.Fatal("ParcelCache.SetByTrack should not be called when ParcelReader.GetByTrackNumber returns error")
		}

		if !mockReader.getByTrackCalled {
			t.Fatal("expected ParcelReader.GetByTrackNumber to be called")
		}

		if mockPhoto.getByIdCalled {
			t.Fatal("ParcelPhotoRepository.GetByParcelID should not be called when ParcelReader.GetByTrackNumber returns error")
		}

		if mockHistory.getByIdCalled {
			t.Fatal("ParcelStatusHistoryRepository.GetByParcelID should not be called when ParcelReader.GetByTrackNumber returns error")
		}
	})
}

func TestParcelService_GetByTrackNumber_PhotoRepoError(t *testing.T) {
	trackNumber := "Q4P405SHH8EG"

	photoErr := errors.New("photo repo error")

	mockCache := &mockParcelCache{
		getErr: cache.ErrCacheMiss,
	}

	mockReader := &mockParcelRepo{
		getResult: &domain.Parcel{
			ID:            1,
			TrackNumber:   "Q4P405SHH8EG",
			ItemName:      "Reina Petite Ring",
			RecipientName: "Иван",
			CurrentStatus: domain.StatusCreated,
		},
	}

	mockPhoto := &mockParcelPhotoRepo{
		getErr: photoErr,
	}

	mockHistory := &mockParcelHistoryRepo{}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		historyRepo:  mockHistory,
		parcelCache:  mockCache,
	}

	t.Run("photo repo error", func(t *testing.T) {
		_, err := svc.GetByTrackNumber(context.Background(), trackNumber)
		if !errors.Is(err, photoErr) {
			t.Fatalf("got %v, want %v", err, photoErr)
		}

		if !mockCache.getCalled {
			t.Fatal("expected ParcelCache.GetByTrack to be called")
		}

		if mockCache.setCalled {
			t.Fatal("ParcelCache.SetByTrack should not be called when ParcelPhotoRepository.GetByParcelID returns error")
		}

		if !mockReader.getByTrackCalled {
			t.Fatal("expected ParcelReader.GetByTrackNumber to be called")
		}

		if !mockPhoto.getByIdCalled {
			t.Fatal("expected ParcelPhotoRepository.GetByParcelID to be called")
		}

		if mockHistory.getByIdCalled {
			t.Fatal("ParcelStatusHistoryRepository.GetByParcelID should not be called when ParcelPhotoRepository.GetByParcelID returns error")
		}
	})
}

func TestParcelService_GetByTrackNumber_HistoryRepoError(t *testing.T) {
	trackNumber := "Q4P405SHH8EG"

	historyErr := errors.New("history repo error")

	mockCache := &mockParcelCache{
		getErr: cache.ErrCacheMiss,
	}

	mockReader := &mockParcelRepo{
		getResult: &domain.Parcel{
			ID:            1,
			TrackNumber:   "Q4P405SHH8EG",
			ItemName:      "Reina Petite Ring",
			RecipientName: "Иван",
			CurrentStatus: domain.StatusCreated,
		},
	}

	mockPhoto := &mockParcelPhotoRepo{
		getResult: []domain.ParcelPhoto{
			{ID: 1, ParcelID: 1, FilePath: "/uploads/__1231132.jpg"},
		},
	}

	mockHistory := &mockParcelHistoryRepo{
		getErr: historyErr,
	}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		historyRepo:  mockHistory,
		parcelCache:  mockCache,
	}

	t.Run("history repo error", func(t *testing.T) {
		_, err := svc.GetByTrackNumber(context.Background(), trackNumber)
		if !errors.Is(err, historyErr) {
			t.Fatalf("got %v, want %v", err, historyErr)
		}

		if !mockCache.getCalled {
			t.Fatal("expected ParcelCache.GetByTrack to be called")
		}

		if mockCache.setCalled {
			t.Fatal("ParcelCache.SetByTrack should not be called when ParcelStatusHistoryRepository.GetByParcelID returns error")
		}

		if !mockReader.getByTrackCalled {
			t.Fatal("expected ParcelReader.GetByTrackNumber to be called")
		}

		if !mockPhoto.getByIdCalled {
			t.Fatal("expected ParcelPhotoRepository.GetByParcelID to be called")
		}

		if !mockHistory.getByIdCalled {
			t.Fatal("expected ParcelStatusHistoryRepository.GetByParcelID to be called")
		}
	})
}
