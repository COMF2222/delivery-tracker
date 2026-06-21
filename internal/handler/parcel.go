package handler

import (
	"delivery-tracker/internal/contextkeys"
	"delivery-tracker/internal/domain"
	"delivery-tracker/internal/dto"
	"delivery-tracker/internal/repository"
	"delivery-tracker/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type ParcelHandler struct {
	parcelService *service.ParcelService
}

func NewParcelHandler(parcelService *service.ParcelService) *ParcelHandler {
	return &ParcelHandler{parcelService: parcelService}
}

func (h *ParcelHandler) CreateParcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req dto.CreateParcelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parcel := domain.Parcel{
		ItemName:         req.ItemName,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientAddress: req.RecipientAddress,
	}

	if err := h.parcelService.CreateParcel(&parcel); err != nil {
		http.Error(w, "failed to create parcel", http.StatusInternalServerError)
		return
	}

	resp := dto.CreateParcelResponse{
		ID:          parcel.ID,
		TrackNumber: parcel.TrackNumber,
	}

	responseJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		return
	}
}

func (h *ParcelHandler) GetByTrackNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackNumber := r.URL.Query().Get("track_number")
	if trackNumber == "" {
		http.Error(w, "track number cannot be empty", http.StatusBadRequest)
		return
	}

	parcel, err := h.parcelService.GetByTrackNumber(trackNumber)
	if err != nil {
		if errors.Is(err, repository.ErrParcelNotFound) {
			http.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get parcel by track number", http.StatusInternalServerError)
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
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		return
	}
}

func (h *ParcelHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "failed convert id to int", http.StatusBadRequest)
		return
	}
	if intId <= 0 {
		http.Error(w, "id must be positive", http.StatusBadRequest)
		return
	}

	var req dto.ChangeStatusRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if req.Status == "" || req.Location == "" {
		http.Error(w, "status and location cannot be empty", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	if err = h.parcelService.ChangeStatus(intId, req.Status, req.Location, userID); err != nil {
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			http.Error(w, "cannot skip statuses", http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrParcelNotFound) {
			http.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update parcel status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ParcelHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "failed convert id to int", http.StatusBadRequest)
		return
	}
	if intId <= 0 {
		http.Error(w, "id must be positive", http.StatusBadRequest)
		return
	}

	var req dto.AddPhotoRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if req.FilePath == "" {
		http.Error(w, "file path cannot be empty", http.StatusBadRequest)
		return
	}

	if err = h.parcelService.AddPhoto(intId, req.FilePath); err != nil {
		if errors.Is(err, repository.ErrParcelNotFound) {
			http.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to create parcel photo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ParcelHandler) Archive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id cannot be empty", http.StatusBadRequest)
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "failed convert id to int", http.StatusBadRequest)
		return
	}
	if intId <= 0 {
		http.Error(w, "id must be positive", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(contextkeys.UserID).(int)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	if err = h.parcelService.Archive(intId, userID); err != nil {
		if errors.Is(err, repository.ErrParcelNotFound) {
			http.Error(w, "parcel not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrParcelNotDelivered) {
			http.Error(w, "parcel not delivered", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrParcelAlreadyArchived) {
			http.Error(w, "parcel already archived", http.StatusBadRequest)
			return
		}
		http.Error(w, "archive parcel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ParcelHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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
			http.Error(w, "failed to atoi page", http.StatusBadRequest)
			return
		}
	}

	if limitStr == "" {
		limit = 20
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "failed to atoi limit", http.StatusBadRequest)
			return
		}
	}

	parcels, err := h.parcelService.List(status, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLimit) {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrInvalidPage) {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}

		http.Error(w, "get parcel list", http.StatusInternalServerError)
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
	}

	responseJSON, err := json.Marshal(listResp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		return
	}
}

func (h *ParcelHandler) Parcels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.CreateParcel(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
