package service

import (
	"context"
	"delivery-tracker/internal/cache"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/generator"
	"delivery-tracker/internal/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"time"
)

type ParcelService struct {
	parcelReader ParcelReader
	parcelWriter ParcelWriter
	parcelLister ParcelLister

	statusRepo  StatusRepository
	photoRepo   ParcelPhotoRepository
	historyRepo ParcelStatusHistoryRepository
	auditRepo   AuditRepository

	parcelCache ParcelCache
	txManager   TransactionManager

	trackGenerator TrackGenerator
}

func NewParcelService(parcelRepo *repository.ParcelRepository,
	statusRepo *repository.StatusRepository,
	photoRepo *repository.ParcelPhotoRepository,
	historyRepo *repository.ParcelStatusHistoryRepository,
	auditRepo *repository.AuditRepository,
	parcelCache *cache.ParcelCache,
	txManager *repository.TransactionManager,
	trackGenerator *generator.TrackNumberGenerator) *ParcelService {
	return &ParcelService{
		parcelReader:   parcelRepo,
		parcelWriter:   parcelRepo,
		parcelLister:   parcelRepo,
		statusRepo:     statusRepo,
		photoRepo:      photoRepo,
		historyRepo:    historyRepo,
		auditRepo:      auditRepo,
		parcelCache:    parcelCache,
		txManager:      txManager,
		trackGenerator: trackGenerator,
	}
}

func (s *ParcelService) CreateParcel(parcel *domain.Parcel) error {
	parcel.CurrentStatus = domain.StatusCreated
	parcel.IsArchived = false

	statusID, err := s.statusRepo.GetStatusID(domain.StatusCreated)
	if err != nil {
		return fmt.Errorf("get created status id: %w", err)
	}

	for attempt := 0; attempt < 5; attempt++ {
		track, err := s.trackGenerator.GenerateTrackNumber()
		if err != nil {
			return fmt.Errorf("generate track number: %w", err)
		}
		parcel.TrackNumber = track

		err = s.parcelWriter.CreateParcel(parcel, statusID)
		if err == nil {
			return nil
		}

		if errors.Is(err, repository.ErrTrackNumberAlreadyExists) {
			continue
		}

		return fmt.Errorf("create parcel: %w", err)

	}

	return ErrFailedToGenerateUniqueTrack
}

func (s *ParcelService) GetByTrackNumber(ctx context.Context, trackNumber string) (*domain.ParcelDetails, error) {
	cachedDetails, err := s.parcelCache.GetByTrack(ctx, trackNumber)
	if err == nil {
		return cachedDetails, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("failed to get cache by track: %v", err)
	}

	parcel, err := s.parcelReader.GetByTrackNumber(trackNumber)
	if err != nil {
		return nil, fmt.Errorf("get by track number: %w", err)
	}

	parcelPhotos, err := s.photoRepo.GetByParcelID(parcel.ID)
	if err != nil {
		return nil, fmt.Errorf("get photo by parcel id(%d): %w", parcel.ID, err)
	}

	parcelHistory, err := s.historyRepo.GetByParcelID(parcel.ID)
	if err != nil {
		return nil, fmt.Errorf("get history by parcel id(%d): %w", parcel.ID, err)
	}

	parcelDetails := &domain.ParcelDetails{
		TrackNumber:     parcel.TrackNumber,
		ItemName:        parcel.ItemName,
		RecipientName:   parcel.RecipientName,
		CurrentStatus:   parcel.CurrentStatus,
		CurrentLocation: parcel.CurrentLocation,
		History:         parcelHistory,
		Photos:          parcelPhotos,
	}

	err = s.parcelCache.SetByTrack(ctx, trackNumber, parcelDetails, 30*time.Minute)
	if err != nil {
		log.Printf("failed to cache track number: %v", err)
	}

	return parcelDetails, nil
}

