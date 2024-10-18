package summary

import (
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
)

// Step Label and user Key constants.
const (
	// stepLabelKey is used in buildMetadata events to provide a human-readable label for build steps.
	stepLabelKey = "BUILD_STEP_LABEL"
	userEmailKey = "user_email"
	userLdapKey  = "user_ldap"
)

// Exit Code constatn.
const (
	ExitCodeSuccess     = 0
	ExitCodeInterrupted = 8
)

// MissReason enum.
type MissReason int32

// MissReason enum.
const (
	MissReasonDIFFERENTACTIONKEY MissReason = iota + 1
	MissReasonDIFFERENTDEPS
	MissReasonDIFFERENTENVIRONMENT
	MissReasonDIFFERENTFILES
	MissReasonCORRUPTEDCACHEENTRY
	MissReasonNOTCACHED
	MissReasonUNCONDITIONALEXECUTION
)

// EnumIndex helper method.
func (r MissReason) EnumIndex() int32 {
	return int32(r)
}

// String Enum helper method.
func (r MissReason) String() string {
	return [...]string{
		"UNKNOWN",
		"DIFFERENT_ACTION_KEY",
		"DIFFERENT_DEPS",
		"DIFFERENT_ENVIRONMENT",
		"DIFFERENT_FILES",
		"CORRUPTED_CACHE_ENTRY",
		"NOT_CACHED",
		"UNCONDITIONAL_EXECUTION",
	}[r]
}

// TestStatus ENUM.
type TestStatus int32

// TestStatus enum.
const (
	TestStatuNOSTATUS TestStatus = iota + 1
	TestStatuPASSED
	TestStatuFLAKY
	TestStatuTIMEOUT
	TestStatuFAILED
	TestStatuINCOMPLETE
	TestStatuREMOTEFAILURE
	TestStatuFAILEDTOBUILD
	TestStatuTOOLHALTEDBEFORETESTING
)

// EnumIndex helper method.
func (r TestStatus) EnumIndex() int32 {
	return int32(r)
}

// String Enum helper method.
func (r TestStatus) String() string {
	return [...]string{
		"NO_STATUS",
		"PASSED",
		"FLAKY",
		"TIMEOUT",
		"FAILED",
		"INCOMPLETE",
		"REMOTE_FAILURE",
		"FAILED_TO_BUILD",
		"TOOL_HALTED_BEFORE_TESTING",
	}[r]
}

// TestSize Enum.
type TestSize int32

// TestSize enum.
const (
	UNKNOWN  TestSize = iota + 1 //nolint
	SMALL                        //nolint
	MEDIUM                       //nolint
	LARGE                        //nolint
	ENORMOUS                     //nolint
)

// EnumIndex Enum helper method.
func (r TestSize) EnumIndex() int32 {
	return int32(r)
}

// Enum helper method.
func (r TestSize) String() string {
	return [...]string{
		"UNKNOWN",
		"SMALL",
		"MEDIUM",
		"LARGE",
		"ENORMOUS",
	}[r]
}

// AbortReason enum.
type AbortReason int32

// AbortReason struct set the order explicitly because of ordering mismatches!
const (
	AbortedUNKNOWN                  AbortReason = 0
	AbortedUSERINTERRUPTED          AbortReason = 1
	AbortedTIMEOUT                  AbortReason = 2
	AbortedREMOTEENVIRONMENTFAILURE AbortReason = 3
	AbortedINTERNAL                 AbortReason = 4
	AbortedLOADINGFAILURE           AbortReason = 5
	AbortedANALYSISFAILURE          AbortReason = 6
	AbortedSKIPPED                  AbortReason = 7
	AbortedNOANALYZE                AbortReason = 8
	AbortedNOBUILD                  AbortReason = 9
	AbortedINCOMPLETE               AbortReason = 10
	AbortedOUTOFMEMORY              AbortReason = 11
)

// EnumIndex helper method.
func (r AbortReason) EnumIndex() int32 {
	return int32(r)
}

// String Enum helper method.
func (r AbortReason) String() string {
	return [...]string{
		"UNKNOWN",
		"USER_INTERRUPTED",
		"TIME_OUT",
		"REMOTE_ENVIRONMENT_FAILURE",
		"INTERNAL",
		"LOADING_FAILURE",
		"ANALYSIS_FAILURE",
		"SKIPPED",
		"NO_ANALYZE",
		"NO_BUILD",
		"INCOMPLETE",
		"OUT_OF_MEMORY",
	}[r]
}

