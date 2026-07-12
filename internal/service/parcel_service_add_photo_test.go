package service

import (
	"context"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"errors"
	"testing"
)

func TestAddPhoto_ParcelNotFound(t *testing.T) {
	parcelID := 2

	mockReader := &mockParcelRepo{
		getByIDErr: repository.ErrParcelNotFound,
	}

	mockPhoto := &mockParcelPhotoRepo{}

	mockCache := &mockParcelCache{}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		parcelCache:  mockCache,
	}

	t.Run("parcel not found error", func(t *testing.T) {
		err := svc.AddPhoto(context.Background(), parcelID, "/uploads")
		if !errors.Is(err, repository.ErrParcelNotFound) {
			t.Fatalf("got %v, want ErrParcelNotFound", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if mockPhoto.createCalled {
			t.Fatal("ParcelPhotoRepository.Create should not be called")
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}
	})
}

func TestAddPhoto_CreateErr(t *testing.T) {
	parcelID := 1

	createErr := errors.New("photo create err")

	mockReader := &mockParcelRepo{
		getByIDResult: testParcel(domain.StatusCreated, false),
	}

	mockPhoto := &mockParcelPhotoRepo{createErr: createErr}

	mockCache := &mockParcelCache{}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		parcelCache:  mockCache,
	}

	t.Run("create photo error", func(t *testing.T) {
		err := svc.AddPhoto(context.Background(), parcelID, "/uploads")
		if !errors.Is(err, createErr) {
			t.Fatalf("got %v, want createErr", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockPhoto.createCalled {
			t.Fatal("expected ParcelPhotoRepository.Create to be called")
		}

		if mockCache.deleteCalled {
			t.Fatal("ParcelCache.DeleteByTrack should not be called")
		}
	})
}

func TestAddPhoto_CacheDeleteError(t *testing.T) {
	parcelID := 1

	deleteCacheError := errors.New("delete cache by track error")

	mockReader := &mockParcelRepo{
		getByIDResult: testParcel(domain.StatusCreated, false),
	}

	mockPhoto := &mockParcelPhotoRepo{}

	mockCache := &mockParcelCache{deleteErr: deleteCacheError}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		parcelCache:  mockCache,
	}

	t.Run("cache delete error", func(t *testing.T) {
		err := svc.AddPhoto(context.Background(), parcelID, "/uploads")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockPhoto.createCalled {
			t.Fatal("expected ParcelPhotoRepository.Create to be called")
		}

		if !mockCache.deleteCalled {
			t.Fatal("expected ParcelCache.DeleteByTrack to be called")
		}
	})
}

func TestAddPhoto_Success(t *testing.T) {
	parcelID := 1

	mockReader := &mockParcelRepo{
		getByIDResult: testParcel(domain.StatusCreated, false),
	}

	mockPhoto := &mockParcelPhotoRepo{}

	mockCache := &mockParcelCache{}

	svc := ParcelService{
		parcelReader: mockReader,
		photoRepo:    mockPhoto,
		parcelCache:  mockCache,
	}

	t.Run("success", func(t *testing.T) {
		err := svc.AddPhoto(context.Background(), parcelID, "/uploads")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockReader.getByIDCalled {
			t.Fatal("expected ParcelReader.GetByID to be called")
		}

		if !mockPhoto.createCalled {
			t.Fatal("expected ParcelPhotoRepository.Create to be called")
		}

		if !mockCache.deleteCalled {
			t.Fatal("expected ParcelCache.DeleteByTrack to be called")
		}
	})
}
