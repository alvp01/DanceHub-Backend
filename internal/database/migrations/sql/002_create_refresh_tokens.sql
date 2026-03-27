-- migrations/002_create_refresh_tokens.sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    academy_id UUID NOT NULL REFERENCES academies (id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE, -- se guarda hasheado
    expires_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL,
        created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW (),
        revoked_at TIMESTAMP
    WITH
        TIME ZONE -- NULL = activo
);

CREATE INDEX idx_refresh_tokens_academy_id ON refresh_tokens (academy_id);

CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens (token_hash);