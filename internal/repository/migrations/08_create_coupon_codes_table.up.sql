CREATE TABLE IF NOT EXISTS coupon_codes (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    code varchar NOT NULL UNIQUE,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);
