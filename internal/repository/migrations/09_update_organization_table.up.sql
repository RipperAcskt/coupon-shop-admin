
ALTER TABLE organization
ADD COLUMN ORGN VARCHAR NOT NULL,
ADD COLUMN KPP VARCHAR NOT NULL,
ADD COLUMN INN VARCHAR NOT NULL,
ADD COLUMN address VARCHAR NOT NULL;

create table if not exists images(
    id varchar PRIMARY KEY,
    organization_id varchar REFERENCES organization(id) ON DELETE CASCADE,
    path varchar not null
);