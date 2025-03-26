CREATE TABLE leads (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    fio VARCHAR(255),
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    internet BOOLEAN NOT NULL,
    cleaning BOOLEAN NOT NULL,
    shipping BOOLEAN NOT NULL,
    address VARCHAR(255),
    status INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    payment_at TIMESTAMP,
    FOREIGN KEY (status) REFERENCES statuses(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
