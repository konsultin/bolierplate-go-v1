-- Create user table
CREATE TABLE IF NOT EXISTS "User" (
    "id" BIGSERIAL PRIMARY KEY,
    "xid" VARCHAR(255) NOT NULL UNIQUE,
    "username" VARCHAR(100) UNIQUE,
    "fullName" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE,
    "phone" VARCHAR(20) UNIQUE,
    "age" VARCHAR(10),
    "avatar" VARCHAR(255),
    "statusId" INT NOT NULL DEFAULT 1,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modifiedBy" JSONB,
    "version" BIGINT NOT NULL DEFAULT 1,
    "metadata" JSONB DEFAULT '{}'
);

-- Create indexes for user table
CREATE INDEX IF NOT EXISTS idx_user_username ON "User"("username");

-- User credential table for flexible authentication
CREATE TABLE IF NOT EXISTS "UserCredential" (
    "id" BIGSERIAL PRIMARY KEY,
    "userId" BIGINT NOT NULL REFERENCES "User"("id") ON DELETE CASCADE,
    "authProviderId" INT NOT NULL DEFAULT 1,
    "credentialKey" VARCHAR(255) NOT NULL,
    "credentialSecret" VARCHAR(255),
    "isVerified" BOOLEAN DEFAULT FALSE,
    "verifiedAt" TIMESTAMP,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("authProviderId", "credentialKey")
);

CREATE INDEX IF NOT EXISTS idx_user_credential_user_id ON "UserCredential"("userId");
CREATE INDEX IF NOT EXISTS idx_user_credential_key ON "UserCredential"("credentialKey");

COMMENT ON COLUMN "UserCredential"."authProviderId" IS '1=PASSWORD, 2=GOOGLE';
COMMENT ON COLUMN "UserCredential"."credentialKey" IS 'email/phone/username for PASSWORD, provider_user_id for OAuth';
COMMENT ON COLUMN "UserCredential"."credentialSecret" IS 'password_hash for PASSWORD, null for OAuth';

