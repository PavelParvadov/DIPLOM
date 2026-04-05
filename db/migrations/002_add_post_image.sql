IF COL_LENGTH('posts', 'image_url') IS NULL
BEGIN
    ALTER TABLE posts
    ADD image_url NVARCHAR(255) NULL;
END
GO
