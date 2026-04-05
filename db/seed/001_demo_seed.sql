BEGIN;

DELETE FROM refresh_tokens
WHERE user_id IN (
    SELECT id
    FROM users
    WHERE login IN ('admin_demo', 'resident_demo')
);

DELETE FROM houses
WHERE name = 'HappyHouse Tower'
  AND created_by IN (
      SELECT id
      FROM users
      WHERE login = 'admin_demo'
  );

DELETE FROM users
WHERE login IN ('admin_demo', 'resident_demo');

INSERT INTO users (login, password_hash, display_name)
VALUES
    ('admin_demo', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Админ Demo'),
    ('resident_demo', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Житель Demo');

INSERT INTO houses (name, address, created_by)
SELECT 'HappyHouse Tower', 'Москва, ул. Соседская, 12', id
FROM users
WHERE login = 'admin_demo';

INSERT INTO house_members (user_id, house_id, role)
SELECT u.id, h.id, role_map.role
FROM users u
JOIN houses h ON h.name = 'HappyHouse Tower'
JOIN (
    VALUES
        ('admin_demo', 'admin'),
        ('resident_demo', 'resident')
) AS role_map(login, role) ON role_map.login = u.login;

INSERT INTO categories (house_id, name)
SELECT h.id, category_name
FROM houses h
CROSS JOIN (
    VALUES
        ('Новости'),
        ('Объявления'),
        ('Благоустройство')
) AS categories(category_name)
WHERE h.name = 'HappyHouse Tower';

INSERT INTO posts (house_id, author_id, category_id, title, content)
SELECT
    h.id,
    u.id,
    c.id,
    payload.title,
    payload.content
FROM (
    VALUES
        (
            'admin_demo',
            'Новости',
            'Добро пожаловать в HappyHouse',
            'Это демонстрационный дом для презентации проекта. Здесь можно публиковать новости, объявления и обсуждать вопросы дома.'
        ),
        (
            'resident_demo',
            'Объявления',
            'Собираем предложения по двору',
            'Если у вас есть идеи по озеленению или освещению двора, напишите в комментариях к этому посту.'
        )
) AS payload(login, category_name, title, content)
JOIN users u ON u.login = payload.login
JOIN houses h ON h.name = 'HappyHouse Tower'
JOIN categories c ON c.house_id = h.id AND c.name = payload.category_name;

INSERT INTO comments (post_id, author_id, content)
SELECT
    p.id,
    u.id,
    payload.content
FROM (
    VALUES
        (
            'Собираем предложения по двору',
            'admin_demo',
            'Отличная тема. На защите можно показать этот сценарий как пример реального общения соседей.'
        ),
        (
            'Собираем предложения по двору',
            'resident_demo',
            'Я бы предложил добавить освещение у второго подъезда и больше лавочек.'
        )
) AS payload(post_title, login, content)
JOIN posts p ON p.title = payload.post_title
JOIN users u ON u.login = payload.login;

INSERT INTO invite_codes (house_id, code, created_by, expires_at)
SELECT h.id, 'HAPPY2026', u.id, NOW() + INTERVAL '30 days'
FROM houses h
JOIN users u ON u.login = 'admin_demo'
WHERE h.name = 'HappyHouse Tower';

COMMIT;
