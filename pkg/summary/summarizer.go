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

type Summarizer struct {
	summary         *Summary
	problemDetector detectors.ProblemDetector
}

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

func NewSummarizer() *Summarizer {
	return newSummarizer("", detectors.NewProblemDetector())
}

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

func (s Summarizer) FinishProcessing() (*Summary, error) {
	// If problems are ignored for the exit code, return immediately.
	if !shouldIgnoreProblems(s.summary.ExitCode) {
		// Add any detected test problems.
		problems, problemsErr := s.problemDetector.Problems()
		if problemsErr != nil {
			return nil, problemsErr
		}
		s.summary.Problems = append(s.summary.Problems, problems...)
	}

	return s.summary, nil
}

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

	case *bes.BuildEventId_StructuredCommandLine:
		err := s.handleStructuredCommandLine(buildEvent.GetStructuredCommandLine())
		if err != nil {
			return err
		}

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

func (s Summarizer) handleBuildMetadata(metadataProto *bes.BuildMetadata) {
	metadataMap := metadataProto.GetMetadata()
	if metadataMap == nil {
		return
	}
	stepLabel, ok := metadataMap[stepLabelKey]
	if !ok {
		return
	}
	s.summary.StepLabel = stepLabel
}

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

func (s Summarizer) handleStructuredCommandLine(structuredCommandLine *bescore.CommandLine) error {
	if structuredCommandLine.GetCommandLineLabel() != "original" {
		return nil
	}

	s.updateEnvVarsAndCommandFromStructuredCommandLine(structuredCommandLine)

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

	return nil
}

func (s Summarizer) handleOptionsParsed(optionsParsed *bes.OptionsParsed) {
	s.summary.InvocationSummary.BazelCommandLine.Options = optionsParsed.GetExplicitCmdLine()
}

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

func (s Summarizer) updateEnvVarsAndCommandFromStructuredCommandLine(structuredCommandLine *bescore.CommandLine) {
	sections := structuredCommandLine.GetSections()
	for _, section := range sections {
		label := section.GetSectionLabel()
		if label == "command options" {
			s.summary.InvocationSummary.EnvVars = map[string]string{}
			ParseEnvVarsFromSectionOptions(section, &s.summary.InvocationSummary.EnvVars)
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

func shouldIgnoreProblems(exitCode *ExitCode) bool {
	return exitCode != nil && (exitCode.Code == ExitCodeSuccess || exitCode.Code == ExitCodeInterrupted)
}

func envToI(envVars map[string]string, name string) (int, error) {
	res, err := strconv.Atoi(envVars[name])
	if err != nil {
		slog.Error("failed to parse env var to int", "envKey", name, "envValue", envVars[name], "err", err)
		return 0, fmt.Errorf("failed to parse %s (value: %s) as an int: %w", name, envVars[name], err)
	}
	return res, nil
}

func ParseEnvVarsFromSectionOptions(section *bescore.CommandLineSection, destMap *map[string]string) {
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
