-- 删除索引
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_tenants_public_key;
DROP INDEX IF EXISTS idx_tenants_secret_key_hash;
DROP INDEX IF EXISTS idx_internal_clients_service_name;
DROP INDEX IF EXISTS idx_scopes_scope_name;
DROP INDEX IF EXISTS idx_client_scopes_client_id;
DROP INDEX IF EXISTS idx_service_access_logs_client_id;
DROP INDEX IF EXISTS idx_service_tokens_client_id;

-- 删除表
DROP TABLE IF EXISTS service_access_logs;
DROP TABLE IF EXISTS service_tokens;
DROP TABLE IF EXISTS client_scopes;
DROP TABLE IF EXISTS scopes;
DROP TABLE IF EXISTS internal_clients;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants; 