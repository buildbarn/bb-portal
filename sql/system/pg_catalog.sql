-- This create table statement informs sqlc about the existence of
-- postgres specific tables, it is not intended to be run anywhere.
CREATE TABLE pg_class (
    oid oid NOT NULL,
    relname name NOT NULL,
    relpages integer NOT NULL,
    relispartition boolean NOT NULL
);
