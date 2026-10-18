CREATE TABLE IF NOT EXISTS followers (
    user_id SERIAL REFERENCES users(id) ON DELETE CASCADE,
    follower_id SERIAL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, follower_id),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

