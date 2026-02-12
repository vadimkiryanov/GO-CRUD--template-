CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('like', 'dislike')),  -- none не храним
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, post_id)  -- один голос на пост от юзера
);

-- Индексы создаются для ускорения запросов
CREATE INDEX idx_votes_post ON votes(post_id);
CREATE INDEX idx_votes_user_post ON votes(user_id, post_id);
