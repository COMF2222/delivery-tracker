package domain

type ParcelDetails struct {
	TrackNumber     string
	ItemName        string
	RecipientName   string
	CurrentStatus   Status
	CurrentLocation string
	History         []ParcelStatusHistory
	Photos          []ParcelPhoto
}
