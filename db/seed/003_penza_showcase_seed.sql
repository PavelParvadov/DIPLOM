BEGIN TRANSACTION;

DECLARE @passwordHash NVARCHAR(255) = N'$2a$10$9biuxUCTZGQVN0nLh9LWjuuBs8QrFIvZjgtbA65bvcu64kJHkKzgW';

DECLARE @showcaseUsers TABLE (login NVARCHAR(64) PRIMARY KEY);
INSERT INTO @showcaseUsers (login)
VALUES
    (N'alexey.morozov'),
    (N'irina.lapina'),
    (N'sergey.kozlov'),
    (N'maria.belova'),
    (N'pavel.romanov'),
    (N'olga.eliseeva'),
    (N'dmitry.tikhonov'),
    (N'ekaterina.vlasova'),
    (N'andrey.nikitin'),
    (N'svetlana.zhukova'),
    (N'artem.sorokin'),
    (N'natalya.voronova'),
    (N'ilya.gromov'),
    (N'yuliya.makarova'),
    (N'konstantin.fedorov'),
    (N'elena.danilova'),
    (N'mikhail.korneev'),
    (N'tatyana.savina'),
    (N'viktor.orlov'),
    (N'alina.shirokova');

DECLARE @showcaseUserIds TABLE (id BIGINT PRIMARY KEY);
INSERT INTO @showcaseUserIds (id)
SELECT id
FROM users
WHERE login IN (SELECT login FROM @showcaseUsers);

DECLARE @showcaseHouseIds TABLE (id BIGINT PRIMARY KEY);
INSERT INTO @showcaseHouseIds (id)
SELECT id
FROM houses
WHERE created_by IN (SELECT id FROM @showcaseUserIds);

DELETE FROM invite_codes
WHERE house_id IN (SELECT id FROM @showcaseHouseIds)
   OR code LIKE N'PENZA-%';

DELETE FROM chat_messages
WHERE house_id IN (SELECT id FROM @showcaseHouseIds)
   OR author_id IN (SELECT id FROM @showcaseUserIds);

DELETE FROM comments
WHERE post_id IN (
    SELECT id
    FROM posts
    WHERE house_id IN (SELECT id FROM @showcaseHouseIds)
)
   OR author_id IN (SELECT id FROM @showcaseUserIds);

DELETE FROM posts
WHERE house_id IN (SELECT id FROM @showcaseHouseIds)
   OR author_id IN (SELECT id FROM @showcaseUserIds);

DELETE FROM categories
WHERE house_id IN (SELECT id FROM @showcaseHouseIds);

DELETE FROM house_members
WHERE house_id IN (SELECT id FROM @showcaseHouseIds)
   OR user_id IN (SELECT id FROM @showcaseUserIds);

DELETE FROM houses
WHERE id IN (SELECT id FROM @showcaseHouseIds);

DELETE FROM refresh_tokens
WHERE user_id IN (SELECT id FROM @showcaseUserIds);

DELETE FROM users
WHERE login IN (SELECT login FROM @showcaseUsers);

INSERT INTO users (login, password_hash, display_name)
VALUES
    (N'alexey.morozov', @passwordHash, N'Алексей Морозов'),
    (N'irina.lapina', @passwordHash, N'Ирина Лапина'),
    (N'sergey.kozlov', @passwordHash, N'Сергей Козлов'),
    (N'maria.belova', @passwordHash, N'Мария Белова'),
    (N'pavel.romanov', @passwordHash, N'Павел Романов'),
    (N'olga.eliseeva', @passwordHash, N'Ольга Елисеева'),
    (N'dmitry.tikhonov', @passwordHash, N'Дмитрий Тихонов'),
    (N'ekaterina.vlasova', @passwordHash, N'Екатерина Власова'),
    (N'andrey.nikitin', @passwordHash, N'Андрей Никитин'),
    (N'svetlana.zhukova', @passwordHash, N'Светлана Жукова'),
    (N'artem.sorokin', @passwordHash, N'Артем Сорокин'),
    (N'natalya.voronova', @passwordHash, N'Наталья Воронова'),
    (N'ilya.gromov', @passwordHash, N'Илья Громов'),
    (N'yuliya.makarova', @passwordHash, N'Юлия Макарова'),
    (N'konstantin.fedorov', @passwordHash, N'Константин Федоров'),
    (N'elena.danilova', @passwordHash, N'Елена Данилова'),
    (N'mikhail.korneev', @passwordHash, N'Михаил Корнеев'),
    (N'tatyana.savina', @passwordHash, N'Татьяна Савина'),
    (N'viktor.orlov', @passwordHash, N'Виктор Орлов'),
    (N'alina.shirokova', @passwordHash, N'Алина Широкова');

