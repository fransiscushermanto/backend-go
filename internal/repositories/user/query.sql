-- name: StoreUser :exec
INSERT INTO core.users (id, app_id, name, email, is_email_verified, email_verified_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: StoreUserAuthProvider :exec
INSERT INTO core.user_auth_providers (user_id, app_id, provider, provider_user_id, password) 
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserAuthenticationByProvider :one
SELECT user_id, app_id, provider, provider_user_id, password, created_at, updated_at 
FROM core.user_auth_providers 
WHERE app_id = $1 AND user_id = $2 AND provider = $3;

-- name: GetAllUsers :many
SELECT id, app_id, name, email 
FROM core.users 
ORDER BY created_at DESC;

-- name: GetAllUsersByAppID :many
SELECT id, name, email 
FROM core.users 
WHERE app_id=$1 
ORDER BY created_at DESC;

-- name: GetAppUserByID :one
SELECT id, app_id, name, email, is_email_verified, email_verified_at 
FROM core.users 
WHERE app_id = $1 AND id = $2;

-- name: GetUserByEmail :one
SELECT id, app_id, name, email, is_email_verified, email_verified_at, created_at, updated_at 
FROM core.users 
WHERE app_id = $1 AND email = $2;