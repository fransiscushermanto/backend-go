ALTER TABLE core.apps
ALTER COLUMN created_at
SET
    NOT NULL,
ALTER COLUMN updated_at
SET
    NOT NULL;

ALTER TABLE core.users
ALTER COLUMN created_at
SET
    NOT NULL,
ALTER COLUMN updated_at
SET
    NOT NULL;

ALTER TABLE core.user_auth_providers
ALTER COLUMN created_at
SET
    NOT NULL,
ALTER COLUMN updated_at
SET
    NOT NULL;

ALTER TABLE core.refresh_tokens
ALTER COLUMN created_at
SET
    NOT NULL;

ALTER TABLE core.blacklist_tokens
ALTER COLUMN blacklisted_at
SET
    NOT NULL;

ALTER TABLE core.app_api_keys
ALTER COLUMN created_at
SET
    NOT NULL;

ALTER TABLE core.reset_password_tokens
ALTER COLUMN created_at
SET
    NOT NULL;