CREATE TABLE
    IF NOT EXISTS core.user_auth_providers (
        user_id UUID NOT NULL,
        app_id UUID NOT NULL,
        -- This field is used to identify the authentication provider
        -- local | google | github | facebook | etc.
        provider VARCHAR(50) NOT NULL,
        -- The provider_user_id is possible to be null when provider is local
        provider_user_id VARCHAR(255) NULL,
        -- Password field is optional, only required for local authentication
        password VARCHAR(255) NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT pk_users_authentication PRIMARY KEY (user_id, app_id, provider),
        CONSTRAINT fk_user_auth_user FOREIGN KEY (user_id, app_id) REFERENCES core.users (id, app_id) ON DELETE CASCADE,
        CONSTRAINT fk_user_auth_app FOREIGN KEY (app_id) REFERENCES core.apps (id) ON DELETE CASCADE
    );