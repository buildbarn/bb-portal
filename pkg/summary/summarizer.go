package summary

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"

	"github.com/buildbarn/bb-portal/pkg/events"
	"github.com/buildbarn/bb-portal/pkg/summary/detectors"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bescore"
)

// Summarizer struct.
type Summarizer struct {
	summary         *Summary
	problemDetector detectors.ProblemDetector
}

// Summarize function.
func Summarize(ctx context.Context, eventFileURL string) (*Summary, error) {
	reader, err := os.Open(eventFileURL)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", eventFileURL, err)
	}
	defer reader.Close()

	problemDetector := detectors.NewProblemDetector()
	summarizer := newSummarizer(eventFileURL, problemDetector)
	it := events.NewBuildEventIterator(ctx, reader)
	return summarizer.summarize(it)
}

// NewSummarizer constructor
func NewSummarizer() *Summarizer {
	return newSummarizer("", detectors.NewProblemDetector())
}

// newSummarizer
func newSummarizer(eventFileURL string, problemDetector detectors.ProblemDetector) *Summarizer {
	return &Summarizer{
		summary: &Summary{
			InvocationSummary: &InvocationSummary{},
			EventFileURL:      eventFileURL,
			RelatedFiles: map[string]string{
				filepath.Base(eventFileURL): eventFileURL,
			},
		},
		problemDetector: problemDetector,
	}
}

// summarize
func (s Summarizer) summarize(it *events.BuildEventIterator) (*Summary, error) {
	for {
		buildEvent, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get build event: %w", err)
		}

		err = s.ProcessEvent(buildEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to process event (with id: %s): %w", buildEvent.Id.String(), err)
		}
	}

	return s.FinishProcessing()
}

// FinishProcessing function
func (s Summarizer) FinishProcessing() (*Summary, error) {
	// If problems are ignored for the exit code, return immediately.
	slog.Debug("processing", "err", "none")

	if s.summary.EndedAt == nil {
		now := time.Now()
		s.summary.EndedAt = &now
	}

	if s.summary.ExitCode == nil {
		s.summary.ExitCode = &ExitCode{
			Code: 1000,
			Name: "UNKNOWN",
		}
	}

	if !shouldIgnoreProblems(s.summary.ExitCode) {
		slog.Debug("problems found", "err", "none")
		// Add any detected test problems.
		problems, problemsErr := s.problemDetector.Problems()
		if problemsErr != nil {
			return nil, problemsErr
		}
		s.summary.Problems = append(s.summary.Problems, problems...)
	}
	slog.Debug("returning from FinishProcessing")

	return s.summary, nil
}

// ProcessEvent function
func (s Summarizer) ProcessEvent(buildEvent *events.BuildEvent) error {
	// Let problem detector process every event.
	s.problemDetector.ProcessBEPEvent(buildEvent)

	switch buildEvent.GetId().GetId().(type) {
	case *bes.BuildEventId_Started:
		s.handleStarted(buildEvent.GetStarted())

	case *bes.BuildEventId_BuildMetadata:
		s.handleBuildMetadata(buildEvent.GetBuildMetadata())

	case *bes.BuildEventId_BuildFinished:
		s.handleBuildFinished(buildEvent.GetFinished())

	case *bes.BuildEventId_BuildMetrics:
		s.handleBuildMetrics(buildEvent.GetBuildMetrics())

	case *bes.BuildEventId_StructuredCommandLine:
		err := s.handleStructuredCommandLine(buildEvent.GetStructuredCommandLine())
		if err != nil {
			return err
		}
	case *bes.BuildEventId_Configuration:
		s.handleBuildConfiguration(buildEvent.GetConfiguration())

	case *bes.BuildEventId_TargetConfigured:
		s.handleTargetConfigured(buildEvent.GetConfigured(), buildEvent.GetTargetConfiguredLabel(), time.Now())

	case *bes.BuildEventId_TargetCompleted:
		s.handleTargetCompleted(buildEvent.GetCompleted(), buildEvent.GetTargetCompletedLabel(), buildEvent.GetAborted(), time.Now())

	case *bes.BuildEventId_Fetch:
		s.handleFetch(buildEvent.GetFetch())

	case *bes.BuildEventId_TestResult:
		s.handleTestResult(buildEvent.GetTestResult(), buildEvent.GetId().GetTestResult().Label)

	case *bes.BuildEventId_TestSummary:
		s.handleTestSummary(buildEvent.GetTestSummary(), buildEvent.GetId().GetTestSummary().Label)

	case *bes.BuildEventId_OptionsParsed:
		s.handleOptionsParsed(buildEvent.GetOptionsParsed())

	case *bes.BuildEventId_BuildToolLogs:
		err := s.handleBuildToolLogs(buildEvent.GetBuildToolLogs())
		if err != nil {
			return err
		}
	}

	s.summary.BEPCompleted = buildEvent.GetLastMessage()
	return nil
}

