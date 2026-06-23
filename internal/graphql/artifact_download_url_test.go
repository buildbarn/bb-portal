package graphql

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDownloadURLFor(t *testing.T) {
	sha256 := strings.Repeat("a", 64)
	sha1 := strings.Repeat("b", 40)
	md5 := strings.Repeat("c", 32)
	sha512 := strings.Repeat("d", 128)

	cases := []struct {
		name, uri, digest, file string
		size                    int64
		want                    string
	}{
		{
			"bytestream sha256 (function implied by 64-char digest)",
			"bytestream://bb/instance/blobs/" + sha256 + "/12", sha256, "foo.bin", 12,
			"/api/v1/servefile/instance/blobs/sha256/file/" + sha256 + "-12/foo.bin",
		},
		{
			"bytestream sha1 (function implied by 40-char digest)",
			"bytestream://bb/instance/blobs/" + sha1 + "/12", sha1, "foo.bin", 12,
			"/api/v1/servefile/instance/blobs/sha1/file/" + sha1 + "-12/foo.bin",
		},
		{
			"bytestream md5 (function implied by 32-char digest)",
			"bytestream://bb/instance/blobs/" + md5 + "/12", md5, "foo.bin", 12,
			"/api/v1/servefile/instance/blobs/md5/file/" + md5 + "-12/foo.bin",
		},
		{
			"bytestream sha512 (function implied by 128-char digest)",
			"bytestream://bb/instance/blobs/" + sha512 + "/12", sha512, "foo.bin", 12,
			"/api/v1/servefile/instance/blobs/sha512/file/" + sha512 + "-12/foo.bin",
		},
		{
			"http passthrough",
			"https://files.example.com/foo.bin", "", "foo.bin", 0,
			"https://files.example.com/foo.bin",
		},
		{
			"file scheme returns empty",
			"file:///tmp/foo", "", "foo", 0, "",
		},
		{
			"empty uri returns empty",
			"", sha256, "foo", 12, "",
		},
		{
			"bytestream missing digest",
			"bytestream://bb/instance/blobs/abc/12", "", "foo", 12, "",
		},
		{
			"bytestream unrecognized digest length returns empty",
			"bytestream://bb/instance/blobs/abc/12", "abc", "foo", 12, "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := downloadURLFor(tc.uri, tc.digest, tc.size, tc.file)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
