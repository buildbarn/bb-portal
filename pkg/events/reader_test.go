package events_test

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/buildbarn/bb-portal/pkg/events"
)

// TestBuildEventIteratorBasic A test build event iterator.
func TestBuildEventIteratorBasic(t *testing.T) {
	type setup struct {
		fileName string
	}
	tests := []struct {
		name           string
		setup          setup
		wantNumOfLines int
		wantErr        error
	}{
		{
			name: "bazelbuild/examples/cpp-tutorial/stage1/build.bep.ndjson",
			setup: setup{
				fileName: "bazelbuild/examples/cpp-tutorial/stage1/build.bep.ndjson",
			},
			wantNumOfLines: 25,
		},
		{
			name: "bazelbuild/examples/cpp-tutorial/stage1/test.bep.ndjson",
			setup: setup{
				fileName: "bazelbuild/examples/cpp-tutorial/stage1/test.bep.ndjson",
			},
			wantNumOfLines: 27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join("testdata", tt.setup.fileName)
			file, err := os.Open(filePath)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, file.Close()) })

			it := events.NewBuildEventIterator(context.Background(), file)
			require.NoError(t, err)

			readBytesFile, err := os.Open(filePath)
			require.NoError(t, err)
			scanner := bufio.NewScanner(readBytesFile)

			lineCount := 0
			for {
				buildEvent, err := it.Next()
				if errors.Is(err, iterator.Done) {
					break
				}

				if tt.wantErr != nil && err != nil {
					require.Contains(t, err.Error(), "unexpected token")
					return
				}
				require.NoError(t, err)

				if lineCount == 0 {
					started := buildEvent.GetStarted()
					require.NotEmpty(t, started)
				}

				// Check if original bytes are kept.
				readStatus := scanner.Scan()
				require.True(t, readStatus)
				expectedBytes := scanner.Bytes()
				require.Equal(t, expectedBytes, []byte(buildEvent.RawMessage()))

				lineCount++
			}

			require.Equal(t, tt.wantNumOfLines, lineCount)
		})
	}
}

// TestBuildEventIterator_SpecificEvent A test build event iterator for a specific event.
func TestBuildEventIterator_SpecificEvent(t *testing.T) {
	filePath := filepath.Join("testdata", "namedSetOfFiles.ndjson")
	file, err := os.Open(filePath)
	require.NoError(t, err)

	readBytesFile, err := os.Open(filePath)
	require.NoError(t, err)
	scanner := bufio.NewScanner(readBytesFile)
	readStatus := scanner.Scan()
	require.True(t, readStatus)
	expectedBytes := scanner.Bytes()

	expectedBESEvent := &bes.BuildEvent{
		Id: &bes.BuildEventId{
			Id: &bes.BuildEventId_NamedSet{
				NamedSet: &bes.BuildEventId_NamedSetOfFilesId{
					Id: "0",
				},
			},
		},
		Payload: &bes.BuildEvent_NamedSetOfFiles{
			NamedSetOfFiles: &bes.NamedSetOfFiles{},
		},
	}
	expectedBuildEvent := events.NewBuildEvent(expectedBESEvent, expectedBytes)
	expectedBuildEvents := []*events.BuildEvent{
		&expectedBuildEvent,
	}

	it := events.NewBuildEventIterator(context.Background(), file)
	var buildEvent *events.BuildEvent
	for _, expectedEvent := range expectedBuildEvents {
		buildEvent, err = it.Next()
		require.NoError(t, err)

		// Compare proto parts separately since reflect.DeepEqual does not always work for them.
		require.True(t, proto.Equal(expectedEvent.BuildEvent, buildEvent.BuildEvent))

		// Compare the rest of events.
		buildEvent.BuildEvent = expectedEvent.BuildEvent
		require.Equal(t, expectedEvent, buildEvent)
	}
	buildEvent, err = it.Next()
	require.Error(t, err, iterator.Done)
	require.Nil(t, buildEvent)
}

// TestBuildEventIterator_EmptyLinesNotAllowed Prevents Empty Lines.
func TestBuildEventIterator_EmptyLinesNotAllowed(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "bazelbuild/examples/cpp-tutorial/stage1/build.bep.ndjson"))
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")
	lines = slices.Insert(lines, len(lines)/2, "")
	testContent := strings.Join(lines, "\n")

	reader := strings.NewReader(testContent)
	it := events.NewBuildEventIterator(context.Background(), reader)
	require.NoError(t, err)

	lineCount := 0
	for {
		_, err = it.Next()
		if err != nil {
			break
		}
		if errors.Is(err, iterator.Done) {
			t.Fatal("Expected an error")
		}
		lineCount++
	}

	assert.ErrorIs(t, err, proto.Error)
	assert.Contains(t, err.Error(), "unexpected token")
}

// TestBuildEventIterator_FinalNewlineIsOptional Final newline is optional helper.
func TestBuildEventIterator_FinalNewlineIsOptional(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "bazelbuild/examples/cpp-tutorial/stage1/build.bep.ndjson"))
	require.NoError(t, err)

	originalContent := string(content)
	require.True(t, strings.HasSuffix(originalContent, "\n"))
	tests := map[string]string{
		"original":              originalContent,
		"without final newline": strings.TrimSuffix(originalContent, "\n"),
	}

	for name, testContent := range tests {
		t.Run(name, func(t *testing.T) {
			reader := strings.NewReader(testContent)

			it := events.NewBuildEventIterator(context.Background(), reader)

			lineCount := 0
			for {
				_, err = it.Next()
				if errors.Is(err, iterator.Done) {
					break
				}
				require.NoError(t, err)
				lineCount++
			}

			assert.Equal(t, lineCount, 25)
		})
	}
}
