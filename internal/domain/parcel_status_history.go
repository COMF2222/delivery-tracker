package domain

import "time"

type ParcelStatusHistory struct {
	ID        int       `db:"id"`
	ParcelID  int       `db:"parcel_id"`
	OldStatus *Status   `db:"old_status"`
	NewStatus Status    `db:"new_status"`
	Location  string    `db:"location"`
	ChangedBy int       `db:"changed_by"`
	CreatedAt time.Time `db:"created_at"`
}
