package schema

import (
	"database/sql"
	"database/sql/driver"
	"strconv"

	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

var postgresUint64SchemaType = map[string]string{
	dialect.Postgres: "NUMERIC(20,0)",
}

// uint64ValueScanner keeps Ent fields as native uint64 values while storing
// them as PostgreSQL numerics, which can represent the full uint64 range.
var uint64ValueScanner = field.ValueScannerFunc[uint64, *sql.NullString]{
	V: func(value uint64) (driver.Value, error) {
		return strconv.FormatUint(value, 10), nil
	},
	S: func(value *sql.NullString) (uint64, error) {
		if !value.Valid {
			return 0, nil
		}
		return strconv.ParseUint(value.String, 10, 64)
	},
}
