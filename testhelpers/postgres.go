package testhelpers

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresContainer(t *testing.T) *sqlx.DB {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("не удалось запустить postgres контейнер: %v", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("ошибка остановки контейнера: %v", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("не удалось получить строку подключения: %v", err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("не удалось подключиться к тестовой БД: %v", err)
	}

	applyMigrations(t, db)

	return db
}

func applyMigrations(t *testing.T, db *sqlx.DB) {
	t.Helper()
	schema := `
		CREATE TYPE users_roles AS ENUM ('manager', 'admin');
		
		CREATE TYPE audit_action AS ENUM (
			'create_parcel',
			'change_status',
			'upload_photo',
			'archive_parcel',
			'deactivate_user'
		);
		
		CREATE TYPE audit_entity_type AS ENUM ('user', 'parcel');
		
		CREATE TABLE IF NOT EXISTS statuses (
			id SERIAL PRIMARY KEY,
			status TEXT UNIQUE NOT NULL
		);
		
		INSERT INTO statuses (status)
		VALUES
			('CREATED'),
			('PURCHASED'),
			('WAREHOUSE'),
			('IN_TRANSIT'),
			('CUSTOMS'),
			('ARRIVED'),
			('DELIVERED')
		ON CONFLICT (status) DO NOTHING;
		
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role users_roles NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
		
		CREATE TABLE IF NOT EXISTS parcels (
			id SERIAL PRIMARY KEY,
			track_number TEXT UNIQUE NOT NULL,
			item_name TEXT NOT NULL,
			recipient_name TEXT NOT NULL,
			recipient_phone TEXT NOT NULL,
			recipient_address TEXT NOT NULL,
			current_status INTEGER REFERENCES statuses(id) NOT NULL,
			current_location TEXT,
			is_archived BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
		
		CREATE TABLE IF NOT EXISTS parcel_photos (
			id SERIAL PRIMARY KEY,
			parcel_id INTEGER REFERENCES parcels(id) NOT NULL,
			file_path TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
		
		CREATE TABLE IF NOT EXISTS parcel_status_history (
			id SERIAL PRIMARY KEY,
			parcel_id INTEGER REFERENCES parcels(id) NOT NULL,
			old_status INTEGER REFERENCES statuses(id),
			new_status INTEGER REFERENCES statuses(id) NOT NULL,
			location TEXT,
			changed_by INTEGER REFERENCES users(id) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
		
		CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) NOT NULL,
			action audit_action NOT NULL,
			old_value TEXT,
			new_value TEXT,
			entity_type audit_entity_type NOT NULL,
			entity_id INTEGER NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);

		CREATE INDEX idx_parcels_current_status
		ON parcels(current_status);
	
		CREATE INDEX idx_parcel_photos_parcel_id
			ON parcel_photos(parcel_id);
		
		CREATE INDEX idx_parcel_status_history_parcel_id_created_at
			ON parcel_status_history(parcel_id, created_at);
		
		CREATE INDEX idx_audit_logs_user_id_created_at
			ON audit_logs(user_id, created_at);
		
		CREATE INDEX idx_audit_logs_created_at
			ON audit_logs(created_at);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("ошибка миграции: %v", err)
	}
}
