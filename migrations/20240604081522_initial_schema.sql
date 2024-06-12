-- +goose Up
-- +goose StatementBegin

CREATE TABLE users
(
    user_id       SERIAL PRIMARY KEY,
    username      VARCHAR(50) UNIQUE  NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE friendships
(
    user_id    INTEGER REFERENCES users (user_id),
    friend_id  INTEGER REFERENCES users (user_id),
    status     VARCHAR(20) CHECK (status IN ('pending', 'accepted', 'blocked')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, friend_id)
);

CREATE TABLE conversations
(
    conversation_id SERIAL PRIMARY KEY,
    user1_id        INTEGER REFERENCES users (user_id),
    user2_id        INTEGER REFERENCES users (user_id),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user1_id, user2_id)
);


CREATE TABLE private_messages
(
    message_id      SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES conversations (conversation_id),
    sender_id       INTEGER REFERENCES users (user_id),
    message_text    TEXT NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE channels
(
    channel_id  SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    is_public   BOOLEAN   DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE channel_members
(
    channel_id INTEGER REFERENCES channels (channel_id),
    user_id    INTEGER REFERENCES users (user_id),
    role       VARCHAR(20) CHECK (role IN ('member', 'admin', 'owner')),
    joined_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (channel_id, user_id)
);

CREATE TABLE chat_rooms
(
    chat_room_id SERIAL PRIMARY KEY,
    channel_id   INTEGER REFERENCES channels (channel_id),
    name         VARCHAR(100) NOT NULL,
    description  TEXT,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chat_room_roles
(
    chat_room_id INTEGER REFERENCES chat_rooms (chat_room_id),
    role         VARCHAR(20) CHECK (role IN ('member', 'admin', 'owner')),
    PRIMARY KEY (chat_room_id, role)
);

CREATE TABLE channel_messages
(
    message_id   SERIAL PRIMARY KEY,
    chat_room_id INTEGER REFERENCES chat_rooms (chat_room_id),
    user_id      INTEGER REFERENCES users (user_id),
    message_text TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_friendships_user_id ON friendships(user_id);
CREATE INDEX idx_friendships_friend_id ON friendships(friend_id);
CREATE INDEX idx_private_messages_sender_id ON private_messages(sender_id);
CREATE INDEX idx_private_messages_conversation_id ON private_messages(conversation_id);
CREATE INDEX idx_channel_members_channel_id ON channel_members(channel_id);
CREATE INDEX idx_channel_members_user_id ON channel_members(user_id);
CREATE INDEX idx_channel_messages_chat_room_id ON channel_messages(chat_room_id);
CREATE INDEX idx_channel_messages_user_id ON channel_messages(user_id);

ALTER TABLE friendships ADD CONSTRAINT unique_friendship UNIQUE (user_id, friend_id);
ALTER TABLE channel_members ADD CONSTRAINT unique_channel_membership UNIQUE (channel_id, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_channel_messages_user_id;
DROP INDEX IF EXISTS idx_channel_messages_chat_room_id;
DROP INDEX IF EXISTS idx_channel_members_user_id;
DROP INDEX IF EXISTS idx_channel_members_channel_id;
DROP INDEX IF EXISTS idx_private_messages_sender_id;
DROP INDEX IF EXISTS idx_private_messages_conversation_id;
DROP INDEX IF EXISTS idx_friendships_friend_id;
DROP INDEX IF EXISTS idx_friendships_user_id;
DROP INDEX IF EXISTS idx_users_username;

DROP TABLE IF EXISTS chat_room_roles;
DROP TABLE IF EXISTS chat_rooms;
DROP TABLE IF EXISTS channel_members;
DROP TABLE IF EXISTS channels;
DROP TABLE IF EXISTS private_messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS friendships;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
