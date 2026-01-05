-- name: StoreRefreshToken :exec
INSERT INTO core.refresh_tokens (jti, user_id, app_id, token, expires_at, is_active) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetRefreshTokenByJTI :one
SELECT jti, user_id, app_id, token, expires_at, is_active, created_at 
FROM core.refresh_tokens 
WHERE app_id = $1 AND jti = $2 
ORDER BY created_at;

-- name: GetUserActiveRefreshTokensByUserID :many
SELECT * FROM core.refresh_tokens
WHERE app_id = $1 AND user_id = $2 AND is_active = true
ORDER BY created_at DESC;

-- name: GetUserActiveRefreshTokensByJTI :many
SELECT * FROM core.refresh_tokens
WHERE app_id = $1 AND jti = $2 AND is_active = true
ORDER BY created_at DESC;

-- name: RevokeRefreshTokens :exec
UPDATE core.refresh_tokens 
SET is_active = false 
WHERE app_id = $1 AND user_id = $2 AND jti = ANY(@jtis::text[]);