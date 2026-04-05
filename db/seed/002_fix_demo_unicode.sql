SET NOCOUNT ON;

UPDATE users
SET display_name = CASE login
    WHEN N'admin_demo' THEN N'Админ Demo'
    WHEN N'resident_demo' THEN N'Житель Demo'
    ELSE display_name
END
WHERE login IN (N'admin_demo', N'resident_demo');

UPDATE houses
SET address = N'Москва, ул. Соседская, 12'
WHERE name = N'HappyHouse Tower';

WITH ordered_categories AS (
    SELECT id, ROW_NUMBER() OVER (PARTITION BY house_id ORDER BY id) AS rn
    FROM categories
    WHERE house_id IN (SELECT id FROM houses WHERE name = N'HappyHouse Tower')
)
UPDATE categories
SET name = CASE ordered_categories.rn
    WHEN 1 THEN N'Новости'
    WHEN 2 THEN N'Объявления'
    ELSE N'Благоустройство'
END
FROM categories
JOIN ordered_categories ON ordered_categories.id = categories.id;

UPDATE posts
SET title = CASE
        WHEN author_id = (SELECT id FROM users WHERE login = N'admin_demo') THEN N'Добро пожаловать в HappyHouse'
        ELSE N'Собираем предложения по двору'
    END,
    content = CASE
        WHEN author_id = (SELECT id FROM users WHERE login = N'admin_demo') THEN N'Это демонстрационный дом для презентации проекта. Здесь можно публиковать новости, объявления и обсуждать вопросы дома.'
        ELSE N'Если у вас есть идеи по озеленению или освещению двора, напишите в комментариях к этому посту.'
    END
WHERE house_id IN (SELECT id FROM houses WHERE name = N'HappyHouse Tower');

WITH ordered_comments AS (
    SELECT comments.id, ROW_NUMBER() OVER (ORDER BY comments.id) AS rn
    FROM comments
    JOIN posts ON posts.id = comments.post_id
    WHERE posts.house_id IN (SELECT id FROM houses WHERE name = N'HappyHouse Tower')
)
UPDATE comments
SET content = CASE ordered_comments.rn
    WHEN 1 THEN N'Отличная тема. На защите можно показать этот сценарий как пример реального общения соседей.'
    WHEN 2 THEN N'Я бы предложил добавить освещение у второго подъезда и больше лавочек.'
    ELSE comments.content
END
FROM comments
JOIN ordered_comments ON ordered_comments.id = comments.id;
