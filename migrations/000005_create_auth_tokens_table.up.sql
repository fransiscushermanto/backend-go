CREATE TABLE
    IF NOT EXISTS core.refresh_tokens (
        jti VARCHAR(255) NOT NULL PRIMARY KEY,
        user_id UUID NOT NULL,
        app_id UUID NOT NULL,
        token VARCHAR(255) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMPTZ NOT NULL,
        CONSTRAINT fk_refresh_token_user FOREIGN KEY (user_id, app_id) REFERENCES core.users (id, app_id) ON DELETE CASCADE,
        CONSTRAINT fk_refresh_token_app FOREIGN KEY (app_id) REFERENCES core.apps (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS core.blacklist_tokens (
        jti VARCHAR(255) NOT NULL PRIMARY KEY,
        user_id UUID NOT NULL,
        app_id UUID NOT NULL,
        token VARCHAR(255) NOT NULL,
        blacklisted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        -- 'logout' or 'revoke'
        reason VARCHAR(100) NULL,
        CONSTRAINT fk_blacklist_token_user FOREIGN KEY (user_id, app_id) REFERENCES core.users (id, app_id) ON DELETE CASCADE,
        CONSTRAINT fk_blacklist_token_app FOREIGN KEY (app_id) REFERENCES core.apps (id) ON DELETE CASCADE
    )