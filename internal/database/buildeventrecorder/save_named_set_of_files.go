package buildeventrecorder

import (
	"context"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/protobuf/proto"
)

// saveNamedSetOfFiles appends the BuildEvent to the in-memory streaming
// artifact-graph buffer. No-op unless save level is
// basic_and_target_and_artifacts (which is the only condition under which
// the buffer is allocated).
func (r *buildEventRecorder) saveNamedSetOfFiles(
	ctx context.Context,
	buildEvent *bes.BuildEvent,
) error {
	if r.artifactGraph == nil || buildEvent == nil {
		return nil
	}
	payload, err := proto.Marshal(buildEvent)
	if err != nil {
		return util.StatusWrap(err, "Failed to marshal NamedSetOfFiles BuildEvent")
	}
	if err := r.artifactGraph.AppendBuildEvent(payload); err != nil {
		return util.StatusWrap(err, "Failed to append NamedSetOfFiles to artifact graph buffer")
	}
	return nil
}
