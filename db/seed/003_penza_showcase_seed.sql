BEGIN;

DELETE FROM refresh_tokens
WHERE user_id IN (
    SELECT id
    FROM users
    WHERE login IN (
        'alexey.morozov', 'irina.lapina', 'sergey.kozlov', 'maria.belova', 'pavel.romanov',
        'olga.eliseeva', 'dmitry.tikhonov', 'ekaterina.vlasova', 'andrey.nikitin', 'svetlana.zhukova',
        'artem.sorokin', 'natalya.voronova', 'ilya.gromov', 'yuliya.makarova', 'konstantin.fedorov',
        'elena.danilova', 'mikhail.korneev', 'tatyana.savina', 'viktor.orlov', 'alina.shirokova'
    )
);

DELETE FROM houses
WHERE name IN ('Арбековский двор', 'Дом на Московской', 'Суворовский квартал', 'Пушкинский дом')
  AND created_by IN (
      SELECT id
      FROM users
      WHERE login IN ('alexey.morozov', 'dmitry.tikhonov', 'ilya.gromov', 'mikhail.korneev')
  );

DELETE FROM users
WHERE login IN (
    'alexey.morozov', 'irina.lapina', 'sergey.kozlov', 'maria.belova', 'pavel.romanov',
    'olga.eliseeva', 'dmitry.tikhonov', 'ekaterina.vlasova', 'andrey.nikitin', 'svetlana.zhukova',
    'artem.sorokin', 'natalya.voronova', 'ilya.gromov', 'yuliya.makarova', 'konstantin.fedorov',
    'elena.danilova', 'mikhail.korneev', 'tatyana.savina', 'viktor.orlov', 'alina.shirokova'
);

