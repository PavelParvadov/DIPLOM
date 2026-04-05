IF OBJECT_ID('chat_messages', 'U') IS NULL
BEGIN
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
END
GO

IF NOT EXISTS (
    SELECT 1
    FROM sys.indexes
    WHERE name = 'IX_chat_messages_house_created_at'
      AND object_id = OBJECT_ID('chat_messages')
)
BEGIN
    CREATE INDEX IX_chat_messages_house_created_at ON chat_messages (house_id, created_at ASC);
END
GO
