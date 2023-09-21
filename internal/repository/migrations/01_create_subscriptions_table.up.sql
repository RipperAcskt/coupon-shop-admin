CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    price int NOT NULL,
    level int PRIMARY KEY
);