// handleStarted
func (s Summarizer) handleStarted(started *bes.BuildStarted) {
	var startedAt time.Time
	if started.GetStartTime() != nil {
		startedAt = started.GetStartTime().AsTime()
	} else {
		//nolint:staticcheck // Keep backwards compatibility until the field is removed.
		startedAt = time.UnixMilli(started.GetStartTimeMillis())
	}
	s.summary.StartedAt = startedAt
	s.summary.InvocationID = started.GetUuid()
	s.summary.BazelVersion = started.GetBuildToolVersion()
}

// handleFetch
func (s Summarizer) handleFetch(fetch *bes.Fetch) {
	if fetch.Success {
		s.summary.NumFetches++
	}
}

// handleBuildConfiguration
func (s Summarizer) handleBuildConfiguration(configuration *bes.Configuration) {
	s.summary.CPU = configuration.Cpu
	s.summary.PlatformName = configuration.PlatformName
	s.summary.ConfigrationMnemonic = configuration.Mnemonic
}

// handleTargetConfigured
func (s Summarizer) handleTargetConfigured(target *bes.TargetConfigured, label string, timestamp time.Time) {
	if len(label) == 0 {
		slog.Warn("missing a target label for target configured event!")
		return
	}
	if target == nil {
		slog.Warn(fmt.Sprintf("missing target for label %s on targetConfigured", label))
		return
	}

	// if this is the first target we've seen, initialize the targets collection
	if s.summary.Targets == nil {
		s.summary.Targets = make(map[string]TargetPair)
	}

	// create a target pair and at it to the targets collection
	s.summary.Targets[label] = TargetPair{
		Configuration: TargetConfigured{
			StartTimeInMs: timestamp.UnixMilli(),
			TargetKind:    target.TargetKind,
			TestSize:      TestSize(target.TestSize),
			Tag:           target.Tag,
		},
		Success:    false, // set it to false, change it when we get a complete
		TargetKind: target.TargetKind,
		TestSize:   TestSize(target.TestSize),
	}
}

// handleTargetCompleted
func (s Summarizer) handleTargetCompleted(target *bes.TargetComplete, label string, aborted *bes.Aborted, timestamp time.Time) {
	if len(label) == 0 {
		slog.Error("label is empty for a target completed event")
		return
	}

	if s.summary.Targets == nil {
		slog.Warn(fmt.Sprintf("target completed event received before any target configured messages for label %s,", label))
		return
	}

	var targetPair TargetPair
	targetPair, ok := s.summary.Targets[label]

	if !ok {
		slog.Warn(fmt.Sprintf("target completed event received for label %s before target configured message received", label))
		return
	}

	var targetCompletion TargetComplete

	if target != nil {
		targetCompletion = TargetComplete{
			Success:     target.Success,
			Tag:         target.Tag,
			EndTimeInMs: timestamp.UnixMilli(),
		}
		if target.TestTimeout != nil {
			targetCompletion.TestTimeoutSeconds = target.TestTimeout.Seconds
			targetCompletion.TestTimeout = target.TestTimeout.Seconds
		}
	} else { // this event was aborted
		targetCompletion = TargetComplete{
			Success:     false,
			EndTimeInMs: timestamp.UnixMilli(),
		}
	}

	targetPair.Completion = targetCompletion
	targetPair.DurationInMs = targetPair.Completion.EndTimeInMs - targetPair.Configuration.StartTimeInMs
	targetPair.Success = targetCompletion.Success

	if aborted != nil {
		targetPair.AbortReason = AbortReason(aborted.Reason)
	}

	s.summary.Targets[label] = targetPair
}

// handleTestResult
func (s Summarizer) handleTestResult(testResult *bes.TestResult, label string) {
	if len(label) == 0 {
		slog.Warn("missing label on TestResult event", "err", nil)
		return
	}
	if testResult == nil {
		slog.Warn("Missing Test Result for label %s", label, nil)
		return
	}
	var testResults []TestResult
	if s.summary.Tests == nil {
		s.summary.Tests = make(map[string]TestsCollection)
	}
	testcollection, ok := s.summary.Tests[label]
	if ok {
		testResults = testcollection.TestResults
	} else { // initailize it if we've never seen this label before
		testcollection = TestsCollection{
			TestSummary:    TestSummary{},
			TestResults:    []TestResult{},
			CachedLocally:  true,
			CachedRemotely: true,
			Strategy:       "INITIALIZED",
			FirstSeen:      time.Now(), // this is primarly used for sorting
		}
		testResults = make([]TestResult, 0)
	}
	tr := TestResult{
		Status:              TestStatus(testResult.Status),
		StatusDetails:       testResult.StatusDetails,
		Label:               label,
		Warning:             testResult.Warning,
		CachedLocally:       testResult.CachedLocally,
		TestAttemptStart:    testResult.TestAttemptStart.AsTime().String(),
		TestAttemptDuration: testResult.TestAttemptDuration.AsDuration().Milliseconds(),
		ExecutionInfo:       processExecutionInfo(testResult),
		TestActionOutput:    make([]TestFile, 0),
	}
	for _, ao := range testResult.TestActionOutput {
		actionOutput := TestFile{
			Digest: ao.Digest,
			File:   ao.GetUri(),
			Length: ao.Length,
			Name:   ao.Name,
			Prefix: ao.PathPrefix,
		}
		tr.TestActionOutput = append(tr.TestActionOutput, actionOutput)
	}
	testResults = append(testResults, tr)
	testcollection.TestResults = testResults
	if testResult.Status != bes.TestStatus_NO_STATUS {
		if !testResult.CachedLocally {
			testcollection.CachedLocally = false
		}
		if !tr.ExecutionInfo.CachedRemotely {
			testcollection.CachedRemotely = false
		}
		if (testcollection.Strategy) == "INITIALIZED" {
			testcollection.Strategy = tr.ExecutionInfo.Strategy
		} else {
			if testcollection.Strategy != tr.ExecutionInfo.Strategy {
				testcollection.Strategy = "indeterminate"
			}
		}
	}

	s.summary.Tests[label] = testcollection
}

