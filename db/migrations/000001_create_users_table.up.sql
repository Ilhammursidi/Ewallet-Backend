CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    pin VARCHAR(255),
    fullname VARCHAR(255),
    photo_path VARCHAR(255),
    phone_number VARCHAR(255) UNIQUE,
    isVerified BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);