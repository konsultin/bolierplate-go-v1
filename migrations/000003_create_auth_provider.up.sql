-- Auth Provider lookup table
CREATE TABLE IF NOT EXISTS auth_provider (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default auth providers
INSERT INTO auth_provider (id, name, description) VALUES
    (1, 'PASSWORD', 'Email/Phone/Username with password authentication'),
    (2, 'GOOGLE', 'Google OAuth authentication');
