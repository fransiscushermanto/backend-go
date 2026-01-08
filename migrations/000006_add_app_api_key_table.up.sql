CREATE TABLE
    core.app_api_keys (
        id UUID PRIMARY KEY,
        app_id UUID NOT NULL REFERENCES core.apps (id) ON DELETE CASCADE,
        key_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        revoked_at TIMESTAMPTZ NULL DEFAULT NULL,
        last_used_at TIMESTAMPTZ NULL DEFAULT NULL,
        CONSTRAINT uk_app_api_keys UNIQUE (app_id, key_hash)
    )