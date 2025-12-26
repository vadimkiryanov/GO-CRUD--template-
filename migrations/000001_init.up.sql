CREATE TABLE IF NOT EXISTS example ( -- создается таблица example
    id SERIAL PRIMARY KEY, -- id поля создается автоматически
    name VARCHAR(255) NOT NULL, -- поле name с длиной от 0 до 255 символов
    email VARCHAR(255) NOT NULL, -- поле email с длиной от 0 до 255 символов
    password VARCHAR(255) NOT NULL, -- поле password с длиной от 0 до 255 символов
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- поле created_at создается автоматически
);