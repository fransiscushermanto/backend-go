ALTER TABLE core.refresh_tokens
DROP COLUMN updated_at;

ALTER TABLE core.blacklist_tokens
DROP COLUMN updated_at;

ALTER TABLE core.app_api_keys
DROP COLUMN updated_at;

ALTER TABLE core.reset_password_tokens
DROP COLUMN updated_at;