INSERT INTO users (login, password_hash, display_name)
VALUES
    ('alexey.morozov', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Алексей Морозов'),
    ('irina.lapina', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Ирина Лапина'),
    ('sergey.kozlov', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Сергей Козлов'),
    ('maria.belova', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Мария Белова'),
    ('pavel.romanov', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Павел Романов'),
    ('olga.eliseeva', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Ольга Елисеева'),
    ('dmitry.tikhonov', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Дмитрий Тихонов'),
    ('ekaterina.vlasova', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Екатерина Власова'),
    ('andrey.nikitin', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Андрей Никитин'),
    ('svetlana.zhukova', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Светлана Жукова'),
    ('artem.sorokin', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Артем Сорокин'),
    ('natalya.voronova', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Наталья Воронова'),
    ('ilya.gromov', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Илья Громов'),
    ('yuliya.makarova', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Юлия Макарова'),
    ('konstantin.fedorov', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Константин Федоров'),
    ('elena.danilova', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Елена Данилова'),
    ('mikhail.korneev', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Михаил Корнеев'),
    ('tatyana.savina', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Татьяна Савина'),
    ('viktor.orlov', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Виктор Орлов'),
    ('alina.shirokova', '$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW', 'Алина Широкова');

INSERT INTO houses (name, address, created_by)
SELECT data.name, data.address, u.id
FROM (
    VALUES
        ('Арбековский двор', 'г. Пенза, ул. Ладожская, 112', 'alexey.morozov'),
        ('Дом на Московской', 'г. Пенза, ул. Московская, 84', 'dmitry.tikhonov'),
        ('Суворовский квартал', 'г. Пенза, ул. Суворова, 172', 'ilya.gromov'),
        ('Пушкинский дом', 'г. Пенза, ул. Пушкина, 15', 'mikhail.korneev')
) AS data(name, address, admin_login)
JOIN users u ON u.login = data.admin_login;

INSERT INTO house_members (user_id, house_id, role)
SELECT u.id, h.id, data.role
FROM (
    VALUES
        ('Арбековский двор', 'alexey.morozov', 'admin'),
        ('Арбековский двор', 'irina.lapina', 'resident'),
        ('Арбековский двор', 'sergey.kozlov', 'resident'),
        ('Арбековский двор', 'maria.belova', 'resident'),
        ('Арбековский двор', 'pavel.romanov', 'resident'),
        ('Дом на Московской', 'dmitry.tikhonov', 'admin'),
        ('Дом на Московской', 'ekaterina.vlasova', 'resident'),
        ('Дом на Московской', 'andrey.nikitin', 'resident'),
        ('Дом на Московской', 'svetlana.zhukova', 'resident'),
        ('Дом на Московской', 'artem.sorokin', 'resident'),
        ('Суворовский квартал', 'ilya.gromov', 'admin'),
        ('Суворовский квартал', 'yuliya.makarova', 'resident'),
        ('Суворовский квартал', 'konstantin.fedorov', 'resident'),
        ('Суворовский квартал', 'elena.danilova', 'resident'),
        ('Суворовский квартал', 'tatyana.savina', 'resident'),
        ('Пушкинский дом', 'mikhail.korneev', 'admin'),
        ('Пушкинский дом', 'viktor.orlov', 'resident'),
        ('Пушкинский дом', 'alina.shirokova', 'resident'),
        ('Пушкинский дом', 'olga.eliseeva', 'resident'),
        ('Пушкинский дом', 'natalya.voronova', 'resident')
) AS data(house_name, login, role)
JOIN users u ON u.login = data.login
JOIN houses h ON h.name = data.house_name;

INSERT INTO categories (house_id, name)
SELECT h.id, category_name
FROM houses h
CROSS JOIN (
    VALUES
        ('Новости'),
        ('Объявления'),
        ('Благоустройство'),
        ('Соседи')
) AS category_list(category_name)
WHERE h.name IN ('Арбековский двор', 'Дом на Московской', 'Суворовский квартал', 'Пушкинский дом');

INSERT INTO posts (house_id, author_id, category_id, title, content, image_url)
SELECT
    h.id,
    u.id,
    c.id,
    payload.title,
    payload.content,
    NULL
FROM (
    VALUES
        ('Арбековский двор', 'alexey.morozov', 'Новости', 'Собрание жильцов в субботу', 'В эту субботу в 11:00 соберемся у первого подъезда и обсудим весенние работы во дворе и план по озеленению.'),
        ('Арбековский двор', 'irina.lapina', 'Объявления', 'Кто может забрать посылку днем', 'Буду на работе до вечера. Если курьер приедет раньше, буду благодарна соседям за помощь с получением посылки.'),
        ('Арбековский двор', 'sergey.kozlov', 'Благоустройство', 'Нужно обновить песок на детской площадке', 'После зимы покрытие стало жестким. Предлагаю включить замену песка в ближайшую заявку от дома.'),
        ('Дом на Московской', 'dmitry.tikhonov', 'Новости', 'Плановое отключение воды во вторник', 'По информации управляющей компании, во вторник с 10:00 до 14:00 отключат горячую воду для профилактических работ.'),
        ('Дом на Московской', 'ekaterina.vlasova', 'Объявления', 'Ищу соседей для совместной закупки цветов', 'Хочу заказать рассаду для клумб у подъездов. Если кто-то хочет присоединиться, напишите в комментариях.'),
        ('Дом на Московской', 'andrey.nikitin', 'Благоустройство', 'Предлагаю поставить велостойку', 'Во дворе все больше велосипедов, и их негде аккуратно парковать. Думаю, стоит вынести этот вопрос на голосование.'),
        ('Суворовский квартал', 'ilya.gromov', 'Новости', 'Открываем чат по благоустройству двора', 'Запускаем отдельное обсуждение по освещению, лавочкам и ремонту покрытия возле арки.'),
        ('Суворовский квартал', 'yuliya.makarova', 'Объявления', 'Найден детский самокат у второго подъезда', 'Самокат синего цвета стоит у консьержа. Если это ваш, можно забрать сегодня до 20:00.'),
        ('Суворовский квартал', 'konstantin.fedorov', 'Благоустройство', 'Нужно подрезать кусты у парковки', 'Кусты сильно разрослись и мешают обзору при выезде. Думаю, стоит подать коллективную заявку.'),
        ('Пушкинский дом', 'mikhail.korneev', 'Новости', 'Весенняя уборка двора в воскресенье', 'Приглашаю соседей в воскресенье к 10:30 на общую уборку территории. Перчатки и мешки возьмем централизованно.'),
        ('Пушкинский дом', 'viktor.orlov', 'Объявления', 'Отдам письменный стол в хорошем состоянии', 'Стол светлого цвета, самовывоз из третьего подъезда. Если нужен кому-то для школьника, забирайте.'),
        ('Пушкинский дом', 'alina.shirokova', 'Благоустройство', 'Хочется добавить освещение у арки', 'По вечерам у арки темно. Предлагаю собрать подписи и попросить установить дополнительный светильник.')
) AS payload(house_name, login, category_name, title, content)
JOIN users u ON u.login = payload.login
JOIN houses h ON h.name = payload.house_name
JOIN categories c ON c.house_id = h.id AND c.name = payload.category_name;

INSERT INTO comments (post_id, author_id, content)
SELECT
    p.id,
    u.id,
    payload.content
FROM (
    VALUES
        ('Собрание жильцов в субботу', 'maria.belova', 'Я буду на собрании, могу принести распечатку с предложениями по озеленению.'),
        ('Собрание жильцов в субботу', 'pavel.romanov', 'Поддерживаю, еще хотелось бы обсудить освещение у парковки.'),
        ('Нужно обновить песок на детской площадке', 'alexey.morozov', 'Хорошее замечание, включу этот вопрос в обращение в управляющую компанию.'),
        ('Нужно обновить песок на детской площадке', 'irina.lapina', 'Если понадобится подпись от жильцов, я готова подключиться.'),
        ('Плановое отключение воды во вторник', 'svetlana.zhukova', 'Спасибо за предупреждение, тогда заранее наберем воду.'),
        ('Плановое отключение воды во вторник', 'artem.sorokin', 'Хорошо, что написали. Передам информацию соседям по этажу.'),
        ('Предлагаю поставить велостойку', 'ekaterina.vlasova', 'Велостойка действительно нужна, сейчас велосипеды стоят в проходе.'),
        ('Предлагаю поставить велостойку', 'dmitry.tikhonov', 'Соберу предложения по месту установки и вынесу на голосование.'),
        ('Открываем чат по благоустройству двора', 'elena.danilova', 'Я за освещение возле арки, вечером там действительно неуютно.'),
        ('Открываем чат по благоустройству двора', 'tatyana.savina', 'Еще нужны лавочки у детской площадки, сейчас их не хватает.'),
        ('Нужно подрезать кусты у парковки', 'yuliya.makarova', 'Полностью согласна, обзор очень плохой после дождя и в сумерках.'),
        ('Нужно подрезать кусты у парковки', 'ilya.gromov', 'Сформирую заявку на этой неделе и отпишусь о сроках.'),
        ('Весенняя уборка двора в воскресенье', 'olga.eliseeva', 'Я приду и возьму с собой мешки для листвы.'),
        ('Весенняя уборка двора в воскресенье', 'natalya.voronova', 'Отлично, давно хотели привести клумбы в порядок.'),
        ('Хочется добавить освещение у арки', 'viktor.orlov', 'Поддерживаю, вечером там темно даже возле входа в арку.'),
        ('Хочется добавить освещение у арки', 'mikhail.korneev', 'Добавим это в коллективное обращение, думаю, шансы хорошие.')
) AS payload(post_title, login, content)
JOIN posts p ON p.title = payload.post_title
JOIN users u ON u.login = payload.login;

INSERT INTO invite_codes (house_id, code, created_by, expires_at)
SELECT h.id, payload.code, u.id, NOW() + INTERVAL '60 days'
FROM (
    VALUES
        ('Арбековский двор', 'alexey.morozov', 'PENZA-ARB-01'),
        ('Дом на Московской', 'dmitry.tikhonov', 'PENZA-MSK-02'),
        ('Суворовский квартал', 'ilya.gromov', 'PENZA-SUV-03'),
        ('Пушкинский дом', 'mikhail.korneev', 'PENZA-PSH-04')
) AS payload(house_name, login, code)
JOIN users u ON u.login = payload.login
JOIN houses h ON h.name = payload.house_name;

INSERT INTO chat_messages (house_id, author_id, content, image_url)
SELECT
    h.id,
    u.id,
    payload.content,
    NULL
FROM (
    VALUES
        ('Арбековский двор', 'alexey.morozov', 'Добрый вечер, соседи. Вынес повестку по субботнему собранию в ленту.'),
        ('Арбековский двор', 'maria.belova', 'Спасибо, я подготовлю список предложений по цветникам.'),
        ('Арбековский двор', 'sergey.kozlov', 'И еще прошу не забыть про покрытие на детской площадке.'),
        ('Дом на Московской', 'dmitry.tikhonov', 'Коллеги, напоминаю про отключение воды во вторник.'),
        ('Дом на Московской', 'ekaterina.vlasova', 'Хорошо, тогда перенесем полив клумб на вечер.'),
        ('Дом на Московской', 'andrey.nikitin', 'И давайте соберем идеи по велостойке в одном треде.'),
        ('Суворовский квартал', 'ilya.gromov', 'Запускаем обсуждение благоустройства прямо здесь, чтобы ничего не потерялось.'),
        ('Суворовский квартал', 'elena.danilova', 'За освещение у арки и за новые лавочки возле площадки.'),
        ('Суворовский квартал', 'konstantin.fedorov', 'И еще нужна подрезка кустов у парковки, там плохой обзор.'),
        ('Пушкинский дом', 'mikhail.korneev', 'В воскресенье встречаемся на уборку двора в 10:30.'),
        ('Пушкинский дом', 'olga.eliseeva', 'Я принесу перчатки и пару лишних мешков.'),
        ('Пушкинский дом', 'alina.shirokova', 'После уборки предлагаю обсудить освещение у арки.')
) AS payload(house_name, login, content)
JOIN users u ON u.login = payload.login
JOIN houses h ON h.name = payload.house_name;

COMMIT;
