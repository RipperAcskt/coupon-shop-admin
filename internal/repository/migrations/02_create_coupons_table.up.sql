CREATE TABLE IF NOT EXISTS coupons (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    price int NOT NULL,
    level int NOT NULL
);