DECLARE @alexeyId BIGINT = (SELECT id FROM users WHERE login = N'alexey.morozov');
DECLARE @irinaId BIGINT = (SELECT id FROM users WHERE login = N'irina.lapina');
DECLARE @sergeyId BIGINT = (SELECT id FROM users WHERE login = N'sergey.kozlov');
DECLARE @mariaId BIGINT = (SELECT id FROM users WHERE login = N'maria.belova');
DECLARE @pavelId BIGINT = (SELECT id FROM users WHERE login = N'pavel.romanov');
DECLARE @olgaId BIGINT = (SELECT id FROM users WHERE login = N'olga.eliseeva');
DECLARE @dmitryId BIGINT = (SELECT id FROM users WHERE login = N'dmitry.tikhonov');
DECLARE @ekaterinaId BIGINT = (SELECT id FROM users WHERE login = N'ekaterina.vlasova');
DECLARE @andreyId BIGINT = (SELECT id FROM users WHERE login = N'andrey.nikitin');
DECLARE @svetlanaId BIGINT = (SELECT id FROM users WHERE login = N'svetlana.zhukova');
DECLARE @artemId BIGINT = (SELECT id FROM users WHERE login = N'artem.sorokin');
DECLARE @natalyaId BIGINT = (SELECT id FROM users WHERE login = N'natalya.voronova');
DECLARE @ilyaId BIGINT = (SELECT id FROM users WHERE login = N'ilya.gromov');
DECLARE @yuliyaId BIGINT = (SELECT id FROM users WHERE login = N'yuliya.makarova');
DECLARE @konstantinId BIGINT = (SELECT id FROM users WHERE login = N'konstantin.fedorov');
DECLARE @elenaId BIGINT = (SELECT id FROM users WHERE login = N'elena.danilova');
DECLARE @mikhailId BIGINT = (SELECT id FROM users WHERE login = N'mikhail.korneev');
DECLARE @tatyanaId BIGINT = (SELECT id FROM users WHERE login = N'tatyana.savina');
DECLARE @viktorId BIGINT = (SELECT id FROM users WHERE login = N'viktor.orlov');
DECLARE @alinaId BIGINT = (SELECT id FROM users WHERE login = N'alina.shirokova');

INSERT INTO houses (name, address, created_by)
VALUES
    (N'Арбековский двор', N'г. Пенза, ул. Ладожская, 112', @alexeyId),
    (N'Дом на Московской', N'г. Пенза, ул. Московская, 84', @dmitryId),
    (N'Суворовский квартал', N'г. Пенза, ул. Суворова, 172', @ilyaId),
    (N'Пушкинский дом', N'г. Пенза, ул. Пушкина, 15', @mikhailId);

DECLARE @arbekovoHouseId BIGINT = (SELECT id FROM houses WHERE name = N'Арбековский двор');
DECLARE @moskovskayaHouseId BIGINT = (SELECT id FROM houses WHERE name = N'Дом на Московской');
DECLARE @suvorovaHouseId BIGINT = (SELECT id FROM houses WHERE name = N'Суворовский квартал');
DECLARE @pushkinaHouseId BIGINT = (SELECT id FROM houses WHERE name = N'Пушкинский дом');

