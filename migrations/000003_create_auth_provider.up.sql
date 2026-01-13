-- Auth Provider lookup table
CREATE TABLE IF NOT EXISTS "AuthProvider" (
    "id" INT PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL UNIQUE,
    "description" VARCHAR(255),
    "isActive" BOOLEAN DEFAULT TRUE,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default auth providers
INSERT INTO "AuthProvider" ("id", "name", "description") VALUES
    (1, 'PASSWORD', 'Email/Phone/Username with password authentication'),
    (2, 'GOOGLE', 'Google OAuth authentication');
