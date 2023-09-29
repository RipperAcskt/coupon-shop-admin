create table if not exists media(
    id varchar PRIMARY KEY,
    coupon_id varchar REFERENCES coupons(id) ON DELETE CASCADE,
    path varchar not null
);