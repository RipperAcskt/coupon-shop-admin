CREATE TABLE IF NOT EXISTS users (
    id VARCHAR PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    code VARCHAR(8)  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    subscription VARCHAR,
    subscription_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + '1 month'::interval
);

CREATE FUNCTION expire_table_delete_old_rows() RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE users SET subscription = '' WHERE subscription_time < NOW() - INTERVAL '1 minute' AND subscription != '';
    RETURN NEW;
END;
$$;

CREATE TRIGGER expire_table_insert_trigger
    AFTER INSERT ON users
    FOR EACH ROW
EXECUTE PROCEDURE expire_table_delete_old_rows();

CREATE TRIGGER expire_table_update_trigger
    AFTER UPDATE ON users
    FOR EACH ROW
EXECUTE PROCEDURE expire_table_delete_old_rows();
