CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    price int NOT NULL,
    level int PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS organization (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    email_admin VARCHAR NOT NULL,
    level_subscription INT
);