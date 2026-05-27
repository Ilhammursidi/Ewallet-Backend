-- =========================
-- PAYMENT METHODS
-- =========================
INSERT INTO payment_methods (payment_name) VALUES
('Bank Rakyat Indonesia'),
('Dana'),
('Bank Central Asia'),
('Gopay'),
('Ovo');


-- =========================
-- USERS (50 DATA)
-- isVerified DIABAIKAN
-- =========================
INSERT INTO users (
    email,
    password,
    pin,
    fullname,
    phone_number
)
SELECT
    'user' || gs || '@gmail.com',
    '$2a$10$examplehashedpassword',
    '$2a$10$examplehashedpin',
    'User ' || gs,
    '08' || LPAD((100000000 + gs)::TEXT, 10, '0')
FROM generate_series(1, 50) AS gs;


-- =========================
-- WALLET
-- SETIAP USER PUNYA WALLET
-- =========================
INSERT INTO wallet (
    user_id,
    balance
)
SELECT
    id,
    (RANDOM() * 5000000 + 100000)::INT
FROM users
WHERE id <= 50;


-- =========================
-- TRANSACTIONS
-- 50 DATA RANDOM
-- =========================
INSERT INTO transactions (
    user_id,
    sender_wallet_id,
    receiver_wallet_id,
    payment_method_id,
    type,
    flow_type,
    amount,
    status
)
SELECT
    ((RANDOM() * 49) + 1)::INT,
    
    CASE
        WHEN RANDOM() > 0.5
        THEN ((RANDOM() * 49) + 1)::INT
        ELSE NULL
    END,

    ((RANDOM() * 49) + 1)::INT,

    CASE
        WHEN RANDOM() > 0.5
        THEN ((RANDOM() * 4) + 1)::INT
        ELSE NULL
    END,

    CASE
        WHEN RANDOM() < 0.33 THEN 'TOPUP'::transaction_type
        WHEN RANDOM() < 0.66 THEN 'TRANSFER_IN'::transaction_type
        ELSE 'TRANSFER_OUT'::transaction_type
    END,

    CASE
        WHEN RANDOM() < 0.5 THEN 'income'::cashflow_type
        ELSE 'expense'::cashflow_type
    END,

    ((RANDOM() * 1000000) + 10000)::INT,

    'SUCCESS'::transaction_status

FROM generate_series(1, 50);


-- =========================
-- TRANSFER DETAILS
-- HANYA UNTUK TRANSFER
-- =========================
INSERT INTO transfer_details (
    transaction_id,
    sender_wallet_id,
    receiver_wallet_id
)
SELECT
    id,
    sender_wallet_id,
    receiver_wallet_id
FROM transactions
WHERE type IN ('TRANSFER_IN', 'TRANSFER_OUT')
AND sender_wallet_id IS NOT NULL
AND receiver_wallet_id IS NOT NULL;


-- =========================
-- TOPUP DETAILS
-- HANYA UNTUK TOPUP
-- =========================
INSERT INTO topup_details (
    transaction_id,
    wallet_id,
    payment_method_id,
    order_amount,
    tax_amount,
    delivery_fee,
    total_amount,
    status
)
SELECT
    t.id,
    t.receiver_wallet_id,
    t.payment_method_id,

    t.amount,

    (t.amount * 0.02)::INT,

    1000,

    t.amount + (t.amount * 0.02)::INT + 1000,

    'SUCCESS'::transaction_status

FROM transactions t
WHERE t.type = 'TOPUP';