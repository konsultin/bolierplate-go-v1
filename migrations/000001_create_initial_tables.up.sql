-- Create privilege table
CREATE TABLE IF NOT EXISTS "Privilege" (
    "id" BIGSERIAL PRIMARY KEY,
    "xid" VARCHAR(255) NOT NULL UNIQUE,
    "name" VARCHAR(255) NOT NULL,
    "exposed" BOOLEAN NOT NULL DEFAULT false,
    "sort" INT NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modifiedBy" JSONB,
    "version" BIGINT NOT NULL DEFAULT 1,
    "metadata" JSONB DEFAULT '{}'
);

-- Create role table
CREATE TABLE IF NOT EXISTS "Role" (
    "id" SERIAL PRIMARY KEY,
    "xid" VARCHAR(255) NOT NULL UNIQUE,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "roleTypeId" INT NOT NULL DEFAULT 1,
    "statusId" INT NOT NULL DEFAULT 1,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modifiedBy" JSONB,
    "version" BIGINT NOT NULL DEFAULT 1,
    "metadata" JSONB DEFAULT '{}'
);

-- Create client_auth table
CREATE TABLE IF NOT EXISTS "ClientAuth" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "clientId" VARCHAR(255) NOT NULL UNIQUE,
    "clientTypeId" INT NOT NULL,
    "options" JSONB NOT NULL DEFAULT '{"clientSecret": "", "tokenLifetime": 2592000}',
    "statusId" INT NOT NULL DEFAULT 1,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modifiedBy" JSONB,
    "version" BIGINT NOT NULL DEFAULT 1,
    "metadata" JSONB DEFAULT '{}'
);

-- Create role_privilege table (junction table)
CREATE TABLE IF NOT EXISTS "RolePrivilege" (
    "id" BIGSERIAL PRIMARY KEY,
    "roleId" INT NOT NULL,
    "privilegeId" BIGINT NOT NULL,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modifiedBy" JSONB,
    "version" BIGINT NOT NULL DEFAULT 1,
    "metadata" JSONB DEFAULT '{}',
    CONSTRAINT fk_role_privilege_role FOREIGN KEY ("roleId") REFERENCES "Role"("id") ON DELETE CASCADE,
    CONSTRAINT fk_role_privilege_privilege FOREIGN KEY ("privilegeId") REFERENCES "Privilege"("id") ON DELETE CASCADE,
    CONSTRAINT unique_role_privilege UNIQUE ("roleId", "privilegeId")
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_privilege_xid ON "Privilege"("xid");
CREATE INDEX IF NOT EXISTS idx_privilege_exposed ON "Privilege"("exposed");
CREATE INDEX IF NOT EXISTS idx_role_xid ON "Role"("xid");
CREATE INDEX IF NOT EXISTS idx_role_status ON "Role"("statusId");
CREATE INDEX IF NOT EXISTS idx_client_auth_client_id ON "ClientAuth"("clientId");
CREATE INDEX IF NOT EXISTS idx_client_auth_status ON "ClientAuth"("statusId");
CREATE INDEX IF NOT EXISTS idx_role_privilege_role_id ON "RolePrivilege"("roleId");
CREATE INDEX IF NOT EXISTS idx_role_privilege_privilege_id ON "RolePrivilege"("privilegeId");
