CREATE TABLE topup_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT UNIQUE NOT NULL,
    wallet_id INT NOT NULL,
    payment_method_id INT NOT NULL,
    order_amount INT NOT NULL,
    tax_amount INT NOT NULL,
    delivery_fee INT NOT NULL,
    total_amount INT NOT NULL,
    status transaction_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (wallet_id) REFERENCES wallet(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);