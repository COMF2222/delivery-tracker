package service

import (
	"context"
	"delivery-tracker/internal/cache"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"errors"
	"testing"
)

func TestParcelService_GetByTrackNumber_CacheHit(t *testing.T) {
	trackNumber := "Q4P405SHH8EG"

	mockCache := &mockParcelCache{
		getResult: testParcelDetails(trackNumber),
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

	mockReader := &mockParcelRepo{getResult: testParcel(domain.StatusCreated)}

	mockPhoto := &mockParcelPhotoRepo{getResult: testPhotos()}

	mockHistory := &mockParcelHistoryRepo{getResult: testHistory()}

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

	mockReader := &mockParcelRepo{getResult: testParcel(domain.StatusCreated)}

	mockPhoto := &mockParcelPhotoRepo{getResult: testPhotos()}

	mockHistory := &mockParcelHistoryRepo{getResult: testHistory()}

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

	t.Run("parcel not found error", func(t *testing.T) {
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

	mockReader := &mockParcelRepo{getResult: testParcel(domain.StatusCreated)}

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

	mockReader := &mockParcelRepo{getResult: testParcel(domain.StatusCreated)}

	mockPhoto := &mockParcelPhotoRepo{getResult: testPhotos()}

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
