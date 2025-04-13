CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    lead_id integer NOT NULL,
    user_id integer NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (lead_id) REFERENCES leads(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
