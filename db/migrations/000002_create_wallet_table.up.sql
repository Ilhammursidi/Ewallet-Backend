CREATE TABLE wallet (
    id SERIAL PRIMARY KEY,  
    user_id INT UNIQUE NOT NULL,
    balance INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);