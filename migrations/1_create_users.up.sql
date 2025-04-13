DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('user', 'admin', 'manager');
    END IF;
END $$;

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
   is_active BOOLEAN DEFAULT false,
   roles user_role[] DEFAULT ARRAY['user'::user_role]
);
