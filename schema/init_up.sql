CREATE TABLE accounts
(
    id SERIAL PRIMARY KEY,
    user_id int not null,
    curr_amount bigint,
    pending_amount bigint,
    last_updated timestamp,
    constraint curr_amount_non_negative check (curr_amount >= 0)
);

CREATE TABLE transactions_log
(
    id SERIAL PRIMARY KEY,
    account_id_from int not null REFERENCES accounts(id),
    account_id_to int REFERENCES accounts(id),
    transaction_sum bigint,
    status varchar(255) not null,
    event_type varchar(255) not null,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE service_log
(
    id SERIAL PRIMARY KEY,
    account_id int not null REFERENCES accounts(id),
    invoice bigint,
    service_id int not null,
    order_id int not null,
    status varchar(255) not null,
    created_at timestamp,
    updated_at timestamp
);

CREATE VIEW TransactionsByAccount AS
SELECT account_id_from AS account_id, account_id_to, transaction_sum, NULL AS service_id, NULL AS order_id, status, event_type, created_at::TIMESTAMP WITHOUT TIME ZONE,  updated_at::TIMESTAMP WITHOUT TIME ZONE
FROM transactions_log WHERE status='Completed'
UNION ALL
SELECT account_id, NULL as account_id_to, invoice as transaction_sum, service_id, order_id, status, 'Service-purchase' AS event_type,
       created_at, updated_at
FROM service_log WHERE status='Approved';