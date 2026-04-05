IF OBJECT_ID('refresh_tokens', 'U') IS NOT NULL DROP TABLE refresh_tokens;
IF OBJECT_ID('chat_messages', 'U') IS NOT NULL DROP TABLE chat_messages;
IF OBJECT_ID('comments', 'U') IS NOT NULL DROP TABLE comments;
IF OBJECT_ID('posts', 'U') IS NOT NULL DROP TABLE posts;
IF OBJECT_ID('invite_codes', 'U') IS NOT NULL DROP TABLE invite_codes;
IF OBJECT_ID('categories', 'U') IS NOT NULL DROP TABLE categories;
IF OBJECT_ID('house_members', 'U') IS NOT NULL DROP TABLE house_members;
IF OBJECT_ID('houses', 'U') IS NOT NULL DROP TABLE houses;
IF OBJECT_ID('users', 'U') IS NOT NULL DROP TABLE users;
GO

CREATE TABLE users (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    login NVARCHAR(64) NOT NULL UNIQUE,
    password_hash NVARCHAR(255) NOT NULL,
    display_name NVARCHAR(120) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME()
);
GO

CREATE TABLE houses (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(120) NOT NULL,
    address NVARCHAR(255) NOT NULL,
    created_by BIGINT NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT FK_houses_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);
GO

CREATE TABLE house_members (
    user_id BIGINT NOT NULL,
    house_id BIGINT NOT NULL,
    role NVARCHAR(20) NOT NULL CHECK (role IN (N'resident', N'admin')),
    joined_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT PK_house_members PRIMARY KEY (user_id, house_id),
    CONSTRAINT FK_house_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT FK_house_members_house FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE
);
GO

CREATE TABLE categories (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    house_id BIGINT NOT NULL,
    name NVARCHAR(80) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT FK_categories_house FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE,
    CONSTRAINT UQ_categories_house_name UNIQUE (house_id, name)
);
GO

CREATE TABLE posts (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    house_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    title NVARCHAR(180) NOT NULL,
    content NVARCHAR(MAX) NOT NULL,
    image_url NVARCHAR(255) NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    updated_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT FK_posts_house FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE,
    CONSTRAINT FK_posts_author FOREIGN KEY (author_id) REFERENCES users(id),
    CONSTRAINT FK_posts_category FOREIGN KEY (category_id) REFERENCES categories(id)
);
GO

CREATE TABLE comments (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    post_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    content NVARCHAR(MAX) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT FK_comments_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT FK_comments_author FOREIGN KEY (author_id) REFERENCES users(id)
);
GO

CREATE TABLE chat_messages (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    house_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    content NVARCHAR(MAX) NULL,
    image_url NVARCHAR(255) NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT FK_chat_messages_house FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE,
    CONSTRAINT FK_chat_messages_author FOREIGN KEY (author_id) REFERENCES users(id)
);
GO

CREATE TABLE invite_codes (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    house_id BIGINT NOT NULL,
    code NVARCHAR(32) NOT NULL UNIQUE,
    created_by BIGINT NOT NULL,
    is_active BIT NOT NULL DEFAULT 1,
    expires_at DATETIME2 NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT FK_invite_codes_house FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE,
    CONSTRAINT FK_invite_codes_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);
GO

CREATE TABLE refresh_tokens (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token NVARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME2 NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT SYSDATETIME(),
    CONSTRAINT FK_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
GO

CREATE INDEX IX_posts_house_category_created_at ON posts (house_id, category_id, created_at DESC);
CREATE INDEX IX_comments_post_created_at ON comments (post_id, created_at ASC);
CREATE INDEX IX_chat_messages_house_created_at ON chat_messages (house_id, created_at ASC);
CREATE INDEX IX_house_members_house_id ON house_members (house_id);
CREATE INDEX IX_invite_codes_house_id ON invite_codes (house_id, is_active);
GO
