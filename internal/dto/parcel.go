package dto

import (
	"delivery-tracker/internal/domain"
	"fmt"
	"time"
)

type CreateParcelRequest struct {
	ItemName         string `json:"item_name"`
	RecipientName    string `json:"recipient_name"`
	RecipientPhone   string `json:"recipient_phone"`
	RecipientAddress string `json:"recipient_address"`
}

type CreateParcelResponse struct {
	ID          int    `json:"id"`
	TrackNumber string `json:"track_number"`
}
type GetParcelResponse struct {
	TrackNumber     string        `json:"track_number"`
	ItemName        string        `json:"item_name"`
	Recipient       string        `json:"recipient"`
	CurrentStatus   domain.Status `json:"current_status"`
	CurrentLocation string        `json:"current_location"`

	History []ParcelHistoryResponse `json:"history"`
	Photos  []ParcelPhotoResponse   `json:"photos"`
}

type ParcelPhotoResponse struct {
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}

type ParcelHistoryResponse struct {
	OldStatus *domain.Status `json:"old_status"`
	NewStatus domain.Status  `json:"new_status"`
	Location  string         `json:"location"`
	CreatedAt time.Time      `json:"created_at"`
}

type ChangeStatusRequest struct {
	Status   domain.Status `json:"status"`
	Location string        `json:"location"`
}

type AddPhotoRequest struct {
	FilePath string `json:"file_path"`
}

func (r CreateParcelRequest) Validate() error {
	if r.ItemName == "" || r.RecipientName == "" || r.RecipientPhone == "" || r.RecipientAddress == "" {
		return fmt.Errorf("failed to validate parcel request")
	}
	return nil
}
