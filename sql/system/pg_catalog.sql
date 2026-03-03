-- These create table statement informs sqlc about the existence of
-- postgres specific tables, it is not intended to be run anywhere.
CREATE TABLE pg_namespace (
    oid oid PRIMARY KEY,
    nspname name NOT NULL
);

CREATE TABLE pg_class (
    oid oid PRIMARY KEY,
    relname name NOT NULL,
    relnamespace oid NOT NULL,
    relpages integer NOT NULL,
    relispartition boolean NOT NULL
);

CREATE TABLE pg_attribute (
    attrelid oid NOT NULL,
    attname name NOT NULL,
    attnum smallint NOT NULL
);

CREATE TABLE pg_constraint (
    oid oid PRIMARY KEY,
    conname name NOT NULL,
    connamespace oid NOT NULL,
    contype char NOT NULL,
    conrelid oid NOT NULL,
    confrelid oid NOT NULL,
    conkey smallint[]
);

CREATE TABLE pg_index (
    indexrelid oid PRIMARY KEY,
    indrelid oid NOT NULL,
    indisvalid boolean NOT NULL,
    indkey int2vector NOT NULL
);