// Summary The Invocation Summary object holds details about an invocation.
type Summary struct {
	*InvocationSummary
	Problems             []detectors.Problem
	RelatedFiles         map[string]string
	EventFileURL         string
	BEPCompleted         bool
	StartedAt            time.Time
	InvocationID         string
	StepLabel            string
	EndedAt              *time.Time
	ChangeNumber         int
	PatchsetNumber       int
	BuildURL             string
	BuildUUID            uuid.UUID
	UserLDAP             string
	UserEmail            string
	BuildLogs            strings.Builder
	Metrics              Metrics
	Tests                map[string]TestsCollection
	Targets              map[string]TargetPair
	NumFetches           int64
	CPU                  string
	PlatformName         string
	ProfileName          string
	ConfigrationMnemonic string
	SkipTargetData       bool
	EnrichTargetData     bool
}

// Metrics holds Build metrics details
// This aligngs with data found in the ent schema
// https://github.com/bazelbuild/bazel/blob/master/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream.proto#L900
type Metrics struct {
	ActionSummary           ActionSummary
	MemoryMetrics           MemoryMetrics
	TargetMetrics           TargetMetrics
	PackageMetrics          PackageMetrics
	TimingMetrics           TimingMetrics
	CumulativeMetrics       CumulativeMetrics
	ArtifactMetrics         ArtifactMetrics
	BuildGraphMetrics       BuildGraphMetrics
	NetworkMetrics          NetworkMetrics
	DynamicExecutionMetrics DynamicExecutionMetrics
}

// InvocationSummary struct.
type InvocationSummary struct {
	EnvVars          map[string]string
	ExitCode         *ExitCode
	BazelVersion     string
	BazelCommandLine BazelCommandLine
}

// ExitCode An Exit Code.
type ExitCode struct {
	Code int
	Name string
}

// BazelCommandLine struct.
type BazelCommandLine struct {
	Executable string
	Command    string
	Residual   string
	Options    []string
}

// Blob holds information about a blob in the CAS. Should be easily converted to/from the one in the
// cas package. Copied into here so this package does not have *any* dependencies except standard
// libraries.
type Blob struct {
	BlobURI  url.URL
	Size     int
	Contents string
	Name     string
}

// ActionSummary struct.
type ActionSummary struct {
	ActionsCreated                    int64
	ActionsCreatedNotIncludingAspects int64
	ActionsExecuted                   int64
	ActionData                        []ActionData

	RemoteCacheHits       int64
	RunnerCount           []RunnerCount
	ActionCacheStatistics ActionCacheStatistics
}

// ActionData struct
type ActionData struct {
	Mnemonic        string
	ActionsExecuted int64
	FirstStartedMs  int64
	LastEndedMs     int64
	SystemTime      int64
	UserTime        int64
}

// RunnerCount struct
type RunnerCount struct {
	Name     string
	Count    int32
	ExecKind string
}

// GarbageMetrics struct
type GarbageMetrics struct {
	Type             string
	GarbageCollected int64
}

// MemoryMetrics struct
type MemoryMetrics struct {
	UsedHeapSizePostBuild          int64
	PeakPostGcHeapSize             int64
	PeakPostGcTenuredSpaceHeapSize int64
	GarbageMetrics                 []GarbageMetrics
}

// TargetMetrics struct
type TargetMetrics struct {
	TargetsLoaded                        int64
	TargetsConfigured                    int64
	TargetsConfiguredNotIncludingAspects int64
}

// PackageMetrics struct
type PackageMetrics struct {
	PackagesLoaded     int64
	PackageLoadMetrics []PackageLoadMetrics
}

// TimingMetrics struct
type TimingMetrics struct {
	CPUTimeInMs            int64
	WallTimeInMs           int64
	AnalysisPhaseTimeInMs  int64
	ExecutionPhaseTimeInMs int64
}

// CumulativeMetrics struct
type CumulativeMetrics struct {
	NumAnalyses int32
	NumBuilds   int32
}

// ArtifactMetrics struct
type ArtifactMetrics struct {
	SourceArtifactsRead            FilesMetric
	OutputArtifactsSeen            FilesMetric
	OutputArtifactsFromActionCache FilesMetric
	TopLevelArtifacts              FilesMetric
}

// FilesMetric struct
type FilesMetric struct {
	SizeInBytes int64
	Count       int32
}

// SystemNetworkStats struct
type SystemNetworkStats struct {
	BytesSent             uint64
	BytesRecv             uint64
	PacketsSent           uint64
	PacketsRecv           uint64
	PeakBytesSentPerSec   uint64
	PeakBytesRecvPerSec   uint64
	PeakPacketsSentPerSec uint64
	PeakPacketsRecvPerSec uint64
}

// NetworkMetrics struct
type NetworkMetrics struct {
	SystemNetworkStats *SystemNetworkStats
}

