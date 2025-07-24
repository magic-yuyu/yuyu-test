-- name: CreateTenant :one
INSERT INTO tenants (id, name, api_secret_key_hash, api_public_key)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants WHERE id = $1;

-- name: GetTenantByPublicKey :one
SELECT * FROM tenants WHERE api_public_key = $1;

-- name: GetTenantBySecretKeyHash :one
SELECT * FROM tenants WHERE api_secret_key_hash = $1;

-- name: ListTenants :many
SELECT * FROM tenants ORDER BY created_at DESC;

-- name: UpdateTenant :one
UPDATE tenants 
SET name = $2, api_secret_key_hash = $3, api_public_key = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTenant :exec
DELETE FROM tenants WHERE id = $1; 