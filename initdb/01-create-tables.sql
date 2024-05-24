-- Установка расширения uuid-ossp
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание таблицы posts
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    is_disabled_comments BOOLEAN DEFAULT FALSE,
    user_id TEXT NOT NULL
);

-- Создание таблицы comments
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    parent_id UUID,
    body TEXT NOT NULL,
    user_id TEXT NOT NULL,
    CONSTRAINT fk_post
      FOREIGN KEY(post_id) 
      REFERENCES posts(id)
      ON DELETE CASCADE,
    CONSTRAINT fk_parent_comment
      FOREIGN KEY(parent_id)
      REFERENCES comments(id)
      ON DELETE CASCADE
);

-- Создание индексов
CREATE INDEX idx_post_id ON comments(post_id);
CREATE INDEX idx_parent_id ON comments(parent_id);
