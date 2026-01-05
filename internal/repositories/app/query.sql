-- name: GetAppByID :one
SELECT id, name, created_at, updated_at FROM core.apps WHERE id = $1;

-- name: GetAllApps :many
SELECT id, name FROM core.apps ORDER BY created_at DESC;

-- name: StoreApp :exec
INSERT INTO core.apps (id, name) 
VALUES ($1, $2);

-- name: StoreAppApiKey :exec
INSERT INTO core.app_api_keys (id, app_id, key_hash, is_active) 
VALUES ($1, $2, $3, $4);

-- name: RevokeActiveAppApiKeys :execrows
UPDATE core.app_api_keys 
SET is_active = false, revoked_at = $2 
WHERE app_id = $1 AND is_active = true;

-- name: LockAppForUpdate :one
SELECT id FROM core.apps WHERE id = $1 FOR UPDATE;