-- migrate:up
CREATE TABLE user_password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_password_reset_tokens_token ON user_password_reset_tokens(token);

-- migrate:down
DROP TABLE IF EXISTS user_password_reset_tokens;
