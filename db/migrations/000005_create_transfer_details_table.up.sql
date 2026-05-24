CREATE TABLE transfer_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT UNIQUE NOT NULL,
    sender_wallet_id INT NOT NULL,
    receiver_wallet_id INT NOT NULL,
    -- amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (sender_wallet_id) REFERENCES wallet(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (receiver_wallet_id) REFERENCES wallet(id)
);