// processExecutionInfo
func processExecutionInfo(testResult *bes.TestResult) ExecutionInfo {
	var result ExecutionInfo
	var timingBreakdown TimingBreakdown
	var children []TimingChild
	if testResult.ExecutionInfo != nil {
		if testResult.ExecutionInfo.TimingBreakdown != nil {
			for _, c := range testResult.ExecutionInfo.TimingBreakdown.Child {
				child := TimingChild{
					Name: c.Name,
					Time: c.Time.AsDuration().String(),
				}
				children = append(children, child)
			}

			timingBreakdown.Name = testResult.ExecutionInfo.TimingBreakdown.Name
			timingBreakdown.Time = testResult.ExecutionInfo.TimingBreakdown.Time.String()
			timingBreakdown.Child = children
		}

		result.Strategy = testResult.ExecutionInfo.Strategy
		result.CachedRemotely = testResult.ExecutionInfo.CachedRemotely
		result.ExitCode = testResult.ExecutionInfo.ExitCode
		result.Hostname = testResult.ExecutionInfo.Hostname
		result.TimingBreakdown = timingBreakdown
	}
	return result
}

// handleTestSummary
func (s Summarizer) handleTestSummary(testSummary *bes.TestSummary, label string) {
	if len(label) == 0 {
		slog.Error("missing label on handleTestSummary event")
		return
	}

	if testSummary == nil {
		slog.Warn("missing test summary object for handleTestSummary event for label %s", label, nil)
		return
	}

	testCollection, ok := s.summary.Tests[label]

	if !ok {
		slog.Warn("received a test summary event but never first saw a test result for label %s", label, nil)
		return
	}

	tSummary := testCollection.TestSummary

	tSummary.AttemptCount = testSummary.AttemptCount
	tSummary.FirstStartTime = testSummary.FirstStartTime.AsTime().Unix()
	tSummary.Label = label
	tSummary.LastStopTime = testSummary.FirstStartTime.AsTime().Unix()
	tSummary.RunCount = testSummary.RunCount
	tSummary.ShardCount = testSummary.ShardCount
	tSummary.Status = TestStatus(testSummary.OverallStatus)
	tSummary.TotalNumCached = testSummary.TotalNumCached
	tSummary.TotalRunCount = testSummary.TotalRunCount
	tSummary.TotalRunDuration = testSummary.TotalRunDuration.AsDuration().Microseconds()

	testCollection.TestSummary = tSummary
	testCollection.OverallStatus = tSummary.Status
	testCollection.DurationMs = tSummary.TotalRunDuration
	s.summary.Tests[label] = testCollection
}

// handleBuildMetadata
func (s Summarizer) handleBuildMetadata(metadataProto *bes.BuildMetadata) {
	metadataMap := metadataProto.GetMetadata()
	// extract metadata
	if metadataMap == nil {
		return
	}
	if stepLabel, ok := metadataMap[stepLabelKey]; ok {
		s.summary.StepLabel = stepLabel
	}
	if userEmail, ok := metadataMap[userEmailKey]; ok {
		s.summary.UserEmail = userEmail
	}
	if userLdap, ok := metadataMap[userLdapKey]; ok {
		s.summary.UserLDAP = userLdap
	}
	if isCiWorkerVal, ok := metadataMap[isCiWorkerKey]; ok {
		if isCiWorkerVal == "TRUE" {
			s.summary.IsCiWorker = true
		}
	}
	if hostnameVal, ok := metadataMap[hostnameKey]; ok {
		s.summary.Hostname = hostnameVal
	}
}

