-- Qube: Initial Schema
-- PostgreSQL 16+

-- ============================================================
-- Extensions
-- ============================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- Users
-- ============================================================
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username        VARCHAR(30) NOT NULL UNIQUE,
    display_name    VARCHAR(50) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    bio             TEXT DEFAULT '',
    avatar_url      TEXT DEFAULT '',
    header_url      TEXT DEFAULT '',
    location        VARCHAR(100) DEFAULT '',
    website         VARCHAR(255) DEFAULT '',
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    is_private      BOOLEAN NOT NULL DEFAULT FALSE,
    follower_count  INTEGER NOT NULL DEFAULT 0,
    following_count INTEGER NOT NULL DEFAULT 0,
    post_count      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_created_at ON users (created_at);

-- ============================================================
-- Posts
-- ============================================================
CREATE TABLE posts (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content         TEXT NOT NULL,
    reply_to_id     UUID REFERENCES posts(id) ON DELETE SET NULL,
    repost_of_id    UUID REFERENCES posts(id) ON DELETE SET NULL,
    quote_of_id     UUID REFERENCES posts(id) ON DELETE SET NULL,
    like_count      INTEGER NOT NULL DEFAULT 0,
    repost_count    INTEGER NOT NULL DEFAULT 0,
    reply_count     INTEGER NOT NULL DEFAULT 0,
    quote_count     INTEGER NOT NULL DEFAULT 0,
    is_deleted      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_posts_created_at ON posts (created_at DESC);
CREATE INDEX idx_posts_reply_to ON posts (reply_to_id) WHERE reply_to_id IS NOT NULL;
CREATE INDEX idx_posts_user_timeline ON posts (user_id, created_at DESC) WHERE is_deleted = FALSE;

-- ============================================================
-- Post Media (images, videos)
-- ============================================================
CREATE TABLE post_media (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id     UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    media_type  VARCHAR(20) NOT NULL CHECK (media_type IN ('image', 'video', 'gif')),
    url         TEXT NOT NULL,
    thumbnail_url TEXT DEFAULT '',
    width       INTEGER,
    height      INTEGER,
    sort_order  SMALLINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_post_media_post_id ON post_media (post_id);

-- ============================================================
-- Follows
-- ============================================================
CREATE TABLE follows (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX idx_follows_follower ON follows (follower_id, created_at DESC);
CREATE INDEX idx_follows_following ON follows (following_id, created_at DESC);

-- ============================================================
-- Likes
-- ============================================================
CREATE TABLE likes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id     UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, post_id)
);

CREATE INDEX idx_likes_post_id ON likes (post_id);
CREATE INDEX idx_likes_user_id ON likes (user_id, created_at DESC);

-- ============================================================
-- Reposts
-- ============================================================
CREATE TABLE reposts (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id     UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, post_id)
);

CREATE INDEX idx_reposts_post_id ON reposts (post_id);
CREATE INDEX idx_reposts_user_id ON reposts (user_id, created_at DESC);

-- ============================================================
-- Bookmarks
-- ============================================================
CREATE TABLE bookmarks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id     UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, post_id)
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks (user_id, created_at DESC);

-- ============================================================
-- Conversations (DM)
-- ============================================================
CREATE TABLE conversations (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    is_group    BOOLEAN NOT NULL DEFAULT FALSE,
    name        VARCHAR(100) DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Conversation Participants
-- ============================================================
CREATE TABLE conversation_participants (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (conversation_id, user_id)
);

CREATE INDEX idx_conv_participants_user ON conversation_participants (user_id);
CREATE INDEX idx_conv_participants_conv ON conversation_participants (conversation_id);

-- ============================================================
-- Messages
-- ============================================================
CREATE TABLE messages (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content         TEXT NOT NULL,
    is_deleted      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_conv ON messages (conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON messages (sender_id);

-- ============================================================
-- Notifications
-- ============================================================
CREATE TABLE notifications (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        VARCHAR(20) NOT NULL CHECK (type IN ('like', 'repost', 'follow', 'reply', 'quote', 'mention', 'dm')),
    post_id     UUID REFERENCES posts(id) ON DELETE CASCADE,
    is_read     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications (user_id, created_at DESC);
CREATE INDEX idx_notifications_unread ON notifications (user_id, is_read, created_at DESC) WHERE is_read = FALSE;

-- ============================================================
-- Hashtags
-- ============================================================
CREATE TABLE hashtags (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    post_count  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_hashtags_name ON hashtags (name);

-- ============================================================
-- Post Hashtags (many-to-many)
-- ============================================================
CREATE TABLE post_hashtags (
    post_id     UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    hashtag_id  UUID NOT NULL REFERENCES hashtags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, hashtag_id)
);

CREATE INDEX idx_post_hashtags_hashtag ON post_hashtags (hashtag_id);

-- ============================================================
-- Timeline Cursor (未読管理: Qubeの差別化ポイント)
-- ============================================================
CREATE TABLE timeline_cursors (
    user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    last_seen_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_post_id UUID REFERENCES posts(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Blocks / Mutes
-- ============================================================
CREATE TABLE blocks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    blocker_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
);

CREATE TABLE mutes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    muter_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    muted_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (muter_id, muted_id),
    CHECK (muter_id != muted_id)
);

CREATE INDEX idx_blocks_blocker ON blocks (blocker_id);
CREATE INDEX idx_mutes_muter ON mutes (muter_id);

-- ============================================================
-- Sessions (JWT refresh token management)
-- ============================================================
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token   TEXT NOT NULL UNIQUE,
    user_agent      TEXT DEFAULT '',
    ip_address      INET,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON sessions (user_id);
CREATE INDEX idx_sessions_token ON sessions (refresh_token);

-- ============================================================
-- Updated_at trigger function
-- ============================================================
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_posts_updated_at BEFORE UPDATE ON posts FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_conversations_updated_at BEFORE UPDATE ON conversations FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_timeline_cursors_updated_at BEFORE UPDATE ON timeline_cursors FOR EACH ROW EXECUTE FUNCTION update_updated_at();
