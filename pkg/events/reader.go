package events

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/buildbarn/bb-portal/third_party/bazel/gen/bes"
)

const (
	// initialBufferSize is the starting size of the scanner buffer.
	initialBufferSize = 20 * 1024 * 1024
	// maxBufferSize is the maximum size the buffer can grow to before returning an error.
	maxBufferSize = 100 * 1024 * 1024

	UndeclaredTestOutputsName = "test.outputs__outputs.zip"
)

// BuildEventIterator follows [googleapi's Iterator Guidelines](https://github.com/googleapis/google-cloud-go/wiki/Iterator-Guidelines).
// End of iteration is indicated with Done error from `google.golang.org/api/iterator` package being returned; it is returned for any subsequent calls, too.
// No pagination is provided.
type BuildEventIterator struct {
	ctx         context.Context
	scanner     *bufio.Scanner
	unmarshaler protojson.UnmarshalOptions
}

type BuildEvent struct {
	*bes.BuildEvent
	rawEvent json.RawMessage
}

// NewBuildEvent creates a BuildEvent.
func NewBuildEvent(event *bes.BuildEvent, eventBytes json.RawMessage) BuildEvent {
	return BuildEvent{
		BuildEvent: event,
		rawEvent:   eventBytes,
	}
}

// Clone returns a deep copy of the BuildEvent. In particular this avoids issues with the underlying rawEvent buffer
// changing.
func (e BuildEvent) Clone() BuildEvent {
	copiedEvent := proto.Clone(e.BuildEvent).(*bes.BuildEvent)

	copiedRawEvent := make(json.RawMessage, len(e.rawEvent))
	copy(copiedRawEvent, e.rawEvent)

	return BuildEvent{
		BuildEvent: copiedEvent,
		rawEvent:   copiedRawEvent,
	}
}

func (e BuildEvent) RawMessage() json.RawMessage {
	return e.rawEvent
}

func AsJSONArray(buildEvents []*BuildEvent) (jsonb json.RawMessage, err error) {
	rawEvents := []json.RawMessage{}
	for _, event := range buildEvents {
		rawEvents = append(rawEvents, event.RawMessage())
	}
	// Create JSON list.
	var bepEvents []byte
	bepEvents, err = json.Marshal(rawEvents)
	if err != nil {
		return
	}
	jsonb = bepEvents
	return
}

func FromJSONArray(events json.RawMessage) (buildEvents []BuildEvent, err error) {
	// Unpack list of bytes.
	var rawEvents []json.RawMessage
	if err = json.Unmarshal(events, &rawEvents); err != nil {
		return
	}

	// Unmarshal every element as build event.
	bepUnmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	for _, rawEvent := range rawEvents {
		bepEvent := bes.BuildEvent{}
		err = bepUnmarshaler.Unmarshal(rawEvent, &bepEvent)
		if err != nil {
			return
		}
		buildEvent := NewBuildEvent(&bepEvent, rawEvent).Clone()
		buildEvents = append(buildEvents, buildEvent)
	}
	return
}

func (e BuildEvent) IsActionCompleted() bool {
	return e.GetId() != nil &&
		e.GetId().GetActionCompleted() != nil
}

// IsTargetConfigured returns true if this event is `targetConfigured` event.
func (e BuildEvent) IsTargetConfigured() bool {
	return e.GetId() != nil &&
		e.GetId().GetTargetConfigured() != nil
}

func (e BuildEvent) IsTargetCompleted() bool {
	return e.GetId() != nil &&
		e.GetId().GetTargetCompleted() != nil
}

func (e BuildEvent) IsTestSummary() bool {
	return e.GetId() != nil &&
		e.GetId().GetTestSummary() != nil &&
		e.GetTestSummary() != nil
}

func (e BuildEvent) IsTestResult() bool {
	return e.GetId().GetTestResult() != nil &&
		e.GetTestResult() != nil
}

func (e BuildEvent) IsWorkspaceStatus() bool {
	return e.GetId() != nil && e.GetId().GetWorkspaceStatus() != nil && e.GetWorkspaceStatus() != nil
}

func (e BuildEvent) IsStructuredCommandLine() bool {
	return e.GetId() != nil && e.GetId().GetStructuredCommandLine() != nil && e.GetStructuredCommandLine() != nil
}

// GetTargetConfiguredLabel returns label of this event if it is `targetConfigured` event.
// Otherwise it returns an empty string.
func (e BuildEvent) GetTargetConfiguredLabel() string {
	if !e.IsTargetConfigured() {
		return ""
	}
	return e.GetId().GetTargetConfigured().GetLabel()
}

func (e BuildEvent) GetTargetCompletedLabel() string {
	if !e.IsTargetCompleted() {
		return ""
	}
	return e.GetId().GetTargetCompleted().GetLabel()
}

func (e BuildEvent) GetActionCompletedLabel() string {
	if !e.IsActionCompleted() {
		return ""
	}
	return e.GetId().GetActionCompleted().GetLabel()
}

func (e BuildEvent) FindUndeclaredTestOutputsURI() string {
	if !e.IsTestResult() {
		return ""
	}
	for _, output := range e.GetTestResult().GetTestActionOutput() {
		if output.GetName() == UndeclaredTestOutputsName {
			return output.GetUri()
		}
	}
	return ""
}

func NewBuildEventIterator(ctx context.Context, reader io.Reader) *BuildEventIterator {
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, initialBufferSize), maxBufferSize)
	it := BuildEventIterator{
		ctx:         ctx,
		scanner:     scanner,
		unmarshaler: unmarshaler,
	}
	return &it
}

// Next returns the next result. Its second return value is iterator.Done if there are no more results.
// Once Next returns Done, all subsequent calls will return Done.
// If unknown field is encountered, it is ignored and no error is returned.
//
// NOTE: Currently Next calls BuildEvent.Clone so it is safe for consumers to save the BuildEvent for later.
// In the future we may change this so that it is the consumers responsibility to clone the event if they are
// using it as more than a temporary variable inside a single iteration of the loop.
func (it *BuildEventIterator) Next() (*BuildEvent, error) {
	if !it.scanner.Scan() {
		err := it.scanner.Err()
		if err == nil {
			err = iterator.Done
		}
		return nil, err
	}

	lineBytes, err := it.scanner.Bytes(), it.scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to read a line from build event file: %w", err)
	}

	bepEvent := &bes.BuildEvent{}
	err = it.unmarshaler.Unmarshal(lineBytes, bepEvent)
	if err != nil {
		slog.Warn("failed to unmarshal JSON", "err", err)
		slog.Debug("JSON that could not be unmarshalled", "content", string(lineBytes))
		// The enum below has been removed and protojson fails to parse message containing it.
		// TODO: Remove the condition once https://github.com/golang/protobuf/issues/1208 is fixed.
		if !strings.Contains(err.Error(), "invalid value for enum type: \"TRIGGERED_BY_ALL_INCOMPATIBLE_CHANGES\"") {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	buildEvent := NewBuildEvent(bepEvent, lineBytes).Clone()
	return &buildEvent, nil
}
