package domain

import "time"

type ParcelPhoto struct {
	ID        int       `db:"id"`
	ParcelID  int       `db:"parcel_id"`
	FilePath  string    `db:"file_path"`
	CreatedAt time.Time `db:"created_at"`
}
