CREATE TABLE
    IF NOT EXISTS core.apps (
        id UUID PRIMARY KEY,
        name BYTEA NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

-- INSERT INTO
--     core.apps (name)
-- SELECT
--     ('fransiscushermanto')
-- WHERE
--     NOT EXISTS (
--         SELECT
--             1
--         FROM
--             core.apps
--         WHERE
--             name = 'fransiscushermanto'
--     );
-- INSERT INTO
--     core.apps (name)
-- SELECT
--     ('bloomify')
-- WHERE
--     NOT EXISTS (
--         SELECT
--             1
--         FROM
--             core.apps
--         WHERE
--             name = 'bloomify'
--     );