INSERT INTO house_members (user_id, house_id, role)
VALUES
    (@alexeyId, @arbekovoHouseId, N'admin'),
    (@irinaId, @arbekovoHouseId, N'resident'),
    (@sergeyId, @arbekovoHouseId, N'resident'),
    (@mariaId, @arbekovoHouseId, N'resident'),
    (@pavelId, @arbekovoHouseId, N'resident'),

    (@dmitryId, @moskovskayaHouseId, N'admin'),
    (@ekaterinaId, @moskovskayaHouseId, N'resident'),
    (@andreyId, @moskovskayaHouseId, N'resident'),
    (@svetlanaId, @moskovskayaHouseId, N'resident'),
    (@artemId, @moskovskayaHouseId, N'resident'),

    (@ilyaId, @suvorovaHouseId, N'admin'),
    (@yuliyaId, @suvorovaHouseId, N'resident'),
    (@konstantinId, @suvorovaHouseId, N'resident'),
    (@elenaId, @suvorovaHouseId, N'resident'),
    (@tatyanaId, @suvorovaHouseId, N'resident'),

    (@mikhailId, @pushkinaHouseId, N'admin'),
    (@viktorId, @pushkinaHouseId, N'resident'),
    (@alinaId, @pushkinaHouseId, N'resident'),
    (@olgaId, @pushkinaHouseId, N'resident'),
    (@natalyaId, @pushkinaHouseId, N'resident');

INSERT INTO categories (house_id, name)
VALUES
    (@arbekovoHouseId, N'Новости'),
    (@arbekovoHouseId, N'Объявления'),
    (@arbekovoHouseId, N'Благоустройство'),
    (@arbekovoHouseId, N'Соседи'),

    (@moskovskayaHouseId, N'Новости'),
    (@moskovskayaHouseId, N'Объявления'),
    (@moskovskayaHouseId, N'Благоустройство'),
    (@moskovskayaHouseId, N'Соседи'),

    (@suvorovaHouseId, N'Новости'),
    (@suvorovaHouseId, N'Объявления'),
    (@suvorovaHouseId, N'Благоустройство'),
    (@suvorovaHouseId, N'Соседи'),

    (@pushkinaHouseId, N'Новости'),
    (@pushkinaHouseId, N'Объявления'),
    (@pushkinaHouseId, N'Благоустройство'),
    (@pushkinaHouseId, N'Соседи');

DECLARE @arbekovoNewsId BIGINT = (SELECT id FROM categories WHERE house_id = @arbekovoHouseId AND name = N'Новости');
DECLARE @arbekovoAdsId BIGINT = (SELECT id FROM categories WHERE house_id = @arbekovoHouseId AND name = N'Объявления');
DECLARE @arbekovoImproveId BIGINT = (SELECT id FROM categories WHERE house_id = @arbekovoHouseId AND name = N'Благоустройство');

DECLARE @moskovskayaNewsId BIGINT = (SELECT id FROM categories WHERE house_id = @moskovskayaHouseId AND name = N'Новости');
DECLARE @moskovskayaAdsId BIGINT = (SELECT id FROM categories WHERE house_id = @moskovskayaHouseId AND name = N'Объявления');
DECLARE @moskovskayaImproveId BIGINT = (SELECT id FROM categories WHERE house_id = @moskovskayaHouseId AND name = N'Благоустройство');

DECLARE @suvorovaNewsId BIGINT = (SELECT id FROM categories WHERE house_id = @suvorovaHouseId AND name = N'Новости');
DECLARE @suvorovaAdsId BIGINT = (SELECT id FROM categories WHERE house_id = @suvorovaHouseId AND name = N'Объявления');
DECLARE @suvorovaImproveId BIGINT = (SELECT id FROM categories WHERE house_id = @suvorovaHouseId AND name = N'Благоустройство');

DECLARE @pushkinaNewsId BIGINT = (SELECT id FROM categories WHERE house_id = @pushkinaHouseId AND name = N'Новости');
DECLARE @pushkinaAdsId BIGINT = (SELECT id FROM categories WHERE house_id = @pushkinaHouseId AND name = N'Объявления');
DECLARE @pushkinaImproveId BIGINT = (SELECT id FROM categories WHERE house_id = @pushkinaHouseId AND name = N'Благоустройство');

