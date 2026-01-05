CREATE TABLE
    IF NOT EXISTS core.users (
        id UUID,
        app_id UUID REFERENCES core.apps (id) ON DELETE CASCADE,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
        email_verified_at TIMESTAMPTZ NULL DEFAULT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        -- Ensure that the combination of app_id and email is unique
        CONSTRAINT unique_email_per_app UNIQUE (app_id, email),
        CONSTRAINT pk_users PRIMARY KEY (id, app_id)
    );