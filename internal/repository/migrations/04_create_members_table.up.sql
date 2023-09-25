CREATE TABLE IF NOT EXISTS members (
        id VARCHAR PRIMARY KEY,
        email VARCHAR UNIQUE NOT NULL,
        first_name VARCHAR NOT NULL,
        second_name VARCHAR NOT NULL,
        organization_id VARCHAR NOT NULL REFERENCES organization(id)
);