CREATE TABLE
    IF NOT EXISTS core.reset_password_tokens (
        jti VARCHAR(255) NOT NULL PRIMARY KEY,
        user_id UUID NOT NULL,
        app_id UUID NOT NULL,
        token VARCHAR(255) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMPTZ NOT NULL,
        CONSTRAINT fk_reset_token_user FOREIGN KEY (user_id, app_id) REFERENCES core.users (id, app_id) ON DELETE CASCADE,
        CONSTRAINT fk_reset_token_app FOREIGN KEY (app_id) REFERENCES core.apps (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_reset_token_expiry ON core.reset_password_tokens (expires_at);

CREATE INDEX IF NOT EXISTS idx_reset_token_user ON core.reset_password_tokens (user_id);

CREATE INDEX IF NOT EXISTS idx_reset_token_user_app ON core.reset_password_tokens (user_id, app_id);