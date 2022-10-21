CREATE TABLE accounts
(
    id SERIAL PRIMARY KEY,
    user_id int not null,
    curr_amount bigint,
    pending_amount bigint,
    last_updated timestamp
);

CREATE TABLE transactions_log
(
    id SERIAL PRIMARY KEY,
    account_id_from int not null REFERENCES accounts(id),
    account_id_to int not null REFERENCES accounts(id),
    transaction_sum bigint,
    status varchar(255) not null,
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
    updated_at timestamp,
);

INSERT INTO accounts(user_id, curr_amount, pending_amount, last_updated)
VALUES (2233, 500, 0, CURRENT_TIMESTAMP),
       (2237, 800, 0, CURRENT_TIMESTAMP);
