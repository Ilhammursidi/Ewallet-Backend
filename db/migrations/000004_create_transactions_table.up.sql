CREATE TYPE cashflow_type AS ENUM ('income','expense');
CREATE TYPE transaction_type AS ENUM ('TOPUP', 'TRANSFER_IN', 'TRANSFER_OUT');
CREATE TYPE transaction_status AS ENUM ('PENDING','SUCCESS','FAILED');

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    sender_wallet_id INT,
    receiver_wallet_id INT,
    payment_method_id INT,
    type transaction_type,
    flow_type cashflow_type,
    amount INT NOT NULL,
    status transaction_status DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (sender_wallet_id) REFERENCES wallet(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (receiver_wallet_id) REFERENCES wallet(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);