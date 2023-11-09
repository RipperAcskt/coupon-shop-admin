CREATE TABLE IF NOT EXISTS coupons (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    price int NOT NULL,
    percent int NOT NULL,
    level int NOT NULL,
    region VARCHAR NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    organization_id VARCHAR NOT NULL REFERENCES organization(id)
);