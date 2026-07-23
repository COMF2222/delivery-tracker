package service

import (
	"delivery-tracker/internal/domain"
	"errors"
	"testing"
)

func TestList_InvalidPage(t *testing.T) {
	status := domain.StatusCreated
	page := 0
	limit := 10

	mockList := &mockLister{}

	svc := ParcelService{parcelLister: mockList}

	t.Run("list invalid page", func(t *testing.T) {
		_, _, err := svc.List(status, page, limit)
		if !errors.Is(err, ErrInvalidPage) {
			t.Fatalf("got %v, want %v", err, ErrInvalidPage)
		}

		if mockList.listCalled {
			t.Fatalf("ParcelLister.List should not be called")
		}

		if mockList.countCalled {
			t.Fatalf("ParcelLister.Count should not be called")
		}
	})
}

func TestList_InvalidLimitZero(t *testing.T) {
	status := domain.StatusCreated
	page := 1
	limit := 0

	mockList := &mockLister{}

	svc := ParcelService{parcelLister: mockList}

	t.Run("list invalid limit zero", func(t *testing.T) {
		_, _, err := svc.List(status, page, limit)
		if !errors.Is(err, ErrInvalidLimit) {
			t.Fatalf("got %v, want %v", err, ErrInvalidLimit)
		}

		if mockList.listCalled {
			t.Fatalf("ParcelLister.List should not be called")
		}

		if mockList.countCalled {
			t.Fatalf("ParcelLister.Count should not be called")
		}
	})
}

func TestList_InvalidLimitTooHigh(t *testing.T) {
	status := domain.StatusCreated
	page := 1
	limit := 101

	mockList := &mockLister{}

	svc := ParcelService{parcelLister: mockList}

	t.Run("list invalid limit to high", func(t *testing.T) {
		_, _, err := svc.List(status, page, limit)
		if !errors.Is(err, ErrInvalidLimit) {
			t.Fatalf("got %v, want %v", err, ErrInvalidLimit)
		}

		if mockList.listCalled {
			t.Fatalf("ParcelLister.List should not be called")
		}

		if mockList.countCalled {
			t.Fatalf("ParcelLister.Count should not be called")
		}
	})
}

