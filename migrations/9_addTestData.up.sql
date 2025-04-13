-- Заполнение таблицы statuses
INSERT INTO statuses (id, name, bitrix_name) VALUES
    (0, 'new', 'Новая заявка'),
    (1, 'noContact', 'Недозвон'),
    (2, 'pending', 'Отложена'),
    (3, 'scheduled', 'Назначена'),
    (4, 'ready', 'Готова'),
    (5, 'paid', 'Оплачено'),
    (6, 'refusal', 'Отказ'),
    (7, 'appointment_control', 'Контроль назначения');

-- Заполнение таблицы users
INSERT INTO users (phone_number, email, name, telegram, city, password_hash, referral_code, is_active, roles) VALUES
    ('+1234567890', 'user1@example.com', 'Иван Иванов', '@ivanov', 'Москва', 'hash1', 'REF123', true, '{"user"}'),
    ('+0987654321', 'user2@example.com', 'Петр Петров', '@petrov', 'Санкт-Петербург', 'hash2', 'REF456', false, '{"admin"}'),
    ('+79963208153','vladik00848@gmail.com','Владислав','@bymaginn','г Тюмень',	'$2a$10$JlVRln/ZAipX5NFtDGN7G.THsoK.E/I89ECiuM9NL65H8EWNOVv1W','0faa472d-265a-4b09-bd34-5ae05aacb3b9',true, '{"user"}'),	
    ('+1122334455', 'user3@example.com', 'Сидоров Алексей', '@sidorov', 'Казань', 'hash3', 'REF789', true, '{"user", "manager"}'),
    ('+71112223344', 'user4@example.com', 'Алина Новикова', '@alinanov', 'Екатеринбург', 'hash4', 'REF101', true, '{"user"}'),
    ('+72223334455', 'user5@example.com', 'Григорий Лебедев', '@grigleb', 'Новосибирск', 'hash5', 'REF102', true, '{"user"}'),
    ('+73334445566', 'user6@example.com', 'Екатерина Цветкова', '@katetsv', 'Нижний Новгород', 'hash6', 'REF103', true, '{"user"}'),
    ('+74445556677', 'user7@example.com', 'Михаил Емельянов', '@mikeemel', 'Ростов-на-Дону', 'hash7', 'REF104', true, '{"manager"}'),
    ('+75556667788', 'user8@example.com', 'Татьяна Павлова', '@tatyapav', 'Самара', 'hash8', 'REF105', true, '{"user"}');

-- Заполнение таблицы leads
INSERT INTO leads (
    user_id, fio, phone_number, internet, cleaning, shipping, address,
    status_id, reward_internet, reward_cleaning, reward_shipping,
    created_at, completed_at, payment_at
) VALUES
(1, 'Анна Смирнова', '+79991234567', true, false, true, 'ул. Ленина, д. 10, Москва', 0, 500, 0, 700, NOW(), NULL, NULL),
(2, 'Василий Кузнецов', '+79997654321', false, true, false, 'ул. Пушкина, д. 5, Санкт-Петербург', 3, 0, 800, 0, NOW(), NULL, NULL),
(3, 'Мария Федорова', '+79994567890', true, true, true, 'ул. Гагарина, д. 20, Казань', 5, 600, 400, 1000, NOW(), NOW(), NOW()),
(1, 'Александр Соловьев', '+79991112233', false, false, true, 'ул. Чехова, д. 15, Москва', 2, 0, 0, 900, NOW(), NULL, NULL),
(2, 'Екатерина Орлова', '+79993334455', true, true, false, 'ул. Кирова, д. 8, Санкт-Петербург', 4, 450, 750, 0, NOW(), NOW(), NULL),
(3, 'Дмитрий Павлов', '+79994455667', false, true, true, 'ул. Лермонтова, д. 3, Казань', 1, 0, 600, 850, NOW(), NULL, NULL),
(1, 'Сергей Миронов', '+79995566778', true, false, false, 'ул. Гоголя, д. 12, Москва', 6, 550, 0, 0, NOW(), NULL, NULL),
(2, 'Ольга Зайцева', '+79992223344', true, true, true, 'ул. Тверская, д. 22, Санкт-Петербург', 5, 700, 500, 900, NOW(), NOW(), NOW()),
(3, 'Игорь Ковалев', '+79997778899', false, false, true, 'ул. Октябрьская, д. 7, Казань', 3, 0, 0, 650, NOW(), NULL, NULL),
(1, 'Наталья Белова', '+79998889900', true, false, true, 'ул. Советская, д. 18, Москва', 2, 600, 0, 700, NOW(), NULL, NULL),
(2, 'Максим Романов', '+79990001122', false, true, false, 'ул. Дзержинского, д. 5, Санкт-Петербург', 0, 0, 850, 0, NOW(), NULL, NULL),
(3, 'Антон Васильев', '+79991112244', true, true, true, 'ул. Садовая, д. 30, Казань', 4, 750, 500, 950, NOW(), NOW(), NULL),
(1, 'Ксения Тихонова', '+79992223355', false, false, false, 'ул. Куйбышева, д. 9, Москва', 1, 0, 0, 0, NOW(), NULL, NULL),
(2, 'Юлия Михайлова', '+79993334466', true, false, true, 'ул. Фрунзе, д. 14, Санкт-Петербург', 3, 650, 0, 800, NOW(), NULL, NULL),
(3, 'Павел Седов', '+79994455688', false, true, false, 'ул. Красная, д. 11, Казань', 6, 0, 700, 0, NOW(), NULL, NULL),
(1, 'Елена Кравцова', '+79995566799', true, true, false, 'ул. Лесная, д. 25, Москва', 5, 550, 850, 0, NOW(), NOW(), NOW()),
(2, 'Алексей Игнатьев', '+79996677800', false, false, true, 'ул. Школьная, д. 16, Санкт-Петербург', 2, 0, 0, 750, NOW(), NULL, NULL),
(3, 'Вероника Фролова', '+79997788911', true, false, false, 'ул. Полевая, д. 19, Казань', 4, 620, 0, 0, NOW(), NOW(), NULL);



-- Примеры рефералов (user_id ссылается на id пользователя, а referral_id — на referral_code другого пользователя)
INSERT INTO referrals (user_id, referral_id) VALUES
    (4, 'REF123'), 
    (5, 'REF456'), 
    (6, '0faa472d-265a-4b09-bd34-5ae05aacb3b9'), 
    (7, 'REF789'), 
    (8, 'REF101'); 