// handleBuildMetrics
func (s Summarizer) handleBuildMetrics(metrics *bes.BuildMetrics) {
	actionSummary := readActionSummary(metrics.ActionSummary)
	memoryMetrics := readMemoryMetrics(metrics.MemoryMetrics)
	targetMetrics := readTargetMetrics(metrics.TargetMetrics)
	packageMetrics := readPackageMetrics(metrics.PackageMetrics)
	timingMetrics := readTimingMetrics(metrics.TimingMetrics)
	artifactMetrics := readArtifactMetrics(metrics.ArtifactMetrics)
	cumulativeMetrics := readCumulativeMetrics(metrics.CumulativeMetrics)
	networkMetrics := readNetworkMetrics(metrics.NetworkMetrics)
	buildGraphMetrics := readBuildGraphMetrics(metrics.BuildGraphMetrics)

	summaryMetrics := Metrics{
		ActionSummary:     actionSummary,
		MemoryMetrics:     memoryMetrics,
		TargetMetrics:     targetMetrics,
		PackageMetrics:    packageMetrics,
		TimingMetrics:     timingMetrics,
		ArtifactMetrics:   artifactMetrics,
		CumulativeMetrics: cumulativeMetrics,
		NetworkMetrics:    networkMetrics,
		BuildGraphMetrics: buildGraphMetrics,
		// DynamicExecutionMetrics: dynamicMetrics,
	}

	s.summary.Metrics = summaryMetrics
}

// readBuildGraphMetrics
func readBuildGraphMetrics(buildGraphMetricsData *bes.BuildMetrics_BuildGraphMetrics) BuildGraphMetrics {
	// TODO: these values are not on the proto currently.  once they are, update this code to pull them out
	// var dirtiedValues = make([]EvaluationStat, 0)
	// var changedValues = make([]EvaluationStat, 0)
	// var builtValues = make([]EvaluationStat, 0)
	// var cleanedValues = make([]EvaluationStat, 0)
	// var evaluatedValues = make([]EvaluationStat, 0)

	// DirtiedValues:                   dirtiedValues,
	// ChangedValues:                   changedValues,
	// BuiltValues:                     builtValues,
	// CleanedValues:                   cleanedValues,
	// EvaluatedValues:                 evaluatedValues,
	buildGraphMetrics := BuildGraphMetrics{
		ActionLookupValueCount:                    buildGraphMetricsData.ActionLookupValueCount,
		ActionLookupValueCountNotIncludingAspects: buildGraphMetricsData.ActionLookupValueCountNotIncludingAspects,
		ActionCount:                     buildGraphMetricsData.ActionCount,
		InputFileConfiguredTargetCount:  buildGraphMetricsData.InputFileConfiguredTargetCount,
		OutputFileConfiguredTargetCount: buildGraphMetricsData.OutputFileConfiguredTargetCount,
		OtherConfiguredTargetCount:      buildGraphMetricsData.OtherConfiguredTargetCount,
		OutputArtifactCount:             buildGraphMetricsData.OutputArtifactCount,
		PostInvocationSkyframeNodeCount: buildGraphMetricsData.PostInvocationSkyframeNodeCount,
	}
	return buildGraphMetrics
}

// readNetworkMetrics
func readNetworkMetrics(networkMetricsData *bes.BuildMetrics_NetworkMetrics) NetworkMetrics {
	if networkMetricsData == nil {
		return NetworkMetrics{}
	}
	systemNetworkStats := readSystemNetworkStats(networkMetricsData.SystemNetworkStats)

	networkMetrics := NetworkMetrics{
		SystemNetworkStats: &systemNetworkStats,
	}
	return networkMetrics
}

// readSystemNetworkStats
func readSystemNetworkStats(systemNetworkStatsData *bes.BuildMetrics_NetworkMetrics_SystemNetworkStats) SystemNetworkStats {
	var systemNetworkStats SystemNetworkStats
	if systemNetworkStatsData != nil {
		systemNetworkStats = SystemNetworkStats{
			BytesSent:             systemNetworkStatsData.BytesSent,
			BytesRecv:             systemNetworkStatsData.BytesRecv,
			PacketsSent:           systemNetworkStatsData.PacketsSent,
			PacketsRecv:           systemNetworkStatsData.PacketsRecv,
			PeakBytesSentPerSec:   systemNetworkStatsData.PeakBytesSentPerSec,
			PeakBytesRecvPerSec:   systemNetworkStatsData.PeakBytesRecvPerSec,
			PeakPacketsSentPerSec: systemNetworkStatsData.PeakPacketsSentPerSec,
			PeakPacketsRecvPerSec: systemNetworkStatsData.PeakPacketsRecvPerSec,
		}
	}
	return systemNetworkStats
}

// readCumulativeMetrics
func readCumulativeMetrics(cumulativeMetricsData *bes.BuildMetrics_CumulativeMetrics) CumulativeMetrics {
	if cumulativeMetricsData == nil {
		return CumulativeMetrics{}
	}

	cumulativeMetrics := CumulativeMetrics{
		NumAnalyses: cumulativeMetricsData.NumAnalyses,
		NumBuilds:   cumulativeMetricsData.NumBuilds,
	}
	return cumulativeMetrics
}

