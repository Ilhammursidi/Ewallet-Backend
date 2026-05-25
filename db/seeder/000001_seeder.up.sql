-- ============================================================
-- SEED DATA
-- ============================================================

-- 1. PAYMENT METHODS
INSERT INTO payment_methods (payment_name) VALUES
('Bank Rakyat Indonesia'),
('Dana'),
('Bank Central Asia'),
('Gopay'),
('Ovo');

-- 2. USERS
INSERT INTO users (email, password, pin, fullname, photo_path, phone_number, isVerified) VALUES
('alice@mail.com',   '$2a$10$abcdefghijklmnopqrstuuVwXyZ012345678901234567890123456', '123456', 'Alice Wonderland', 'uploads/profiles/alice.jpg',   '081111111111', true),
('bob@mail.com',     '$2a$10$abcdefghijklmnopqrstuuVwXyZ012345678901234567890123456', '234567', 'Bob Marley',       'uploads/profiles/bob.jpg',     '082222222222', true),
('charlie@mail.com', '$2a$10$abcdefghijklmnopqrstuuVwXyZ012345678901234567890123456', '345678', 'Charlie Chaplin',  'uploads/profiles/charlie.jpg', '083333333333', true),
('diana@mail.com',   '$2a$10$abcdefghijklmnopqrstuuVwXyZ012345678901234567890123456', '456789', 'Diana Prince',    'uploads/profiles/diana.jpg',   '084444444444', true),
('eve@mail.com',     '$2a$10$abcdefghijklmnopqrstuuVwXyZ012345678901234567890123456', '567890', 'Eve Online',      'uploads/profiles/eve.jpg',     '085555555555', true);

-- 3. WALLETS (1 wallet per user)
-- user_id: alice=1, bob=2, charlie=3, diana=4, eve=5
INSERT INTO wallet (user_id, balance) VALUES
(1, 5000000),
(2, 3000000),
(3, 1500000),
(4, 750000),
(5, 2000000);

-- ============================================================
-- 4. TOPUP TRANSACTIONS
-- wallet_id: alice=1, bob=2, charlie=3, diana=4, eve=5
-- ============================================================

-- Alice topup via BCA
INSERT INTO transactions (user_id, payment_method_id, type, flow_type, amount, status) VALUES
(1, 1, 'TOPUP', 'income', 1000000, 'SUCCESS');  -- tx id=1

INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status) VALUES
(1, 1, 1, 1000000, 100000, 5000, 1105000, 'SUCCESS');

-- Bob topup via Mandiri
INSERT INTO transactions (user_id, payment_method_id, type, flow_type, amount, status) VALUES
(2, 2, 'TOPUP', 'income', 500000, 'SUCCESS');   -- tx id=2

INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status) VALUES
(2, 2, 2, 500000, 50000, 5000, 555000, 'SUCCESS');

-- Charlie topup via OVO
INSERT INTO transactions (user_id, payment_method_id, type, flow_type, amount, status) VALUES
(3, 3, 'TOPUP', 'income', 750000, 'SUCCESS');   -- tx id=3

INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status) VALUES
(3, 3, 3, 750000, 75000, 5000, 830000, 'SUCCESS');

-- Diana topup via GoPay
INSERT INTO transactions (user_id, payment_method_id, type, flow_type, amount, status) VALUES
(4, 4, 'TOPUP', 'income', 200000, 'SUCCESS');   -- tx id=4

INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status) VALUES
(4, 4, 4, 200000, 20000, 5000, 225000, 'SUCCESS');

-- Eve topup via DANA
INSERT INTO transactions (user_id, payment_method_id, type, flow_type, amount, status) VALUES
(5, 5, 'TOPUP', 'income', 300000, 'SUCCESS');   -- tx id=5

INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status) VALUES
(5, 5, 5, 300000, 30000, 5000, 335000, 'SUCCESS');

-- Alice topup lagi via GoPay (PENDING)
INSERT INTO transactions (user_id, payment_method_id, type, flow_type, amount, status) VALUES
(1, 4, 'TOPUP', 'income', 2000000, 'PENDING');  -- tx id=6

INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status) VALUES
(6, 1, 4, 2000000, 200000, 5000, 2205000, 'PENDING');

-- Bob topup via DANA (FAILED)
INSERT INTO transactions (user_id, payment_method_id, type, flow_type, amount, status) VALUES
(2, 5, 'TOPUP', 'income', 1000000, 'FAILED');   -- tx id=7

INSERT INTO topup_details (transaction_id, wallet_id, payment_method_id, order_amount, tax_amount, delivery_fee, total_amount, status) VALUES
(7, 2, 5, 1000000, 100000, 5000, 1105000, 'FAILED');

-- ============================================================
-- 5. TRANSFER TRANSACTIONS
-- wallet_id: alice=1, bob=2, charlie=3, diana=4, eve=5
-- ============================================================

-- Alice → Bob (200rb)
INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(1, 1, 2, 'TRANSFER_OUT', 'expense', 200000, 'SUCCESS');  -- tx id=8

INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(2, 1, 2, 'TRANSFER_IN', 'income', 200000, 'SUCCESS');    -- tx id=9

INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id) VALUES
(8, 1, 2);

-- Bob → Charlie (100rb)
INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(2, 2, 3, 'TRANSFER_OUT', 'expense', 100000, 'SUCCESS');  -- tx id=10

INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(3, 2, 3, 'TRANSFER_IN', 'income', 100000, 'SUCCESS');    -- tx id=11

INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id) VALUES
(10, 2, 3);

-- Charlie → Diana (150rb)
INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(3, 3, 4, 'TRANSFER_OUT', 'expense', 150000, 'SUCCESS');  -- tx id=12

INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(4, 3, 4, 'TRANSFER_IN', 'income', 150000, 'SUCCESS');    -- tx id=13

INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id) VALUES
(12, 3, 4);

-- Diana → Eve (50rb)
INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(4, 4, 5, 'TRANSFER_OUT', 'expense', 50000, 'SUCCESS');   -- tx id=14

INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(5, 4, 5, 'TRANSFER_IN', 'income', 50000, 'SUCCESS');     -- tx id=15

INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id) VALUES
(14, 4, 5);

-- Eve → Alice (300rb)
INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(5, 5, 1, 'TRANSFER_OUT', 'expense', 300000, 'SUCCESS');  -- tx id=16

INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(1, 5, 1, 'TRANSFER_IN', 'income', 300000, 'SUCCESS');    -- tx id=17

INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id) VALUES
(16, 5, 1);

-- Alice → Charlie (500rb) FAILED
INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(1, 1, 3, 'TRANSFER_OUT', 'expense', 500000, 'FAILED');   -- tx id=18

INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(3, 1, 3, 'TRANSFER_IN', 'income', 500000, 'FAILED');     -- tx id=19

INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id) VALUES
(18, 1, 3);

-- Bob → Eve (250rb) PENDING
INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(2, 2, 5, 'TRANSFER_OUT', 'expense', 250000, 'PENDING');  -- tx id=20

INSERT INTO transactions (user_id, sender_wallet_id, receiver_wallet_id, type, flow_type, amount, status) VALUES
(5, 2, 5, 'TRANSFER_IN', 'income', 250000, 'PENDING');    -- tx id=21

INSERT INTO transfer_details (transaction_id, sender_wallet_id, receiver_wallet_id) VALUES
(20, 2, 5);