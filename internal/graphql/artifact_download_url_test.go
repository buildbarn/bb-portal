package graphql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDownloadURLFor(t *testing.T) {
	cases := []struct {
		name, uri, digest, fn, file string
		size                        int64
		want                        string
	}{
		{
			"bytestream",
			"bytestream://bb/instance/blobs/abc/12", "abc", "sha256", "foo.bin", 12,
			"/api/v1/servefile/instance/blobs/sha256/file/abc-12/foo.bin",
		},
		{
			"http passthrough",
			"https://files.example.com/foo.bin", "", "", "foo.bin", 0,
			"https://files.example.com/foo.bin",
		},
		{
			"file scheme returns empty",
			"file:///tmp/foo", "", "", "foo", 0, "",
		},
		{
			"empty uri returns empty",
			"", "abc", "sha256", "foo", 12, "",
		},
		{
			"bytestream missing digest",
			"bytestream://bb/instance/blobs/abc/12", "", "sha256", "foo", 12, "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := downloadURLFor(tc.uri, tc.digest, tc.fn, tc.size, tc.file)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
