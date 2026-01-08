CREATE TABLE
    core.refresh_tokens (
        jti VARCHAR(255) NOT NULL PRIMARY KEY,
        user_id UUID NOT NULL,
        app_id UUID NOT NULL,
        device_id VARCHAR(255) NOT NULL,
        device_name VARCHAR(255) NULL,
        token VARCHAR(255) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMPTZ NOT NULL,
        CONSTRAINT fk_refresh_token_user FOREIGN KEY (user_id, app_id) REFERENCES core.users (id, app_id) ON DELETE CASCADE,
        CONSTRAINT fk_refresh_token_app FOREIGN KEY (app_id) REFERENCES core.apps (id) ON DELETE CASCADE
    );

CREATE INDEX idx_user_sessions ON core.refresh_tokens (user_id, device_id);

CREATE TABLE
    core.blacklist_tokens (
        jti VARCHAR(255) NOT NULL PRIMARY KEY,
        token VARCHAR(255) NOT NULL,
        blacklisted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMPTZ NOT NULL,
        -- 'logout' or 'revoke'
        reason VARCHAR(100) NULL
    );

CREATE INDEX idx_blacklist_expiry ON core.blacklist_tokens (expires_at);