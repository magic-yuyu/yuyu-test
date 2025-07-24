-- name: CreateInternalClient :one
INSERT INTO internal_clients (client_id, client_secret_hash, service_name, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetInternalClient :one
SELECT * FROM internal_clients WHERE client_id = $1 AND is_active = true;

-- name: GetInternalClientByID :one
SELECT * FROM internal_clients WHERE client_id = $1;

-- name: ListInternalClients :many
SELECT * FROM internal_clients WHERE is_active = true ORDER BY created_at DESC;

-- name: UpdateInternalClient :one
UPDATE internal_clients 
SET service_name = $2, description = $3, updated_at = CURRENT_TIMESTAMP
WHERE client_id = $1 AND is_active = true
RETURNING *;

-- name: DeactivateInternalClient :exec
UPDATE internal_clients SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE client_id = $1;

-- name: ActivateInternalClient :exec
UPDATE internal_clients SET is_active = true, updated_at = CURRENT_TIMESTAMP WHERE client_id = $1;

-- name: DeleteInternalClient :exec
DELETE FROM internal_clients WHERE client_id = $1;

-- name: GetClientScopes :many
SELECT s.id, s.scope_name, s.description, cs.granted_at, cs.granted_by
FROM client_scopes cs
JOIN scopes s ON cs.scope_id = s.id
WHERE cs.client_id = $1 AND s.is_active = true
ORDER BY s.scope_name;

-- name: GrantScopeToClient :exec
INSERT INTO client_scopes (client_id, scope_id, granted_by)
VALUES ($1, $2, $3)
ON CONFLICT (client_id, scope_id) DO NOTHING;

-- name: RevokeScopeFromClient :exec
DELETE FROM client_scopes WHERE client_id = $1 AND scope_id = $2;

-- name: CheckClientHasScope :one
SELECT EXISTS(
    SELECT 1 FROM client_scopes cs
    JOIN scopes s ON cs.scope_id = s.id
    WHERE cs.client_id = $1 AND s.scope_name = $2 AND s.is_active = true
) as has_scope;

-- name: ListAllScopes :many
SELECT * FROM scopes WHERE is_active = true ORDER BY scope_name;

-- name: GetScopeByName :one
SELECT * FROM scopes WHERE scope_name = $1 AND is_active = true;

-- name: CreateScope :one
INSERT INTO scopes (scope_name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateScope :one
UPDATE scopes 
SET description = $2
WHERE scope_name = $1 AND is_active = true
RETURNING *;

-- name: DeactivateScope :exec
UPDATE scopes SET is_active = false WHERE scope_name = $1;

-- name: LogServiceAccess :exec
INSERT INTO service_access_logs (
    client_id, endpoint, method, status_code, response_time_ms, 
    ip_address, user_agent, request_body, response_body
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetServiceAccessLogs :many
SELECT * FROM service_access_logs 
WHERE client_id = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: StoreServiceToken :exec
INSERT INTO service_tokens (client_id, token_hash, scopes, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetServiceToken :one
SELECT * FROM service_tokens 
WHERE token_hash = $1 AND expires_at > CURRENT_TIMESTAMP AND is_revoked = false;

-- name: RevokeServiceToken :exec
UPDATE service_tokens SET is_revoked = true WHERE token_hash = $1;

-- name: CleanupExpiredTokens :exec
DELETE FROM service_tokens WHERE expires_at < CURRENT_TIMESTAMP;

-- name: GetClientStatistics :one
SELECT 
    COUNT(*) as total_requests,
    AVG(response_time_ms) as avg_response_time,
    COUNT(CASE WHEN status_code >= 400 THEN 1 END) as error_count
FROM service_access_logs 
WHERE client_id = $1 AND created_at >= $2; 