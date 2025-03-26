CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id)
);