// readArtifactMetrics
func readArtifactMetrics(artifactMetricsData *bes.BuildMetrics_ArtifactMetrics) ArtifactMetrics {
	if artifactMetricsData == nil {
		return ArtifactMetrics{}
	}

	sourceArtifactsRead := FilesMetric{}
	if artifactMetricsData.SourceArtifactsRead != nil {
		sourceArtifactsRead.SizeInBytes = artifactMetricsData.SourceArtifactsRead.SizeInBytes
		sourceArtifactsRead.Count = artifactMetricsData.SourceArtifactsRead.Count
	}

	outputArtifactsSeen := FilesMetric{}
	if artifactMetricsData.OutputArtifactsSeen != nil {
		outputArtifactsSeen.SizeInBytes = artifactMetricsData.OutputArtifactsSeen.SizeInBytes
		outputArtifactsSeen.Count = artifactMetricsData.OutputArtifactsSeen.Count
	}

	outputArtifactsFromActionCache := FilesMetric{}
	if artifactMetricsData.OutputArtifactsFromActionCache != nil {
		outputArtifactsFromActionCache.SizeInBytes = artifactMetricsData.OutputArtifactsFromActionCache.SizeInBytes
		outputArtifactsFromActionCache.Count = artifactMetricsData.OutputArtifactsFromActionCache.Count
	}

	topLevelArtifacts := FilesMetric{}
	if artifactMetricsData.TopLevelArtifacts != nil {
		topLevelArtifacts.SizeInBytes = artifactMetricsData.TopLevelArtifacts.SizeInBytes
		topLevelArtifacts.Count = artifactMetricsData.TopLevelArtifacts.Count
	}

	artifactMetrics := ArtifactMetrics{
		SourceArtifactsRead:            sourceArtifactsRead,
		OutputArtifactsSeen:            outputArtifactsSeen,
		OutputArtifactsFromActionCache: outputArtifactsFromActionCache,
		TopLevelArtifacts:              topLevelArtifacts,
	}
	return artifactMetrics
}

// readTimingMetrics
func readTimingMetrics(timingMetricsData *bes.BuildMetrics_TimingMetrics) TimingMetrics {
	if timingMetricsData == nil {
		return TimingMetrics{}
	}

	timingMetrics := TimingMetrics{
		CPUTimeInMs:            timingMetricsData.CpuTimeInMs,
		WallTimeInMs:           timingMetricsData.WallTimeInMs,
		ExecutionPhaseTimeInMs: timingMetricsData.ExecutionPhaseTimeInMs,
		AnalysisPhaseTimeInMs:  timingMetricsData.AnalysisPhaseTimeInMs,
	}
	return timingMetrics
}

// readPackageMetrics
func readPackageMetrics(packageMetricsData *bes.BuildMetrics_PackageMetrics) PackageMetrics {
	if packageMetricsData == nil {
		return PackageMetrics{}
	}

	packageLoadMetrics := readPackageLoadMetrics(packageMetricsData.PackageLoadMetrics)

	packageMetrics := PackageMetrics{
		PackagesLoaded:     packageMetricsData.PackagesLoaded,
		PackageLoadMetrics: packageLoadMetrics,
	}
	return packageMetrics
}

// readTargetMetrics
func readTargetMetrics(targetMetricsData *bes.BuildMetrics_TargetMetrics) TargetMetrics {
	if targetMetricsData == nil {
		return TargetMetrics{}
	}

	targetMetrics := TargetMetrics{
		TargetsConfigured:                    targetMetricsData.TargetsConfigured,
		TargetsConfiguredNotIncludingAspects: targetMetricsData.TargetsConfiguredNotIncludingAspects,
		TargetsLoaded:                        targetMetricsData.TargetsLoaded,
	}
	return targetMetrics
}

// readMemoryMetrics
func readMemoryMetrics(memoryMetricsData *bes.BuildMetrics_MemoryMetrics) MemoryMetrics {
	if memoryMetricsData == nil {
		return MemoryMetrics{}
	}

	garbageMetrics := readGarbageMetrics(memoryMetricsData.GarbageMetrics)

	memoryMetrics := MemoryMetrics{
		PeakPostGcHeapSize:             memoryMetricsData.PeakPostGcHeapSize,
		PeakPostGcTenuredSpaceHeapSize: memoryMetricsData.PeakPostGcTenuredSpaceHeapSize,
		UsedHeapSizePostBuild:          memoryMetricsData.UsedHeapSizePostBuild,
		GarbageMetrics:                 garbageMetrics,
	}
	return memoryMetrics
}

// readActionSummary
func readActionSummary(actionSummaryData *bes.BuildMetrics_ActionSummary) ActionSummary {
	if actionSummaryData == nil {
		return ActionSummary{}
	}

	actionCacheStatistics := readActionCacheStatistics(actionSummaryData.ActionCacheStatistics)

	runnerCounts := readRunnerCounts(actionSummaryData.RunnerCount)

	actionDatas := readActionDatas(actionSummaryData.ActionData)

	actionSummary := ActionSummary{
		ActionsCreated:                    actionSummaryData.ActionsCreated,
		ActionsExecuted:                   actionSummaryData.ActionsExecuted,
		ActionsCreatedNotIncludingAspects: actionSummaryData.ActionsCreatedNotIncludingAspects,
		ActionCacheStatistics:             actionCacheStatistics,
		RunnerCount:                       runnerCounts,
		ActionData:                        actionDatas,
	}
	return actionSummary
}

