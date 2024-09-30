package schema

import (
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
		field.Uint64("bytes_sent").Optional(),

		// Total bytes received during the invocation.
		field.Uint64("bytes_recv").Optional(),

		// Total packets sent during the invocation.
		field.Uint64("packets_sent").Optional(),

		// Total packets received during the invocation
		field.Uint64("packets_recv").Optional(),

		// Peak bytes/sec sent during the invocation.
		field.Uint64("peak_bytes_sent_per_sec").Optional(),

		// Peak bytes/sec received during the invocation.
		field.Uint64("peak_bytes_recv_per_sec").Optional(),

		// Peak packets/sec sent during the invocation.
		field.Uint64("peak_packets_sent_per_sec").Optional(),

		// Peak packets/sec received during the invocation.
		field.Uint64("peak_packets_recv_per_sec").Optional(),
	}
}

// Edges of the SystemNetworkStats.
func (SystemNetworkStats) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge back to the network metrics object.
		edge.From("network_metrics", NetworkMetrics.Type).Ref("system_network_stats").Unique(),
	}
}
