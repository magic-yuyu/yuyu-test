-- 创建租户表
CREATE TABLE IF NOT EXISTS tenants (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_secret_key_hash VARCHAR(255) UNIQUE NOT NULL,
    api_public_key VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(255),
    profile JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(tenant_id, email)
);

-- 内部服务表
CREATE TABLE IF NOT EXISTS internal_clients (
    client_id VARCHAR(255) PRIMARY KEY,
    client_secret_hash VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 权限表
CREATE TABLE IF NOT EXISTS scopes (
    id SERIAL PRIMARY KEY,
    scope_name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 服务-权限关联表
CREATE TABLE IF NOT EXISTS client_scopes (
    client_id VARCHAR(255) NOT NULL REFERENCES internal_clients(client_id) ON DELETE CASCADE,
    scope_id INT NOT NULL REFERENCES scopes(id) ON DELETE CASCADE,
    granted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    granted_by VARCHAR(255),
    PRIMARY KEY (client_id, scope_id)
);

-- 服务访问日志表
CREATE TABLE IF NOT EXISTS service_access_logs (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES internal_clients(client_id) ON DELETE CASCADE,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INT NOT NULL,
    response_time_ms INT,
    ip_address INET,
    user_agent TEXT,
    request_body TEXT,
    response_body TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 服务令牌表（用于存储已颁发的JWT令牌信息）
CREATE TABLE IF NOT EXISTS service_tokens (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES internal_clients(client_id) ON DELETE CASCADE,
    token_hash VARCHAR(512) NOT NULL,
    scopes TEXT[] NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_tenants_public_key ON tenants(api_public_key);
CREATE INDEX IF NOT EXISTS idx_tenants_secret_key_hash ON tenants(api_secret_key_hash);
CREATE INDEX IF NOT EXISTS idx_internal_clients_active ON internal_clients(is_active);
CREATE INDEX IF NOT EXISTS idx_scopes_active ON scopes(is_active);
CREATE INDEX IF NOT EXISTS idx_service_access_logs_client_id ON service_access_logs(client_id);
CREATE INDEX IF NOT EXISTS idx_service_access_logs_created_at ON service_access_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_service_tokens_client_id ON service_tokens(client_id);
CREATE INDEX IF NOT EXISTS idx_service_tokens_expires_at ON service_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_service_tokens_revoked ON service_tokens(is_revoked);

-- 插入默认权限
INSERT INTO scopes (scope_name, description) VALUES 
('user:read', '读取用户信息'),
('user:write', '创建和更新用户信息'),
('user:delete', '删除用户'),
('tenant:read', '读取租户信息'),
('tenant:write', '创建和更新租户信息'),
('tenant:delete', '删除租户'),
('auth:token', '生成认证令牌'),
('auth:validate', '验证令牌'),
('internal:admin', '内部服务管理权限')
ON CONFLICT (scope_name) DO NOTHING; 