// readActionCacheStatistics
func readActionCacheStatistics(actionCacheStatisticsData *bescore.ActionCacheStatistics) ActionCacheStatistics {
	if actionCacheStatisticsData == nil {
		return ActionCacheStatistics{}
	}
	missDetails := readMissDetails(actionCacheStatisticsData.MissDetails)
	actionCacheStatistics := ActionCacheStatistics{
		SizeInBytes:  actionCacheStatisticsData.SizeInBytes,
		SaveTimeInMs: actionCacheStatisticsData.SaveTimeInMs,

		Hits:        actionCacheStatisticsData.Hits,
		Misses:      actionCacheStatisticsData.Misses,
		MissDetails: missDetails,
	}
	return actionCacheStatistics
}

// readPackageLoadMetrics
func readPackageLoadMetrics(packageLoadMetricsData []*bescore.PackageLoadMetrics) []PackageLoadMetrics {
	packageLoadMetrics := make([]PackageLoadMetrics, len(packageLoadMetricsData))

	for i, plm := range packageLoadMetricsData {
		packageLoadMetric := PackageLoadMetrics{
			Name:               *plm.Name,
			NumTargets:         *plm.NumTargets,
			LoadDuration:       plm.LoadDuration.AsDuration().Milliseconds(),
			ComputationSteps:   *plm.ComputationSteps,
			NumTransitiveLoads: *plm.NumTransitiveLoads,
			PackageOverhead:    *plm.PackageOverhead,
		}
		packageLoadMetrics[i] = packageLoadMetric
	}
	return packageLoadMetrics
}

// readGarbageMetrics
func readGarbageMetrics(garbageMetricsData []*bes.BuildMetrics_MemoryMetrics_GarbageMetrics) []GarbageMetrics {
	garbageMetrics := make([]GarbageMetrics, len(garbageMetricsData))

	for i, gm := range garbageMetricsData {
		garbageMetric := GarbageMetrics{
			Type:             gm.Type,
			GarbageCollected: gm.GarbageCollected,
		}
		garbageMetrics[i] = garbageMetric
	}
	return garbageMetrics
}

// readActionDatas
func readActionDatas(actionDataData []*bes.BuildMetrics_ActionSummary_ActionData) []ActionData {
	actionDatas := make([]ActionData, len(actionDataData))
	for i, ad := range actionDataData {
		actionData := ActionData{
			Mnemonic:        ad.Mnemonic,
			UserTime:        ad.UserTime.AsDuration().Milliseconds(),
			SystemTime:      ad.SystemTime.AsDuration().Milliseconds(),
			ActionsExecuted: ad.ActionsExecuted,
			FirstStartedMs:  ad.FirstStartedMs,
			LastEndedMs:     ad.LastEndedMs,
		}
		actionDatas[i] = actionData
	}
	return actionDatas
}

// readRunnerCounts
func readRunnerCounts(runnerCountsData []*bes.BuildMetrics_ActionSummary_RunnerCount) []RunnerCount {
	runnerCounts := make([]RunnerCount, len(runnerCountsData))
	for i, rc := range runnerCountsData {
		runnerCount := RunnerCount{
			ExecKind: rc.ExecKind,
			Count:    rc.Count,
			Name:     rc.Name,
		}
		runnerCounts[i] = runnerCount
	}
	return runnerCounts
}

// readMissDetails
func readMissDetails(missDetailsData []*bescore.ActionCacheStatistics_MissDetail) []MissDetail {
	missDetails := make([]MissDetail, len(missDetailsData))
	for i, md := range missDetailsData {
		missDetail := MissDetail{
			Count:  md.Count,
			Reason: MissReason(*md.Reason.Enum()),
		}
		missDetails[i] = missDetail
	}
	return missDetails
}

// handleBuildFinished
func (s Summarizer) handleBuildFinished(finished *bes.BuildFinished) {
	var endedAt time.Time
	if finished.GetFinishTime() != nil {
		endedAt = finished.GetFinishTime().AsTime()
	} else {
		//nolint:staticcheck // Keep backwards compatibility until the field is removed.
		endedAt = time.UnixMilli(finished.GetFinishTimeMillis())
	}
	s.summary.EndedAt = &endedAt
	s.summary.InvocationSummary.ExitCode = &ExitCode{
		Code: int(finished.GetExitCode().GetCode()),
		Name: finished.GetExitCode().GetName(),
	}
}

