CREATE SCHEMA IF NOT EXISTS comments;

CREATE TABLE IF NOT EXISTS comments.posts (
    id SERIAL PRIMARY KEY,
    parent_id INT NOT NULL,
    news_id INT NOT NULL,
    text TEXT NOT NULL,
    pub_time BIGINT NOT NULL CHECK (pubTime > 0)
)