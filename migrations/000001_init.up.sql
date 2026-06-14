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