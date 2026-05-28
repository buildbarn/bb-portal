package schema

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/99designs/gqlgen/graphql"
)

type Uint64Numeric uint64

func (v Uint64Numeric) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.FormatUint(uint64(v), 10))
}

func (v *Uint64Numeric) UnmarshalGQL(src any) error {
	parsed, err := graphql.UnmarshalUint64(src)
	if err != nil {
		return err
	}
	*v = Uint64Numeric(parsed)
	return nil
}

func (v Uint64Numeric) Value() (driver.Value, error) {
	return strconv.FormatUint(uint64(v), 10), nil
}

func (v *Uint64Numeric) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		*v = 0
		return nil
	case int64:
		if src < 0 {
			return fmt.Errorf("cannot scan negative int64 %d into Uint64Numeric", src)
		}
		*v = Uint64Numeric(src)
		return nil
	case string:
		parsed, err := strconv.ParseUint(src, 10, 64)
		if err != nil {
			return err
		}
		*v = Uint64Numeric(parsed)
		return nil
	case []byte:
		parsed, err := strconv.ParseUint(string(src), 10, 64)
		if err != nil {
			return err
		}
		*v = Uint64Numeric(parsed)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Uint64Numeric", src)
	}
}

var postgresUint64SchemaType = map[string]string{
	dialect.Postgres: "numeric",
}

// SystemNetworkStats holds the schema definition for the SystemNetworkStats entity.
type SystemNetworkStats struct {
	ent.Schema
}

// Fields of the SystemNetworkStats.
func (SystemNetworkStats) Fields() []ent.Field {
	return []ent.Field{
		// Total bytes sent during the invocation.
		field.Uint64("bytes_sent").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Total bytes received during the invocation.
		field.Uint64("bytes_recv").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Total packets sent during the invocation.
		field.Uint64("packets_sent").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Total packets received during the invocation
		field.Uint64("packets_recv").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Peak bytes/sec sent during the invocation.
		field.Uint64("peak_bytes_sent_per_sec").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Peak bytes/sec received during the invocation.
		field.Uint64("peak_bytes_recv_per_sec").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Peak packets/sec sent during the invocation.
		field.Uint64("peak_packets_sent_per_sec").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),

		// Peak packets/sec received during the invocation.
		field.Uint64("peak_packets_recv_per_sec").
			GoType(Uint64Numeric(0)).
			Optional().
			SchemaType(postgresUint64SchemaType),
	}
}

// Edges of the SystemNetworkStats.
func (SystemNetworkStats) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the network metrics object.
		edge.From("network_metrics", NetworkMetrics.Type).
			Ref("system_network_stats").
			Unique(),
	}
}

// Indexes of the SystemNetworkStats.
func (SystemNetworkStats) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("network_metrics"),
	}
}

// Mixin of the SystemNetworkStats.
func (SystemNetworkStats) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
