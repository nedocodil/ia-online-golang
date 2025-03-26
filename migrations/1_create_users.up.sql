CREATE TABLE users (
   id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    telegram VARCHAR(255),
    city VARCHAR(100),
    password_hash TEXT NOT NULL,
    referral_code VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active boolean DEFAULT false,
    role TEXT NOT NULL,
    reward_internet FLOAT,
    reward_cleaning FLOAT,
    reward_shipping FLOAT
);