// handleStructuredCommandLine
func (s Summarizer) handleStructuredCommandLine(structuredCommandLine *bescore.CommandLine) error {
	if structuredCommandLine.GetCommandLineLabel() != "original" {
		return nil
	}

	s.updateSummaryFromStructuredCommandLine(structuredCommandLine)

	// Parse Gerrit change number if available.
	if changeNumberStr, ok := s.summary.InvocationSummary.EnvVars["GERRIT_CHANGE_NUMBER"]; ok && changeNumberStr != "" {
		changeNumber, err := envToI(s.summary.InvocationSummary.EnvVars, "GERRIT_CHANGE_NUMBER")
		if err != nil {
			return err
		}
		s.summary.ChangeNumber = changeNumber
	}

	// Parse Gerrit patchset number if available.
	if patchsetNumberStr, ok := s.summary.InvocationSummary.EnvVars["GERRIT_PATCHSET_NUMBER"]; ok && patchsetNumberStr != "" {
		patchsetNumber, err := envToI(s.summary.InvocationSummary.EnvVars, "GERRIT_PATCHSET_NUMBER")
		if err != nil {
			return err
		}
		s.summary.PatchsetNumber = patchsetNumber
	}

	// Decode commit message, so that client doesn't have to.
	commitMessage := s.summary.InvocationSummary.EnvVars["GERRIT_CHANGE_COMMIT_MESSAGE"]
	if commitMessage != "" {
		decodedCommitMessage, err := base64.StdEncoding.DecodeString(commitMessage)
		if err == nil {
			s.summary.InvocationSummary.EnvVars["GERRIT_CHANGE_COMMIT_MESSAGE"] = string(decodedCommitMessage)
		} else {
			slog.Debug("GERRIT_CHANGE_COMMIT_MESSAGE was not base64 encoded, assuming it is normal string")
		}
	}

	// Set build URL and UUID
	s.summary.BuildURL = s.summary.InvocationSummary.EnvVars["BUILD_URL"]
	s.summary.BuildUUID = uuid.NewSHA1(uuid.NameSpaceURL, []byte(s.summary.BuildURL))

	// Set Hostname
	if hostNameVal, ok := s.summary.EnvVars["RUNNER_NAME"]; ok {
		s.summary.Hostname = hostNameVal
	}
	if hostNameVal, ok := s.summary.EnvVars["HOSTNAME"]; ok {
		s.summary.Hostname = hostNameVal
	}
	if hostNameVal, ok := s.summary.EnvVars["BB_PORTAL_HOSTNAME"]; ok {
		s.summary.Hostname = hostNameVal
	}

	// Set CI Worker Role from environment variables (can also come from metadata)
	if isCiWorkerVal, ok := s.summary.EnvVars["BB_PORTAL_IS_CI_WORKER"]; ok {
		if isCiWorkerVal == "TRUE" {
			s.summary.IsCiWorker = true
		}
	}

	// github actions default env var
	if isCiWorkerVal, ok := s.summary.EnvVars["CI"]; ok {
		if isCiWorkerVal == "true" {
			s.summary.IsCiWorker = true
		}
	}

	// default step label to workfow + job
	if ghWfVal, ok := s.summary.EnvVars["GITHUB_WORKFLOW"]; ok {
		s.summary.StepLabel = ghWfVal
		if ghJobNameVal, ok := s.summary.EnvVars["GITHUB_JOB"]; ok {
			s.summary.StepLabel += "+" + ghJobNameVal
		}
	}

	// Set Step Label from environment variables
	if stepLabelVal, ok := s.summary.EnvVars["BB_PORTAL_STEP_LABEL"]; ok {
		s.summary.StepLabel = stepLabelVal
	}

	// Set SkipTargetData
	if skipTargetSaveEnvVarVal, ok := s.summary.EnvVars["BB_PORTAL_SKIP_SAVE_TARGETS"]; ok {
		if skipTargetSaveEnvVarVal == "TRUE" {
			s.summary.SkipTargetData = true
		}
	}

	// Set EnrichTargetData
	if enrichTargetDataVal, ok := s.summary.EnvVars["BB_PORTAL_ENRICH_TARGET_DATA"]; ok {
		if enrichTargetDataVal == "TRUE" {
			s.summary.EnrichTargetData = true
		}
	}

	// repo
	if ghRepo, ok := s.summary.EnvVars["GITHUB_REPOSITORY"]; ok {
		s.summary.SourceControlData.RepositoryURL = ghRepo
	}

	// head ref
	if ghHeadRef, ok := s.summary.EnvVars["GITHUB_HEAD_REF"]; ok {
		s.summary.SourceControlData.Branch = ghHeadRef
	}

	// commit sha
	if ghSha, ok := s.summary.EnvVars["GITHUB_SHA"]; ok {
		s.summary.SourceControlData.CommitSHA = ghSha
	}

	// actor
	if ghActor, ok := s.summary.EnvVars["GITHUB_ACTOR"]; ok {
		s.summary.SourceControlData.Actor = ghActor
		if s.summary.UserLDAP == "" {
			s.summary.UserLDAP = ghActor
		}
	}

	// refs
	if ghRefs, ok := s.summary.EnvVars["GITHUB_REF"]; ok {
		s.summary.SourceControlData.Refs = ghRefs
	}

	// run id
	if ghRunID, ok := s.summary.EnvVars["GITHUB_RUN_ID"]; ok {
		s.summary.SourceControlData.RunID = ghRunID
	}

	// workflow
	if ghWorkflow, ok := s.summary.EnvVars["GITHUB_WORKFLOW"]; ok {
		s.summary.SourceControlData.Workflow = ghWorkflow
	}

	// action
	if ghAction, ok := s.summary.EnvVars["GITHUB_ACTION"]; ok {
		s.summary.SourceControlData.Action = ghAction
	}

	// workspace
	if ghWorkspace, ok := s.summary.EnvVars["GITHUB_WORKSPACE"]; ok {
		s.summary.SourceControlData.Workspace = ghWorkspace
	}

	// event_name
	if ghEventName, ok := s.summary.EnvVars["GITHUB_EVENT_NAME"]; ok {
		s.summary.SourceControlData.EventName = ghEventName
	}

	// job
	if ghJob, ok := s.summary.EnvVars["GITHUB_JOB"]; ok {
		s.summary.SourceControlData.Job = ghJob
	}

	// runner arch
	if runnerArch, ok := s.summary.EnvVars["RUNNER_ARCH"]; ok {
		s.summary.SourceControlData.RunnerArch = runnerArch
	}

	// runner name
	if runnerName, ok := s.summary.EnvVars["RUNNER_NAME"]; ok {
		s.summary.SourceControlData.RunnerName = runnerName
	}

	// runner os
	if runnerOs, ok := s.summary.EnvVars["RUNNER_OS"]; ok {
		s.summary.SourceControlData.RunnerOs = runnerOs
	}

	return nil
}

