-- name: CreateUser :one
INSERT INTO users (id, tenant_id, email, hashed_password, profile)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE tenant_id = $1 AND email = $2;

-- name: GetUsersByTenant :many
SELECT * FROM users WHERE tenant_id = $1 ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users 
SET email = $3, hashed_password = $4, profile = $5
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1 AND tenant_id = $2;

-- name: GetUserCountByTenant :one
SELECT COUNT(*) FROM users WHERE tenant_id = $1; 

-- 用户Refresh Token表
-- name: CreateRefreshToken :exec
INSERT INTO user_refresh_tokens (user_id, token_hash, expires_at, created_at, client_ip, user_agent)
VALUES ($1, $2, $3, NOW(), $4, $5);

-- name: GetRefreshToken :one
SELECT * FROM user_refresh_tokens WHERE user_id = $1 AND token_hash = $2 AND expires_at > NOW();

-- name: DeleteRefreshToken :exec
DELETE FROM user_refresh_tokens WHERE user_id = $1 AND token_hash = $2;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM user_refresh_tokens WHERE user_id = $1;

-- 表结构建议（请在migrations中建表）
-- CREATE TABLE user_refresh_tokens (
--   id SERIAL PRIMARY KEY,
--   user_id VARCHAR NOT NULL,
--   token_hash VARCHAR NOT NULL,
--   expires_at TIMESTAMP NOT NULL,
--   created_at TIMESTAMP NOT NULL,
--   client_ip VARCHAR,
--   user_agent VARCHAR
-- ); 