CREATE TABLE referrals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    referral_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    active boolean DEFAULT FALSE,
    cost FLOAT DEFAULT 500,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (referral_id) REFERENCES users(referral_code)
);
