DO $$
DECLARE
    r RECORD;
BEGIN
    -- Loop through all tables in the current schema
    FOR r IN SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename != 'accounts' LOOP
        -- Execute the drop table command
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;

-- Create the accounts table
CREATE TABLE IF NOT EXISTS accounts(
    id CHAR(27) PRIMARY KEY,
    name VARCHAR(24) NOT NULL
);
