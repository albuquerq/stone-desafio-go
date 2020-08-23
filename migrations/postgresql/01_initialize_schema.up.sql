CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS accounts(
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    secret VARCHAR(200) NOT NULL,
    balance BIGINT CHECK (balance > 0),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transfers(
    id UUID PRIMARY KEY,
    account_origin_id UUID REFERENCES accounts (id),
    account_destination_id UUID REFERENCES accounts (id),
    amount BIGINT NOT NULL CHECK (amount > -1),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
);