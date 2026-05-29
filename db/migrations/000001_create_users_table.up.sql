CREATE TYPE transaction_status AS ENUM ('PENDING','SUCCESS','FAILED');
CREATE TYPE transaction_type AS ENUM ('TOPUP', 'TRANSFER');
CREATE TYPE cashflow_type AS ENUM ('income','expense');

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