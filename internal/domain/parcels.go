package domain

import "time"

type Parcel struct {
	ID               int       `db:"id"`
	TrackNumber      string    `db:"track_number"`
	ItemName         string    `db:"item_name"`
	RecipientName    string    `db:"recipient_name"`
	RecipientPhone   string    `db:"recipient_phone"`
	RecipientAddress string    `db:"recipient_address"`
	CurrentStatus    Status    `db:"current_status"`
	CurrentLocation  string    `db:"current_location"`
	IsArchived       bool      `db:"is_archived"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