// ActionCacheStatistics struct
type ActionCacheStatistics struct {
	SizeInBytes  uint64
	SaveTimeInMs uint64
	LoadTimeInMs int64
	Hits         int32
	Misses       int32
	MissDetails  []MissDetail
}

// MissDetail struct
type MissDetail struct {
	Reason MissReason
	Count  int32
}

// PackageLoadMetrics struct
type PackageLoadMetrics struct {
	Name               string
	LoadDuration       int64
	NumTargets         uint64
	ComputationSteps   uint64
	NumTransitiveLoads uint64
	PackageOverhead    uint64
}

// DynamicExecutionMetrics struct
type DynamicExecutionMetrics struct {
	RaceStatistics []RaceStatistics
}

// RaceStatistics struct
type RaceStatistics struct {
	Mnemonic     string
	LocalRunner  string
	RemoteRunner string
	LocalWins    int64
	RemoteWins   int64
}

// EvaluationStat struct
type EvaluationStat struct {
	SkyfunctionName string
	Count           int64
}

// BuildGraphMetrics struct
type BuildGraphMetrics struct {
	ActionLookupValueCount                    int32
	ActionLookupValueCountNotIncludingAspects int32
	ActionCount                               int32
	InputFileConfiguredTargetCount            int32
	OutputFileConfiguredTargetCount           int32
	OtherConfiguredTargetCount                int32
	OutputArtifactCount                       int32
	PostInvocationSkyframeNodeCount           int32
	DirtiedValues                             []EvaluationStat
	ChangedValues                             []EvaluationStat
	BuiltValues                               []EvaluationStat
	CleanedValues                             []EvaluationStat
	EvaluatedValues                           []EvaluationStat
}

// ExecutionInfo struct
type ExecutionInfo struct {
	Strategy        string
	CachedRemotely  bool
	ExitCode        int32
	Hostname        string
	TimingBreakdown TimingBreakdown
	ResourceUsage   []ResourceUsage
}

// TestResult struct
type TestResult struct {
	Status              TestStatus
	StatusDetails       string
	Label               string
	Warning             []string
	CachedLocally       bool
	TestAttemptDuration int64
	TestAttemptStart    string // timestamp
	TestActionOutput    []TestFile
	ExecutionInfo       ExecutionInfo
}

// TimingBreakdown struct
type TimingBreakdown struct {
	Name  string
	Time  string
	Child []TimingChild
}

// TimingChild struct
type TimingChild struct {
	Name string
	Time string
}

// ResourceUsage struct
type ResourceUsage struct {
	Name  string
	Value string
}

// TestFile struct
type TestFile struct {
	Digest string
	File   string
	Length int64
	Name   string
	Prefix []string
}

// TestSummary struct
type TestSummary struct {
	Label            string
	Status           TestStatus
	TotalRunCount    int32
	RunCount         int32
	AttemptCount     int32
	ShardCount       int32
	TotalNumCached   int32
	FirstStartTime   int64
	LastStopTime     int64
	TotalRunDuration int64
	Passed           []TestFile
	Failed           []TestFile
}

// TargetConfigured struct
type TargetConfigured struct {
	Tag        []string
	TargetKind string
	TestSize   TestSize

	// adding this to track time for a target
	// not ideal, TODO: can we somehow get a more accurate measure for this data
	StartTimeInMs int64
}

// TargetComplete struct
type TargetComplete struct {
	Success            bool
	TargetKind         string
	TestSize           TestSize
	OutputGroup        OutputGroup
	ImportantOutput    []TestFile
	DirectoryOutput    []TestFile
	Tag                []string
	TestTimeoutSeconds int64
	TestTimeout        int64
	// adding this to track time
	// not ideal, TODO: can we somehow get a more accurate measure for this data
	EndTimeInMs int64

	// TODO: Implement Faillure detail
	// FailureDetail FailureDetail
}

// OutputGroup struct
type OutputGroup struct {
	Name        string
	Incomplete  bool
	InlineFiles []TestFile
	FileSets    NamedSetOfFiles
}

// NamedSetOfFiles struct
type NamedSetOfFiles struct {
	Files    []TestFile
	FileSets *NamedSetOfFiles
}

// TestsCollection struct summary object for a test.
type TestsCollection struct {
	TestSummary    TestSummary
	TestResults    []TestResult
	OverallStatus  TestStatus
	Strategy       string
	CachedLocally  bool
	CachedRemotely bool
	DurationMs     int64
	FirstSeen      time.Time
}

// TargetPair struct summary object for a target.
type TargetPair struct {
	Configuration TargetConfigured
	Completion    TargetComplete
	DurationInMs  int64
	Success       bool
	TargetKind    string
	TestSize      TestSize
	AbortReason   AbortReason
}
