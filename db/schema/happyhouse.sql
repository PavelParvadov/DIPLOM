CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(120) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE houses (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    address VARCHAR(255) NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE house_members (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    house_id BIGINT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('resident', 'admin')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, house_id)
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    house_id BIGINT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    name VARCHAR(80) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (house_id, name)
);

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    house_id BIGINT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES users(id),
    category_id BIGINT NOT NULL REFERENCES categories(id),
    title VARCHAR(180) NOT NULL,
    content TEXT NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_messages (
    id BIGSERIAL PRIMARY KEY,
    house_id BIGINT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT,
    image_url VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE media_assets (
    id BIGSERIAL PRIMARY KEY,
    public_id VARCHAR(255) NOT NULL UNIQUE,
    content_type VARCHAR(120) NOT NULL,
    data BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE invite_codes (
    id BIGSERIAL PRIMARY KEY,
    house_id BIGINT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    code VARCHAR(32) NOT NULL UNIQUE,
    created_by BIGINT NOT NULL REFERENCES users(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_house_members_house_id ON house_members(house_id);
CREATE INDEX idx_categories_house_id ON categories(house_id);
CREATE INDEX idx_posts_house_category_created_at ON posts(house_id, category_id, created_at DESC);
CREATE INDEX idx_posts_house_created_at ON posts(house_id, created_at DESC);
CREATE INDEX idx_comments_post_created_at ON comments(post_id, created_at ASC);
CREATE INDEX idx_chat_messages_house_created_at ON chat_messages(house_id, created_at DESC);
CREATE INDEX idx_media_assets_public_id ON media_assets(public_id);
CREATE INDEX idx_invite_codes_house_created_at ON invite_codes(house_id, created_at DESC);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