INSERT INTO posts (house_id, author_id, category_id, title, content, image_url)
VALUES
    (@arbekovoHouseId, @alexeyId, @arbekovoNewsId, N'Собрание жильцов в субботу', N'В эту субботу в 11:00 соберемся у первого подъезда и обсудим весенние работы во дворе и план по озеленению.', NULL),
    (@arbekovoHouseId, @irinaId, @arbekovoAdsId, N'Кто может забрать посылку днем', N'Буду на работе до вечера. Если курьер приедет раньше, буду благодарна соседям за помощь с получением небольшой посылки.', NULL),
    (@arbekovoHouseId, @sergeyId, @arbekovoImproveId, N'Нужно обновить песок на детской площадке', N'После зимы покрытие стало жестким. Предлагаю включить замену песка в ближайшую заявку от дома.', NULL),

    (@moskovskayaHouseId, @dmitryId, @moskovskayaNewsId, N'Плановое отключение воды во вторник', N'По информации управляющей компании, во вторник с 10:00 до 14:00 отключат горячую воду для профилактических работ.', NULL),
    (@moskovskayaHouseId, @ekaterinaId, @moskovskayaAdsId, N'Ищу соседей для совместной закупки цветов', N'Хочу заказать рассаду для клумб у подъездов. Если кто-то хочет присоединиться, напишите в комментариях.', NULL),
    (@moskovskayaHouseId, @andreyId, @moskovskayaImproveId, N'Предлагаю поставить велостойку', N'Во дворе все больше велосипедов, и их негде аккуратно парковать. Думаю, стоит вынести этот вопрос на голосование.', NULL),

    (@suvorovaHouseId, @ilyaId, @suvorovaNewsId, N'Открываем чат по благоустройству двора', N'Запускаем отдельное обсуждение по освещению, лавочкам и ремонту покрытия возле арки. Пишите предложения, чтобы ничего не потерялось.', NULL),
    (@suvorovaHouseId, @yuliyaId, @suvorovaAdsId, N'Найден детский самокат у второго подъезда', N'Самокат синего цвета стоит у консьержа. Если это ваш, можно забрать сегодня до 20:00.', NULL),
    (@suvorovaHouseId, @konstantinId, @suvorovaImproveId, N'Нужно подрезать кусты у парковки', N'Кусты сильно разрослись и мешают обзору при выезде. Думаю, стоит подать коллективную заявку.', NULL),

    (@pushkinaHouseId, @mikhailId, @pushkinaNewsId, N'Весенняя уборка двора в воскресенье', N'Приглашаю соседей в воскресенье к 10:30 на общую уборку территории. Перчатки и мешки возьмем централизованно.', NULL),
    (@pushkinaHouseId, @viktorId, @pushkinaAdsId, N'Отдам письменный стол в хорошем состоянии', N'Стол светлого цвета, самовывоз из третьего подъезда. Если нужен кому-то для школьника, забирайте.', NULL),
    (@pushkinaHouseId, @alinaId, @pushkinaImproveId, N'Хочется добавить освещение у арки', N'По вечерам у арки темно. Предлагаю собрать подписи и попросить установить дополнительный светильник.', NULL);

DECLARE @arbekovoMeetingPostId BIGINT = (SELECT id FROM posts WHERE house_id = @arbekovoHouseId AND title = N'Собрание жильцов в субботу');
DECLARE @arbekovoSandPostId BIGINT = (SELECT id FROM posts WHERE house_id = @arbekovoHouseId AND title = N'Нужно обновить песок на детской площадке');
DECLARE @moskovskayaWaterPostId BIGINT = (SELECT id FROM posts WHERE house_id = @moskovskayaHouseId AND title = N'Плановое отключение воды во вторник');
DECLARE @moskovskayaBikePostId BIGINT = (SELECT id FROM posts WHERE house_id = @moskovskayaHouseId AND title = N'Предлагаю поставить велостойку');
DECLARE @suvorovaChatPostId BIGINT = (SELECT id FROM posts WHERE house_id = @suvorovaHouseId AND title = N'Открываем чат по благоустройству двора');
DECLARE @suvorovaBushesPostId BIGINT = (SELECT id FROM posts WHERE house_id = @suvorovaHouseId AND title = N'Нужно подрезать кусты у парковки');
DECLARE @pushkinaCleanupPostId BIGINT = (SELECT id FROM posts WHERE house_id = @pushkinaHouseId AND title = N'Весенняя уборка двора в воскресенье');
DECLARE @pushkinaLightPostId BIGINT = (SELECT id FROM posts WHERE house_id = @pushkinaHouseId AND title = N'Хочется добавить освещение у арки');

