-- 用户refresh token表
CREATE TABLE IF NOT EXISTS user_refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    token_hash VARCHAR NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    client_ip VARCHAR,
    user_agent VARCHAR
);
CREATE INDEX IF NOT EXISTS idx_user_refresh_tokens_user_id ON user_refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_user_refresh_tokens_token_hash ON user_refresh_tokens(token_hash); 