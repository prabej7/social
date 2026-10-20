ALTER TABLE user_invitation
ALTER COLUMN token TYPE bytea USING decode(token, 'hex');