INSERT INTO comments (post_id, author_id, content)
VALUES
    (@arbekovoMeetingPostId, @mariaId, N'Я буду на собрании, могу принести распечатку с предложениями по озеленению.'),
    (@arbekovoMeetingPostId, @pavelId, N'Поддерживаю, еще хотелось бы обсудить освещение у парковки.'),
    (@arbekovoSandPostId, @alexeyId, N'Хорошее замечание, включу этот вопрос в обращение в управляющую компанию.'),
    (@arbekovoSandPostId, @irinaId, N'Если понадобится подпись от жильцов, я готова подключиться.'),

    (@moskovskayaWaterPostId, @svetlanaId, N'Спасибо за предупреждение, тогда заранее наберем воду.'),
    (@moskovskayaWaterPostId, @artemId, N'Хорошо, что написали. Передам информацию соседям по этажу.'),
    (@moskovskayaBikePostId, @ekaterinaId, N'Велостойка действительно нужна, сейчас велосипеды стоят в проходе.'),
    (@moskovskayaBikePostId, @dmitryId, N'Соберу предложения по месту установки и вынесу на голосование.'),

    (@suvorovaChatPostId, @elenaId, N'Я за освещение возле арки, вечером там действительно неуютно.'),
    (@suvorovaChatPostId, @tatyanaId, N'Еще нужны лавочки у детской площадки, сейчас их не хватает.'),
    (@suvorovaBushesPostId, @yuliyaId, N'Полностью согласна, обзор очень плохой после дождя и в сумерках.'),
    (@suvorovaBushesPostId, @ilyaId, N'Сформирую заявку на этой неделе и отпишусь о сроках.'),

    (@pushkinaCleanupPostId, @olgaId, N'Я приду и возьму с собой мешки для листвы.'),
    (@pushkinaCleanupPostId, @natalyaId, N'Отлично, давно хотели привести клумбы в порядок.'),
    (@pushkinaLightPostId, @viktorId, N'Поддерживаю, вечером там темно даже возле входа в арку.'),
    (@pushkinaLightPostId, @mikhailId, N'Добавим это в коллективное обращение, думаю, шансы хорошие.');

INSERT INTO invite_codes (house_id, code, created_by, expires_at)
VALUES
    (@arbekovoHouseId, N'PENZA-ARB-01', @alexeyId, DATEADD(DAY, 60, SYSDATETIME())),
    (@moskovskayaHouseId, N'PENZA-MSK-02', @dmitryId, DATEADD(DAY, 60, SYSDATETIME())),
    (@suvorovaHouseId, N'PENZA-SUV-03', @ilyaId, DATEADD(DAY, 60, SYSDATETIME())),
    (@pushkinaHouseId, N'PENZA-PSH-04', @mikhailId, DATEADD(DAY, 60, SYSDATETIME()));

INSERT INTO chat_messages (house_id, author_id, content, image_url)
VALUES
    (@arbekovoHouseId, @alexeyId, N'Добрый вечер, соседи. Вынес повестку по субботнему собранию в ленту.', NULL),
    (@arbekovoHouseId, @mariaId, N'Спасибо, я подготовлю список предложений по цветникам.', NULL),
    (@arbekovoHouseId, @sergeyId, N'И еще прошу не забыть про покрытие на детской площадке.', NULL),

    (@moskovskayaHouseId, @dmitryId, N'Коллеги, напомню про отключение воды во вторник.', NULL),
    (@moskovskayaHouseId, @ekaterinaId, N'Хорошо, тогда перенесем полив клумб на вечер.', NULL),
    (@moskovskayaHouseId, @andreyId, N'И давайте соберем идеи по велостойке в одном треде.', NULL),

    (@suvorovaHouseId, @ilyaId, N'Запускаем обсуждение благоустройства прямо здесь, чтобы ничего не потерялось.', NULL),
    (@suvorovaHouseId, @elenaId, N'За освещение у арки и за новые лавочки возле детской площадки.', NULL),
    (@suvorovaHouseId, @konstantinId, N'И еще нужна подрезка кустов у парковки, там плохой обзор.', NULL),

    (@pushkinaHouseId, @mikhailId, N'В воскресенье встречаемся на уборку двора в 10:30.', NULL),
    (@pushkinaHouseId, @olgaId, N'Я принесу перчатки и пару лишних мешков.', NULL),
    (@pushkinaHouseId, @alinaId, N'После уборки предлагаю обсудить освещение у арки.', NULL);

COMMIT TRANSACTION;