func TestList_ListSuccessWithoutStatus(t *testing.T) {
	status := domain.Status("")
	page := 2
	limit := 10
	offset := (page - 1) * limit

	mockList := &mockLister{listResult: testListParcel(status, false), countResult: 2}

	svc := ParcelService{parcelLister: mockList}

	t.Run("list success without status", func(t *testing.T) {
		parcels, total, err := svc.List(status, page, limit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockList.listCalled {
			t.Fatalf("expected ParcelLister.List to be called")
		}

		if !mockList.countCalled {
			t.Fatalf("expected ParcelLister.Count to be called")
		}

		if mockList.listByStatusCalled {
			t.Fatal("ParcelLister.ListByStatus should not be called")
		}

		if mockList.countByStatusCalled {
			t.Fatal("ParcelLister.CountByStatus should not be called")
		}

		if mockList.listOffset != offset {
			t.Fatalf("got %v, want %v offset", mockList.listOffset, offset)
		}

		for i := range parcels {
			if mockList.listResult[i] != parcels[i] {
				t.Fatalf("got %v, want %v parcel", mockList.listResult[i], parcels[i])
			}
		}

		if mockList.countResult != total {
			t.Fatalf("got %v, want %v total", mockList.countResult, total)
		}

		if mockList.listLimit != limit {
			t.Fatalf("got %d, want %d limit", mockList.listLimit, limit)
		}
	})
}

func TestList_ListErrorWithoutStatus(t *testing.T) {
	status := domain.Status("")
	page := 2
	limit := 10

	listErr := errors.New("failed to get parcel list")

	mockList := &mockLister{listErr: listErr}

	svc := ParcelService{parcelLister: mockList}

	t.Run("list error without status", func(t *testing.T) {
		_, _, err := svc.List(status, page, limit)
		if !errors.Is(err, listErr) {
			t.Fatalf("got %v, want %v", err, listErr)
		}

		if !mockList.listCalled {
			t.Fatalf("expected ParcelLister.List to be called")
		}

		if mockList.countCalled {
			t.Fatalf("ParcelLister.Count should not be called")
		}
	})
}

func TestList_CountErrorWithoutStatus(t *testing.T) {
	status := domain.Status("")
	page := 2
	limit := 10

	countErr := errors.New("failed to count parcels")

	mockList := &mockLister{countErr: countErr}

	svc := ParcelService{parcelLister: mockList}

	t.Run("count error without status", func(t *testing.T) {
		parcels, total, err := svc.List(status, page, limit)
		if !errors.Is(err, countErr) {
			t.Fatalf("got %v, want %v", err, countErr)
		}

		if !mockList.listCalled {
			t.Fatalf("expected ParcelLister.List to be called")
		}

		if !mockList.countCalled {
			t.Fatalf("expected ParcelLister.Count to be called")
		}

		if parcels != nil {
			t.Fatalf("parcels should be nil, but got %v", parcels)
		}

		if total != 0 {
			t.Fatalf("total should be nil, but got %v", total)
		}
	})
}

func TestList_ListSuccessWithStatus(t *testing.T) {
	status := domain.StatusCreated
	page := 3
	limit := 20
	offset := (page - 1) * limit

	mockList := &mockLister{listResult: testListParcel(status, false), countResult: 2}

	svc := ParcelService{parcelLister: mockList}

	t.Run("list success with status", func(t *testing.T) {
		parcels, total, err := svc.List(status, page, limit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !mockList.listByStatusCalled {
			t.Fatalf("expected ParcelLister.ListByStatus to be called")
		}

		if !mockList.countByStatusCalled {
			t.Fatalf("expected ParcelLister.CountByStatus to be called")
		}

		if mockList.listCalled {
			t.Fatal("ParcelLister.List should not be called")
		}

		if mockList.countCalled {
			t.Fatal("ParcelLister.Count should not be called")
		}

		if mockList.listOffset != offset {
			t.Fatalf("got %v, want %v offset", mockList.listOffset, offset)
		}

		for i := range parcels {
			if mockList.listResult[i] != parcels[i] {
				t.Fatalf("got %v, want %v parcel", mockList.listResult[i], parcels[i])
			}
		}

		if mockList.countResult != total {
			t.Fatalf("got %v, want %v total", mockList.countResult, total)
		}

		if mockList.listLimit != limit {
			t.Fatalf("got %d, want %d limit", mockList.listLimit, limit)
		}

		if mockList.receivedStatus != status {
			t.Fatalf("got %v, want %v status", mockList.receivedStatus, status)
		}
	})
}

func TestList_ListErrorWithStatus(t *testing.T) {
	status := domain.StatusCreated
	page := 3
	limit := 20

	listErr := errors.New("failed to get parcel list by status")

	mockList := &mockLister{listErr: listErr}

	svc := ParcelService{parcelLister: mockList}

	t.Run("list error with status", func(t *testing.T) {
		_, _, err := svc.List(status, page, limit)
		if !errors.Is(err, listErr) {
			t.Fatalf("got %v, want %v", err, listErr)
		}

		if !mockList.listByStatusCalled {
			t.Fatalf("expected ParcelLister.ListByStatus to be called")
		}

		if mockList.countByStatusCalled {
			t.Fatalf("ParcelLister.CountByStatus should not be called")
		}

		if mockList.listCalled {
			t.Fatal("ParcelLister.List should not be called")
		}

		if mockList.countCalled {
			t.Fatal("ParcelLister.Count should not be called")
		}
	})
}

func TestList_CountErrorWithStatus(t *testing.T) {
	status := domain.StatusCreated
	page := 2
	limit := 10

	countErr := errors.New("failed to count parcels by status")

	mockList := &mockLister{countErr: countErr}

	svc := ParcelService{parcelLister: mockList}

	t.Run("count error with status", func(t *testing.T) {
		parcels, total, err := svc.List(status, page, limit)
		if !errors.Is(err, countErr) {
			t.Fatalf("got %v, want %v", err, countErr)
		}

		if !mockList.listByStatusCalled {
			t.Fatalf("expected ParcelLister.ListByStatus to be called")
		}

		if !mockList.countByStatusCalled {
			t.Fatalf("expected ParcelLister.CountByStatus to be called")
		}

		if mockList.listCalled {
			t.Fatal("ParcelLister.List should not be called")
		}

		if mockList.countCalled {
			t.Fatal("ParcelLister.Count should not be called")
		}

		if parcels != nil {
			t.Fatalf("parcels should be nil, but got %v", parcels)
		}

		if total != 0 {
			t.Fatalf("total should be nil, but got %v", total)
		}
	})
}