// handleOptionsParsed
func (s Summarizer) handleOptionsParsed(optionsParsed *bes.OptionsParsed) {
	s.summary.InvocationSummary.BazelCommandLine.ExplicitCmdLine = optionsParsed.GetExplicitCmdLine()
	s.summary.InvocationSummary.BazelCommandLine.CmdLine = optionsParsed.GetCmdLine()
	s.summary.InvocationSummary.BazelCommandLine.ExplicitStartupOptions = optionsParsed.GetExplicitStartupOptions()
	s.summary.InvocationSummary.BazelCommandLine.StartUpOptions = optionsParsed.GetStartupOptions()
}

// handleBuildToolLogs
func (s Summarizer) handleBuildToolLogs(buildToolLogs *bes.BuildToolLogs) error {
	for _, logs := range buildToolLogs.GetLog() {
		uri := logs.GetUri()
		blobURI := detectors.BlobURI(uri)

		if s.summary.RelatedFiles == nil {
			s.summary.RelatedFiles = map[string]string{}
		}
		if logs.GetUri() != "" {
			s.summary.RelatedFiles[logs.GetName()] = string(blobURI)
		}
	}
	return nil
}

// updateSummaryFromStructuredCommandLine
func (s Summarizer) updateSummaryFromStructuredCommandLine(structuredCommandLine *bescore.CommandLine) {
	sections := structuredCommandLine.GetSections()
	for _, section := range sections {
		label := section.GetSectionLabel()
		if label == "command options" {
			s.summary.InvocationSummary.EnvVars = map[string]string{}
			parseEnvVarsFromSectionOptions(section, &s.summary.InvocationSummary.EnvVars)
			s.summary.ProfileName = parseProfileNameFromSectionOptions(section)
		} else if section.GetChunkList() != nil {
			sectionChunksStr := strings.Join(section.GetChunkList().GetChunk(), " ")
			switch label {
			case "executable":
				s.summary.InvocationSummary.BazelCommandLine.Executable = sectionChunksStr
			case "command":
				s.summary.InvocationSummary.BazelCommandLine.Command = sectionChunksStr
			case "residual":
				s.summary.InvocationSummary.BazelCommandLine.Residual = sectionChunksStr
			}
		}
	}
}

// shouldIgnoreProblems
func shouldIgnoreProblems(exitCode *ExitCode) bool {
	return exitCode != nil && (exitCode.Code == ExitCodeSuccess || exitCode.Code == ExitCodeInterrupted)
}

// envToI
func envToI(envVars map[string]string, name string) (int, error) {
	res, err := strconv.Atoi(envVars[name])
	if err != nil {
		slog.Error("failed to parse env var to int", "envKey", name, "envValue", envVars[name], "err", err)
		return 0, fmt.Errorf("failed to parse %s (value: %s) as an int: %w", name, envVars[name], err)
	}
	return res, nil
}

func parseEnvVarsFromSectionOptions(section *bescore.CommandLineSection, destMap *map[string]string) {
	if section.GetOptionList() == nil {
		return
	}
	options := section.GetOptionList().GetOption()
	for _, option := range options {
		if option.GetOptionName() != "client_env" {
			// Only looking for env vars from the client env
			continue
		}
		envPair := option.GetOptionValue()
		equalIndex := strings.Index(envPair, "=")
		if equalIndex <= 0 {
			// Skip anything missing an equals sign. The env vars come in the format key=value
			continue
		}
		envName := envPair[:equalIndex]
		envValue := envPair[equalIndex+1:]
		(*destMap)[envName] = envValue
	}
}

func parseProfileNameFromSectionOptions(section *bescore.CommandLineSection) string {
	if section.GetOptionList() != nil {
		options := section.GetOptionList().GetOption()
		for _, option := range options {
			if option.GetOptionName() == "profile" {
				return option.GetOptionValue()
			}
		}
	}

	// Default value if --profile is not set
	return "command.profile.gz"
}
