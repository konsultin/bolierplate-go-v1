-- Create user table
CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL PRIMARY KEY,
    xid VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    age VARCHAR(10),
    avatar VARCHAR(255),
    status_id INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_by JSONB,
    version BIGINT NOT NULL DEFAULT 1,
    metadata JSONB DEFAULT '{}'
);

-- Create indexes for user table
CREATE INDEX IF NOT EXISTS idx_user_username ON "user"(username);

-- User credential table for flexible authentication
CREATE TABLE IF NOT EXISTS user_credential (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    auth_provider_id INT NOT NULL DEFAULT 1,
    credential_key VARCHAR(255) NOT NULL,
    credential_secret VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(auth_provider_id, credential_key)
);

CREATE INDEX IF NOT EXISTS idx_user_credential_user_id ON user_credential(user_id);
CREATE INDEX IF NOT EXISTS idx_user_credential_key ON user_credential(credential_key);

COMMENT ON COLUMN user_credential.auth_provider_id IS '1=PASSWORD, 2=GOOGLE';
COMMENT ON COLUMN user_credential.credential_key IS 'email/phone/username for PASSWORD, provider_user_id for OAuth';
COMMENT ON COLUMN user_credential.credential_secret IS 'password_hash for PASSWORD, null for OAuth';

