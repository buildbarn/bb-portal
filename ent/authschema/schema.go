package authschema

import (
	"github.com/buildbarn/bb-portal/ent/schema"
)

type (
	// Action reexport with auth policy added
	Action struct{ schema.Action }
	// ActionCacheStatistics reexport with auth policy added
	ActionCacheStatistics struct{ schema.ActionCacheStatistics }
	// ActionData reexport with auth policy added
	ActionData struct{ schema.ActionData }
	// ActionSummary reexport with auth policy added
	ActionSummary struct{ schema.ActionSummary }
	// ArtifactMetrics reexport with auth policy added
	ArtifactMetrics struct{ schema.ArtifactMetrics }
	// AuthenticatedUser reexport with auth policy added
	AuthenticatedUser struct{ schema.AuthenticatedUser }
	// BazelInvocation reexport with auth policy added
	BazelInvocation struct{ schema.BazelInvocation }
	// Build reexport with auth policy added
	Build struct{ schema.Build }
	// BuildLogChunk reexport with auth policy added
	BuildLogChunk struct{ schema.BuildLogChunk }
	// BuildGraphMetrics reexport with auth policy added
	BuildGraphMetrics struct{ schema.BuildGraphMetrics }
	// Configuration reexport with auth policy added
	Configuration struct{ schema.Configuration }
	// ConnectionMetadata reexport with auth policy added
	ConnectionMetadata struct{ schema.ConnectionMetadata }
	// CumulativeMetrics reexport with auth policy added
	CumulativeMetrics struct{ schema.CumulativeMetrics }
	// EvaluationStat reexport with auth policy added
	EvaluationStat struct{ schema.EvaluationStat }
	// EventMetadata reexport with auth policy added
	EventMetadata struct{ schema.EventMetadata }
	// GarbageMetrics reexport with auth policy added
	GarbageMetrics struct{ schema.GarbageMetrics }
	// IncompleteBuildLog reexport with auth policy added
	IncompleteBuildLog struct{ schema.IncompleteBuildLog }
	// InstanceName reexport with auth policy added
	InstanceName struct{ schema.InstanceName }
	// InvocationFiles reexport with auth policy added
	InvocationFiles struct{ schema.InvocationFiles }
	// InvocationTarget reexport with auth policy added
	InvocationTarget struct{ schema.InvocationTarget }
	// MemoryMetrics reexport with auth policy added
	MemoryMetrics struct{ schema.MemoryMetrics }
	// Metrics reexport with auth policy added
	Metrics struct{ schema.Metrics }
	// MissDetail reexport with auth policy added
	MissDetail struct{ schema.MissDetail }
	// NetworkMetrics reexport with auth policy added
	NetworkMetrics struct{ schema.NetworkMetrics }
	// PackageLoadMetrics reexport with auth policy added
	PackageLoadMetrics struct{ schema.PackageLoadMetrics }
	// PackageMetrics reexport with auth policy added
	PackageMetrics struct{ schema.PackageMetrics }
	// RunnerCount reexport with auth policy added
	RunnerCount struct{ schema.RunnerCount }
	// SourceControl reexport with auth policy added
	SourceControl struct{ schema.SourceControl }
	// SystemNetworkStats reexport with auth policy added
	SystemNetworkStats struct{ schema.SystemNetworkStats }
	// Target reexport with auth policy added
	Target struct{ schema.Target }
	// TargetKindMapping reexport with auth policy added
	TargetKindMapping struct{ schema.TargetKindMapping }
	// TargetMetrics reexport with auth policy added
	TargetMetrics struct{ schema.TargetMetrics }
	// TestResult reexport with auth policy added
	TestResult struct{ schema.TestResult }
	// TestSummary reexport with auth policy added
	TestSummary struct{ schema.TestSummary }
	// TimingMetrics reexport with auth policy added
	TimingMetrics struct{ schema.TimingMetrics }
)
