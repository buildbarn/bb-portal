package graphql

import (
	"encoding/binary"
	"fmt"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/internal/graphql/model"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/proto"
)

// decodeArtifactGraphBlob decompresses and decodes the compacted BEP-graph
// blob into the structured GraphQL model. The blob is a zstd stream of
// length-prefixed serialized bes.BuildEvent messages (NamedSetOfFiles and
// TargetCompleted variants only), produced by
// dbcleanupservice.CompactArtifactGraphs.
//
// Decoding happens entirely server-side using the generated BES protos;
// the client receives structured data and never sees the compressed
// bytes or has to parse a wire format.
func decodeArtifactGraphBlob(compressed []byte) (*model.ArtifactGraph, error) {
	dec, err := zstd.NewReader(nil)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create zstd reader for artifact graph")
	}
	defer dec.Close()

	raw, err := dec.DecodeAll(compressed, nil)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to decompress artifact graph")
	}

	graph := newArtifactGraph()
	for len(raw) > 0 {
		msgLen, n := binary.Uvarint(raw)
		if n <= 0 {
			return nil, util.StatusWrap(fmt.Errorf("invalid length prefix"), "Failed to decode artifact graph")
		}
		raw = raw[n:]
		if uint64(len(raw)) < msgLen {
			return nil, util.StatusWrap(fmt.Errorf("length-prefixed message past end of stream"), "Failed to decode artifact graph")
		}
		if err := foldArtifactEventBytes(graph, raw[:msgLen]); err != nil {
			return nil, err
		}
		raw = raw[msgLen:]
	}
	return graph, nil
}

// decodeArtifactGraphEvents decodes the structured graph from the
// uncompacted staging events (one serialized bes.BuildEvent per element),
// so the graph can be served in its partial state before compaction has
// run.
func decodeArtifactGraphEvents(events [][]byte) (*model.ArtifactGraph, error) {
	graph := newArtifactGraph()
	for _, event := range events {
		if err := foldArtifactEventBytes(graph, event); err != nil {
			return nil, err
		}
	}
	return graph, nil
}

func newArtifactGraph() *model.ArtifactGraph {
	return &model.ArtifactGraph{
		NamedSets: []*model.ArtifactNamedSet{},
		Targets:   []*model.ArtifactTarget{},
	}
}

func foldArtifactEventBytes(graph *model.ArtifactGraph, raw []byte) error {
	var event bes.BuildEvent
	if err := proto.Unmarshal(raw, &event); err != nil {
		return util.StatusWrap(err, "Failed to unmarshal artifact graph BuildEvent")
	}
	appendArtifactEvent(graph, &event)
	return nil
}

// appendArtifactEvent folds a single BuildEvent into the graph, handling
// the two variants the recorder stores.
func appendArtifactEvent(graph *model.ArtifactGraph, event *bes.BuildEvent) {
	if nsf := event.GetNamedSetOfFiles(); nsf != nil {
		graph.NamedSets = append(graph.NamedSets, &model.ArtifactNamedSet{
			ID:          event.GetId().GetNamedSet().GetId(),
			Files:       artifactFiles(nsf.GetFiles()),
			ChildSetIds: namedSetIDs(nsf.GetFileSets()),
		})
		return
	}
	if completed := event.GetCompleted(); completed != nil {
		id := event.GetId().GetTargetCompleted()
		var aspect *string
		if a := id.GetAspect(); a != "" {
			aspect = &a
		}
		graph.Targets = append(graph.Targets, &model.ArtifactTarget{
			Label:        id.GetLabel(),
			Aspect:       aspect,
			OutputGroups: artifactOutputGroups(completed.GetOutputGroup()),
		})
	}
}

func artifactFiles(files []*bes.File) []*model.ArtifactFile {
	out := make([]*model.ArtifactFile, 0, len(files))
	for _, f := range files {
		af := &model.ArtifactFile{Name: f.GetName()}
		if uri := f.GetUri(); uri != "" {
			af.URI = &uri
		}
		if digest := f.GetDigest(); digest != "" {
			af.Digest = &digest
		}
		if length := f.GetLength(); length != 0 {
			sz := int(length)
			af.SizeBytes = &sz
		}
		out = append(out, af)
	}
	return out
}

func artifactOutputGroups(groups []*bes.OutputGroup) []*model.ArtifactOutputGroup {
	out := make([]*model.ArtifactOutputGroup, 0, len(groups))
	for _, g := range groups {
		out = append(out, &model.ArtifactOutputGroup{
			Name:       g.GetName(),
			Incomplete: g.GetIncomplete(),
			RootSetIds: namedSetIDs(g.GetFileSets()),
		})
	}
	return out
}

func namedSetIDs(sets []*bes.BuildEventId_NamedSetOfFilesId) []string {
	out := make([]string, 0, len(sets))
	for _, s := range sets {
		out = append(out, s.GetId())
	}
	return out
}