func (s *ParcelService) ChangeStatus(
	ctx context.Context,
	parcelID int,
	newStatus domain.Status,
	location string,
	changedBy int) error {
	parcel, err := s.parcelReader.GetByID(parcelID)
	if err != nil {
		return fmt.Errorf("failed to get parcel by ID(%d): %w", parcelID, err)
	}

	canChangeStatus := domain.CanChangeStatus(parcel.CurrentStatus, newStatus)
	if !canChangeStatus {
		return ErrInvalidStatusTransition
	}

	oldStatusID, err := s.statusRepo.GetStatusID(parcel.CurrentStatus)
	if err != nil {
		return fmt.Errorf("failed to get status ID: %w", err)
	}

	newStatusID, err := s.statusRepo.GetStatusID(newStatus)
	if err != nil {
		return fmt.Errorf("failed to get status ID: %w", err)
	}

	err = s.txManager.Do(func(tx *sqlx.Tx) error {
		if err := s.parcelWriter.UpdateStatusTx(tx, parcelID, newStatusID, location); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}

		oldStatus := parcel.CurrentStatus

		history := domain.ParcelStatusHistory{
			ParcelID:  parcelID,
			OldStatus: &oldStatus,
			NewStatus: newStatus,
			Location:  location,
			ChangedBy: changedBy,
		}
		if err := s.historyRepo.CreateTx(tx, &history, oldStatusID, newStatusID); err != nil {
			return fmt.Errorf("failed to create parcel status history: %w", err)
		}

		auditLog := domain.AuditLog{
			UserID:     changedBy,
			Action:     domain.ActionChangeStatus,
			OldValue:   string(parcel.CurrentStatus),
			NewValue:   string(newStatus),
			EntityType: domain.EntityTypeParcel,
			EntityID:   parcelID,
		}

		if err := s.auditRepo.CreateTx(tx, &auditLog); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("change parcel status transaction: %w", err)
	}

	err = s.parcelCache.DeleteByTrack(ctx, parcel.TrackNumber)
	if err != nil {
		log.Printf("failed to delete cache by track: %v", err)
	}

	return nil
}

func (s *ParcelService) AddPhoto(ctx context.Context, parcelID int, filePath string) error {
	parcel, err := s.parcelReader.GetByID(parcelID)
	if err != nil {
		return fmt.Errorf("get parcel by id: %w", err)
	}

	photo := domain.ParcelPhoto{
		ParcelID: parcelID,
		FilePath: filePath,
	}

	err = s.photoRepo.Create(&photo)
	if err != nil {
		return fmt.Errorf("create parcel photo: %w", err)
	}

	err = s.parcelCache.DeleteByTrack(ctx, parcel.TrackNumber)
	if err != nil {
		log.Printf("failed to delete cache by track: %v", err)
	}

	return nil
}

func (s *ParcelService) Archive(ctx context.Context, parcelID, changedBy int) error {
	parcel, err := s.parcelReader.GetByID(parcelID)
	if err != nil {
		return fmt.Errorf("get parcel id: %w", err)
	}

	if parcel.IsArchived {
		return ErrParcelAlreadyArchived
	}

	if parcel.CurrentStatus != domain.StatusDelivered {
		return ErrParcelNotDelivered
	}

	err = s.txManager.Do(func(tx *sqlx.Tx) error {
		oldValue := parcel.IsArchived
		if err := s.parcelWriter.ArchiveTx(tx, parcelID); err != nil {
			return fmt.Errorf("failed to archive parcel: %w", err)
		}

		auditLog := domain.AuditLog{
			UserID:     changedBy,
			Action:     domain.ActionArchiveParcel,
			OldValue:   strconv.FormatBool(oldValue),
			NewValue:   strconv.FormatBool(true),
			EntityType: domain.EntityTypeParcel,
			EntityID:   parcelID,
		}

		if err = s.auditRepo.CreateTx(tx, &auditLog); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("archive parcel transaction: %w", err)
	}

	err = s.parcelCache.DeleteByTrack(ctx, parcel.TrackNumber)
	if err != nil {
		log.Printf("failed to delete cache by track: %v", err)
	}

	return nil
}

func (s *ParcelService) List(status domain.Status, page, limit int) ([]domain.Parcel, int, error) {
	if page < 1 {
		return nil, 0, ErrInvalidPage
	}

	if limit < 1 {
		return nil, 0, ErrInvalidLimit
	}

	if limit > 100 {
		return nil, 0, ErrInvalidLimit
	}

	offset := (page - 1) * limit

	if status == "" {
		parcels, err := s.parcelLister.List(limit, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get parcel list: %w", err)
		}

		total, err := s.parcelLister.Count()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count parcels: %w", err)
		}

		return parcels, total, nil
	}

	parcels, err := s.parcelLister.ListByStatus(status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get parcel list by status: %w", err)
	}

	total, err := s.parcelLister.CountByStatus(status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count parcels: %w", err)
	}

	return parcels, total, nil
}
