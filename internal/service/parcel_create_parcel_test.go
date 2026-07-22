package service

import (
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/repository"
	"errors"
	"testing"
)

func TestCreateParcel_StatusRepoError(t *testing.T) {
	parcel := domain.Parcel{
		ItemName:         "ItemName",
		RecipientName:    "RecipientName",
		RecipientPhone:   "RecipientPhone",
		RecipientAddress: "RecipientAddress",
	}

	statusErr := errors.New("status repo err")

	mockWriter := &mockParcelRepo{}

	mockStatus := &mockStatusRepo{
		getErr: statusErr,
	}

	mockTrackGenerator := &mockGenerator{
		generatorResults: []string{"TRACK001", "TRACK002"},
		generatorErrs:    []error{nil, nil},
	}

	svc := ParcelService{parcelWriter: mockWriter, statusRepo: mockStatus, trackGenerator: mockTrackGenerator}

	t.Run("status repo error", func(t *testing.T) {
		err := svc.CreateParcel(&parcel)
		if !errors.Is(err, statusErr) {
			t.Fatalf("got %v, want %v", err, statusErr)
		}

		if mockWriter.createParcelCalls != 0 {
			t.Fatalf("CreateParcel called %d times, want 0", mockWriter.createParcelCalls)
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if mockTrackGenerator.generatorCalls != 0 {
			t.Fatalf(
				"GenerateTrackNumber called %d times, want 0",
				mockTrackGenerator.generatorCalls,
			)
		}
	})
}

func TestCreateParcel_CreateError(t *testing.T) {
	parcel := domain.Parcel{
		ItemName:         "ItemName",
		RecipientName:    "RecipientName",
		RecipientPhone:   "RecipientPhone",
		RecipientAddress: "RecipientAddress",
	}

	createError := errors.New("create error")

	mockWriter := &mockParcelRepo{createParcelErr: createError}

	mockStatus := &mockStatusRepo{getResult: 1}

	mockTrackGenerator := &mockGenerator{
		generatorResults: []string{"TRACK001"},
		generatorErrs:    []error{nil},
	}

	svc := ParcelService{parcelWriter: mockWriter, statusRepo: mockStatus, trackGenerator: mockTrackGenerator}

	t.Run("create error", func(t *testing.T) {
		err := svc.CreateParcel(&parcel)
		if !errors.Is(err, createError) {
			t.Fatalf("got %v, want %v", err, createError)
		}

		if mockWriter.createParcelCalls != 1 {
			t.Fatalf("CreateParcel called %d times, want 1", mockWriter.createParcelCalls)
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}
	})
}

func TestCreateParcel_CreateSuccess(t *testing.T) {
	parcel := domain.Parcel{
		ItemName:         "ItemName",
		RecipientName:    "RecipientName",
		RecipientPhone:   "RecipientPhone",
		RecipientAddress: "RecipientAddress",
	}

	mockWriter := &mockParcelRepo{}

	mockStatus := &mockStatusRepo{getResult: 1}

	mockTrackGenerator := &mockGenerator{
		generatorResults: []string{"TRACK001"},
		generatorErrs:    []error{nil},
	}

	svc := ParcelService{parcelWriter: mockWriter, statusRepo: mockStatus, trackGenerator: mockTrackGenerator}

	t.Run("success", func(t *testing.T) {
		err := svc.CreateParcel(&parcel)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if mockWriter.createParcelCalls != 1 {
			t.Fatalf("CreateParcel called %d times, want 1", mockWriter.createParcelCalls)
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if mockTrackGenerator.generatorCalls != 1 {
			t.Fatalf(
				"GenerateTrackNumber called %d times, want 1",
				mockTrackGenerator.generatorCalls,
			)
		}

		if parcel.TrackNumber != "TRACK001" {
			t.Fatalf("got track number %q, want %q", parcel.TrackNumber, "TRACK001")
		}

		if parcel.CurrentStatus != domain.StatusCreated {
			t.Fatalf("got current status %q, want %q", parcel.CurrentStatus, domain.StatusCreated)
		}

		if parcel.IsArchived {
			t.Fatal("parcel should not be archived")
		}
	})
}

func TestCreateParcel_TrackCollisionThenSuccess(t *testing.T) {
	parcel := domain.Parcel{
		ItemName:         "ItemName",
		RecipientName:    "RecipientName",
		RecipientPhone:   "RecipientPhone",
		RecipientAddress: "RecipientAddress",
	}

	mockWriter := &mockParcelRepo{
		createParcelErrs: []error{repository.ErrTrackNumberAlreadyExists, nil},
	}

	mockStatus := &mockStatusRepo{getResult: 1}

	mockTrackGenerator := &mockGenerator{
		generatorResults: []string{"TRACK001", "TRACK002"},
		generatorErrs:    []error{nil, nil},
	}

	svc := ParcelService{parcelWriter: mockWriter, statusRepo: mockStatus, trackGenerator: mockTrackGenerator}

	t.Run("track collision then success", func(t *testing.T) {
		err := svc.CreateParcel(&parcel)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if mockWriter.createParcelCalls != 2 {
			t.Fatalf("CreateParcel called %d times, want 2", mockWriter.createParcelCalls)
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if mockTrackGenerator.generatorCalls != 2 {
			t.Fatalf(
				"GenerateTrackNumber called %d times, want 2",
				mockTrackGenerator.generatorCalls,
			)
		}

		if parcel.TrackNumber != "TRACK002" {
			t.Fatalf("got track number %q, want %q", parcel.TrackNumber, "TRACK002")
		}
	})
}

func TestCreateParcel_TrackCollisionAllAttempts(t *testing.T) {
	parcel := domain.Parcel{
		ItemName:         "ItemName",
		RecipientName:    "RecipientName",
		RecipientPhone:   "RecipientPhone",
		RecipientAddress: "RecipientAddress",
	}

	mockWriter := &mockParcelRepo{
		createParcelErrs: []error{
			repository.ErrTrackNumberAlreadyExists,
			repository.ErrTrackNumberAlreadyExists,
			repository.ErrTrackNumberAlreadyExists,
			repository.ErrTrackNumberAlreadyExists,
			repository.ErrTrackNumberAlreadyExists,
		},
	}

	mockStatus := &mockStatusRepo{getResult: 1}

	mockTrackGenerator := &mockGenerator{
		generatorResults: []string{"TRACK001", "TRACK002", "TRACK003", "TRACK004", "TRACK005"},
		generatorErrs:    []error{nil, nil, nil, nil, nil},
	}

	svc := ParcelService{parcelWriter: mockWriter, statusRepo: mockStatus, trackGenerator: mockTrackGenerator}

	t.Run("track collision all attempts", func(t *testing.T) {
		err := svc.CreateParcel(&parcel)
		if !errors.Is(err, ErrFailedToGenerateUniqueTrack) {
			t.Fatalf("got %v, want %v", err, ErrFailedToGenerateUniqueTrack)
		}

		if mockWriter.createParcelCalls != 5 {
			t.Fatalf("CreateParcel called %d times, want 5", mockWriter.createParcelCalls)
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if mockTrackGenerator.generatorCalls != 5 {
			t.Fatalf(
				"GenerateTrackNumber called %d times, want 5",
				mockTrackGenerator.generatorCalls,
			)
		}
	})
}

func TestCreateParcel_TrackGeneratorError(t *testing.T) {
	parcel := domain.Parcel{
		ItemName:         "ItemName",
		RecipientName:    "RecipientName",
		RecipientPhone:   "RecipientPhone",
		RecipientAddress: "RecipientAddress",
	}

	generatorError := errors.New("generator err")

	mockWriter := &mockParcelRepo{}

	mockStatus := &mockStatusRepo{getResult: 1}

	mockTrackGenerator := &mockGenerator{
		generatorErr: generatorError,
	}

	svc := ParcelService{parcelWriter: mockWriter, statusRepo: mockStatus, trackGenerator: mockTrackGenerator}

	t.Run("track generator error", func(t *testing.T) {
		err := svc.CreateParcel(&parcel)
		if !errors.Is(err, generatorError) {
			t.Fatalf("got %v, want %v", err, generatorError)
		}

		if mockWriter.createParcelCalls != 0 {
			t.Fatalf("CreateParcel called %d times, want 0", mockWriter.createParcelCalls)
		}

		if !mockStatus.getStatusIDCalled {
			t.Fatal("expected StatusRepository.GetStatusID to be called")
		}

		if mockTrackGenerator.generatorCalls != 1 {
			t.Fatalf(
				"GenerateTrackNumber called %d times, want 1",
				mockTrackGenerator.generatorCalls,
			)
		}
	})
}
