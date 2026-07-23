package service

import (
	"delivery-tracker/internal/domain"
)

func testParcel(status domain.Status, isArchived bool) *domain.Parcel {
	return &domain.Parcel{
		ID:              1,
		TrackNumber:     "Q4P405SHH8EG",
		ItemName:        "Reina Petite Ring",
		RecipientName:   "Иван",
		CurrentStatus:   status,
		CurrentLocation: "Berlin",
		IsArchived:      isArchived,
	}
}

func testParcelDetails(trackNumber string) *domain.ParcelDetails {
	return &domain.ParcelDetails{
		TrackNumber: trackNumber,
		ItemName:    "Reina Petite Ring",
	}
}

func testPhotos() []domain.ParcelPhoto {
	return []domain.ParcelPhoto{
		{ID: 1, ParcelID: 1, FilePath: "/uploads/__1231132.jpg"},
	}
}

func testHistory() []domain.ParcelStatusHistory {
	return []domain.ParcelStatusHistory{
		{
			ID:        1,
			ParcelID:  1,
			OldStatus: nil,
			NewStatus: domain.StatusPurchased,
			Location:  "loc",
			ChangedBy: 1,
		},
	}
}

func testListParcel(status domain.Status, isArchived bool) []domain.Parcel {
	return []domain.Parcel{
		{
			ID:              1,
			TrackNumber:     "Q4P405SHH8EG",
			ItemName:        "Reina Petite Ring",
			RecipientName:   "Иван",
			CurrentStatus:   status,
			CurrentLocation: "Berlin",
			IsArchived:      isArchived,
		},
		{
			ID:              2,
			TrackNumber:     "Y0P408SJH8ER",
			ItemName:        "Reina Petite Ring",
			RecipientName:   "Андрей",
			CurrentStatus:   status,
			CurrentLocation: "Berlin",
			IsArchived:      isArchived,
		},
	}
}
