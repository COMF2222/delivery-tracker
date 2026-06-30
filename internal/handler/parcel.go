package handler

import (
	"delivery-tracker/internal/contextkeys"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/dto"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/request"
	"delivery-tracker/internal/response"
	"delivery-tracker/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ParcelHandler struct {
	parcelService *service.ParcelService
}

func NewParcelHandler(parcelService *service.ParcelService) *ParcelHandler {
	return &ParcelHandler{parcelService: parcelService}
}

// CreateParcel
//
//	@Summary		Создание посылки
//	@Description	Создаёт посылку и генерирует трек-номер
//	@Tags			Parcel
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateParcelRequest	true	"Create parcel request"
//	@Success		201		{object}	dto.CreateParcelResponse
//	@Failure		400		{object}	response.ErrorResponse	"Bad request"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	response.ErrorResponse	"Forbidden"
//	@Failure		405		{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/parcels [post]
func (h *ParcelHandler) CreateParcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateParcelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := domain.Parcel{
		ItemName:         req.ItemName,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientAddress: req.RecipientAddress,
	}

	if err := h.parcelService.CreateParcel(&parcel); err != nil {
		response.Error(w, "failed to create parcel", http.StatusInternalServerError)
		return
	}

	resp := dto.CreateParcelResponse{
		ID:          parcel.ID,
		TrackNumber: parcel.TrackNumber,
	}

	response.JSON(w, http.StatusCreated, resp)
}

// GetByTrackNumber
//
//	@Summary		Получение посылки по трек-номеру
//	@Description	Возвращает информацию о посылке, историю статусов и фотографии
//	@Tags			Parcel
//	@Produce		json
//	@Param			track_number	query		string	true	"Track number"
//	@Success		200				{object}	dto.GetParcelResponse
//	@Failure		400				{object}	response.ErrorResponse	"Bad request"
//	@Failure		404				{object}	response.ErrorResponse	"Not found"
//	@Failure		405				{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500				{object}	response.ErrorResponse	"Internal server error"
//	@Router			/parcels/track [get]
func (h *ParcelHandler) GetByTrackNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackNumber := r.URL.Query().Get("track_number")
	if trackNumber == "" {
		response.Error(w, "track number cannot be empty", http.StatusBadRequest)
		return
	}

	parcel, err := h.parcelService.GetByTrackNumber(trackNumber)
	if err != nil {
		if errors.Is(err, repository.ErrParcelNotFound) {
			response.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		response.Error(w, "failed to get parcel by track number", http.StatusInternalServerError)
		return
	}
	photosResponse := make([]dto.ParcelPhotoResponse, 0, len(parcel.Photos))
	for _, photo := range parcel.Photos {
		photoResponse := dto.ParcelPhotoResponse{
			FilePath:  photo.FilePath,
			CreatedAt: photo.CreatedAt,
		}
		photosResponse = append(photosResponse, photoResponse)
	}

	historyResponse := make([]dto.ParcelHistoryResponse, 0, len(parcel.History))
	for _, history := range parcel.History {
		historyResponseItem := dto.ParcelHistoryResponse{
			OldStatus: history.OldStatus,
			NewStatus: history.NewStatus,
			Location:  history.Location,
			CreatedAt: history.CreatedAt,
		}
		historyResponse = append(historyResponse, historyResponseItem)
	}
	resp := dto.GetParcelResponse{
		TrackNumber:     parcel.TrackNumber,
		ItemName:        parcel.ItemName,
		Recipient:       parcel.RecipientName,
		CurrentStatus:   parcel.CurrentStatus,
		CurrentLocation: parcel.CurrentLocation,
		History:         historyResponse,
		Photos:          photosResponse,
	}
	response.JSON(w, http.StatusOK, resp)
}

