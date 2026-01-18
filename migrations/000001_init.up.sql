CREATE TABLE users
(
    id serial not null unique, -- поле id будет автоматически создано и уникальным
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE posts
(
    id              serial       not null unique,
    user_id         serial      not null,
    title           varchar(255) not null,
    description     text         not null,
    created_at      timestamp    not null default now(),
        updated_at      timestamp    not null default now()
    );

ALTER TABLE posts
    ADD CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE;

