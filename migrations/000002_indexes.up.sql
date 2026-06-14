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