// UpdateStatus
//
//	@Summary		Обновление статуса посылки
//	@Description	Обновляет статус посылки и добавляет запись в историю
//	@Tags			Parcel
//	@Produce		json
//	@Param			id		query	int						true	"ID"
//	@Param			request	body	dto.ChangeStatusRequest	true	"Change status request"
//	@Success		204		"No Content"
//	@Failure		400		{object}	response.ErrorResponse	"Bad request"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		404		{object}	response.ErrorResponse	"Not found"
//	@Failure		405		{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/parcels/status [patch]
func (h *ParcelHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parcelID, err := request.PositiveIntQuery(r, "id")
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req dto.ChangeStatusRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if req.Status == "" || req.Location == "" {
		response.Error(w, "status and location cannot be empty", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		response.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	if err = h.parcelService.ChangeStatus(parcelID, req.Status, req.Location, userID); err != nil {
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			response.Error(w, "cannot skip statuses", http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrParcelNotFound) {
			response.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		response.Error(w, "failed to update parcel status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddPhoto
//
// @Summary Добавление фото к посылке
// @Description Загружает фото посылки и сохраняет путь к файлу
// @Tags Parcel
// @Accept multipart/form-data
// @Produce json
// @Param id query int true "Parcel ID"
// @Param file formData file true "Photo file"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Not found"
// @Failure 405 {object} response.ErrorResponse "Method not allowed"
// @Failure 413 {object} response.ErrorResponse "File too large"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /parcels/photos [post]
func (h *ParcelHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parcelID, err := request.PositiveIntQuery(r, "id")
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	filePath, err := saveUploadedFile(r)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.parcelService.AddPhoto(parcelID, filePath); err != nil {
		if errors.Is(err, repository.ErrParcelNotFound) {
			response.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		response.Error(w, "failed to create parcel photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func saveUploadedFile(r *http.Request) (string, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("read file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	if header.Filename == "" {
		return "", fmt.Errorf("empty file name")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("unsupported file extension: %s", ext)
	}

	if err = os.MkdirAll("uploads", 0755); err != nil {
		return "", fmt.Errorf("make dir: %w", err)
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join("uploads", fileName)
	filePath = filepath.ToSlash(filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("file copy: %w", err)
	}

	return filePath, nil
}

// Archive
//
//	@Summary		Архивирует посылку
//	@Description	Отправляет посылку в архив
//	@Tags			Parcel
//	@Produce		json
//	@Param			id	query	int	true	"ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	response.ErrorResponse	"Bad request"
//	@Failure		401	{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		404	{object}	response.ErrorResponse	"Not found"
//	@Failure		405	{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/parcels/archive [patch]
func (h *ParcelHandler) Archive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parcelID, err := request.PositiveIntQuery(r, "id")
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		response.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	if err = h.parcelService.Archive(parcelID, userID); err != nil {
		if errors.Is(err, repository.ErrParcelNotFound) {
			response.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrParcelNotDelivered) {
			response.Error(w, "parcel not delivered", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrParcelAlreadyArchived) {
			response.Error(w, "parcel already archived", http.StatusBadRequest)
			return
		}
		response.Error(w, "archive parcel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List
//
//	@Summary		Получение списка посылок
//	@Description	Возвращает список посылок с пагинацией и фильтром по статусу
//	@Tags			Parcel
//	@Produce		json
//	@Param			status	query		string	false	"Parcel status filter"
//	@Param			page	query		int		false	"Page number"
//	@Param			limit	query		int		false	"Items per page"
//	@Success		200		{object}	dto.ListParcelResponse
//	@Failure		400		{object}	response.ErrorResponse	"Bad request"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	response.ErrorResponse	"Forbidden"
//	@Failure		405		{object}	response.ErrorResponse	"Method not allowed"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Security		BearerAuth
//	@Router			/parcels [get]
func (h *ParcelHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	status := domain.Status(q.Get("status"))
	pageStr := q.Get("page")
	limitStr := q.Get("limit")

	var page int
	var limit int
	var err error

	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			response.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
	}

	if limitStr == "" {
		limit = 20
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			response.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
	}

	parcels, total, err := h.parcelService.List(status, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLimit) {
			response.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrInvalidPage) {
			response.Error(w, "invalid page", http.StatusBadRequest)
			return
		}

		response.Error(w, "get parcel list", http.StatusInternalServerError)
		return
	}

	listItemResp := make([]dto.ParcelListItemResponse, 0, len(parcels))
	for _, parcel := range parcels {
		listItemResp = append(listItemResp, dto.ParcelListItemResponse{
			ID:              parcel.ID,
			TrackNumber:     parcel.TrackNumber,
			ItemName:        parcel.ItemName,
			RecipientName:   parcel.RecipientName,
			CurrentStatus:   parcel.CurrentStatus,
			CurrentLocation: parcel.CurrentLocation,
		})
	}

	listResp := dto.ListParcelResponse{
		Items: listItemResp,
		Page:  page,
		Limit: limit,
		Total: total,
	}

	response.JSON(w, http.StatusOK, listResp)
}

func (h *ParcelHandler) Parcels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.CreateParcel(w, r)
	default:
		response.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
