CREATE TABLE IF NOT EXISTS users (
    id	BIGSERIAL NOT NULL PRIMARY KEY,
    username	TEXT NOT NULL UNIQUE,
    password_hash	BYTEA NOT NULL,
    online_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE friendships (
    id	BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    friend_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE messages (
    id	BIGSERIAL NOT NULL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES users(id),
    recipient_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);