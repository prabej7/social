ALTER TABLE user_invitation
ALTER COLUMN token TYPE VARCHAR(255) USING encode(token, 'hex');