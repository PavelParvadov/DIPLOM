DECLARE @adminPasswordHash NVARCHAR(255) = N'$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW';
DECLARE @residentPasswordHash NVARCHAR(255) = N'$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW';

INSERT INTO users (login, password_hash, display_name)
VALUES
    (N'admin_demo', @adminPasswordHash, N'Админ Demo'),
    (N'resident_demo', @residentPasswordHash, N'Житель Demo');

DECLARE @adminId BIGINT = (SELECT id FROM users WHERE login = N'admin_demo');
DECLARE @residentId BIGINT = (SELECT id FROM users WHERE login = N'resident_demo');

INSERT INTO houses (name, address, created_by)
VALUES (N'HappyHouse Tower', N'Москва, ул. Соседская, 12', @adminId);

DECLARE @houseId BIGINT = SCOPE_IDENTITY();

INSERT INTO house_members (user_id, house_id, role)
VALUES
    (@adminId, @houseId, N'admin'),
    (@residentId, @houseId, N'resident');

INSERT INTO categories (house_id, name)
VALUES
    (@houseId, N'Новости'),
    (@houseId, N'Объявления'),
    (@houseId, N'Благоустройство');

DECLARE @newsCategoryId BIGINT = (SELECT id FROM categories WHERE house_id = @houseId AND name = N'Новости');
DECLARE @adsCategoryId BIGINT = (SELECT id FROM categories WHERE house_id = @houseId AND name = N'Объявления');

INSERT INTO posts (house_id, author_id, category_id, title, content)
VALUES
    (
        @houseId,
        @adminId,
        @newsCategoryId,
        N'Добро пожаловать в HappyHouse',
        N'Это демонстрационный дом для презентации проекта. Здесь можно публиковать новости, объявления и обсуждать вопросы дома.'
    ),
    (
        @houseId,
        @residentId,
        @adsCategoryId,
        N'Собираем предложения по двору',
        N'Если у вас есть идеи по озеленению или освещению двора, напишите в комментариях к этому посту.'
    );

DECLARE @discussionPostId BIGINT = (
    SELECT TOP 1 id
    FROM posts
    WHERE house_id = @houseId AND title = N'Собираем предложения по двору'
    ORDER BY id DESC
);

INSERT INTO comments (post_id, author_id, content)
VALUES
    (@discussionPostId, @adminId, N'Отличная тема. На защите можно показать этот сценарий как пример реального общения соседей.'),
    (@discussionPostId, @residentId, N'Я бы предложил добавить освещение у второго подъезда и больше лавочек.');

INSERT INTO invite_codes (house_id, code, created_by, expires_at)
VALUES (@houseId, N'HAPPY2026', @adminId, DATEADD(DAY, 30, SYSDATETIME()));
