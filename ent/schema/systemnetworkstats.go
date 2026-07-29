package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SystemNetworkStats holds the schema definition for the SystemNetworkStats entity.
type SystemNetworkStats struct {
	ent.Schema
}

// Fields of the SystemNetworkStats.
func (SystemNetworkStats) Fields() []ent.Field {
	return []ent.Field{
		// Total bytes sent during the invocation.
		field.Uint64("bytes_sent").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
			SchemaType(postgresUint64SchemaType),

		// Total bytes received during the invocation.
		field.Uint64("bytes_recv").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
			SchemaType(postgresUint64SchemaType),

		// Total packets sent during the invocation.
		field.Uint64("packets_sent").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
			SchemaType(postgresUint64SchemaType),

		// Total packets received during the invocation
		field.Uint64("packets_recv").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
			SchemaType(postgresUint64SchemaType),

		// Peak bytes/sec sent during the invocation.
		field.Uint64("peak_bytes_sent_per_sec").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
			SchemaType(postgresUint64SchemaType),

		// Peak bytes/sec received during the invocation.
		field.Uint64("peak_bytes_recv_per_sec").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
			SchemaType(postgresUint64SchemaType),

		// Peak packets/sec sent during the invocation.
		field.Uint64("peak_packets_sent_per_sec").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
			SchemaType(postgresUint64SchemaType),

		// Peak packets/sec received during the invocation.
		field.Uint64("peak_packets_recv_per_sec").
			Annotations(entgql.Type("UnsignedLong")).
			ValueScanner(uint64ValueScanner).
			Optional().
			Immutable().
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
	return []ent.Index{}
}

// Mixin of the SystemNetworkStats.
func (SystemNetworkStats) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Int64IdMixin{},
	}
}
