CREATE TABLE comments
(
    id          serial    not null unique,
    post_id     integer   not null,
    user_id     integer   not null,
    author      varchar(255) not null,
    content     text      not null,
    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now()
);

ALTER TABLE comments
    ADD CONSTRAINT fk_comments_post
        FOREIGN KEY (post_id) REFERENCES posts (id)
        ON DELETE CASCADE;

ALTER TABLE comments
    ADD CONSTRAINT fk_comments_user
        FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE;

CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_user_id ON comments (user_id);
