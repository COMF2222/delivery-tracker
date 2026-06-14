package domain

import "time"

type AuditLog struct {
	ID         int        `db:"id"`
	UserID     int        `db:"user_id"`
	Action     Action     `db:"action"`
	OldValue   string     `db:"old_value"`
	NewValue   string     `db:"new_value"`
	EntityType EntityType `db:"entity_type"`
	EntityID   int        `db:"entity_id"`
	CreatedAt  time.Time  `db:"created_at"`
}

type EntityType string

type Action string

const (
	EntityTypeUser   EntityType = "user"
	EntityTypeParcel EntityType = "parcel"
)

const (
	ActionCreateParcel   Action = "create_parcel"
	ActionChangeStatus   Action = "change_status"
	ActionUploadPhoto    Action = "upload_photo"
	ActionArchiveParcel  Action = "archive_parcel"
	ActionDeactivateUser Action = "deactivate_user"
)
