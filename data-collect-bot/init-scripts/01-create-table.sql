-- Create tables based on ER diagram


-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    category_id VARCHAR(256) PRIMARY KEY,
    category_name VARCHAR(256) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Channels table
CREATE TABLE IF NOT EXISTS channels (
    channel_id VARCHAR(256) PRIMARY KEY,
    channel_name VARCHAR(256) NOT NULL,
    category_id VARCHAR(256) REFERENCES categories(category_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(256) PRIMARY KEY,
    username VARCHAR(256) NOT NULL,          -- @username (英数字のID)
    display_name VARCHAR(256),                -- 表示名（設定されている場合）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    message_id VARCHAR(256) PRIMARY KEY,
    user_id VARCHAR(256) REFERENCES users(user_id),
    content VARCHAR(256),
    channel_id VARCHAR(256) REFERENCES channels(channel_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Reactions table
CREATE TABLE IF NOT EXISTS reactions (
    reaction_id VARCHAR(256) PRIMARY KEY,
    user_id VARCHAR(256) REFERENCES users(user_id),
    message_id VARCHAR(256) REFERENCES messages(message_id),
    emoji VARCHAR(256),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_channel_id ON messages(channel_id);
CREATE INDEX IF NOT EXISTS idx_reactions_user_id ON reactions(user_id);
CREATE INDEX IF NOT EXISTS idx_reactions_message_id ON reactions(message_id);
CREATE INDEX IF NOT EXISTS idx_channels_category_id ON channels(category_id);