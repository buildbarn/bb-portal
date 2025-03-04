package casproxy

import (
	"testing"
)

func TestGetInstanceName(t *testing.T) {
	tests := []struct {
		resourceName string
		expected     string
	}{
		{"/instance1/blobs/abc123", "instance1"},
		{"instance2/blobs/abc123", "instance2"},
		{"/instance3/blobs", "instance3"},
		{"instance4/blobs", "instance4"},
		{"/blobs/abc123", ""},
		{"/blobs", ""},
		{"", ""},
		{"/in/stance/5/blobs/abc123", "in/stance/5"},
		{"in/stance/6/blobs/abc123", "in/stance/6"},
		{"/in/stance/7/blobs", "in/stance/7"},
		{"in/stance/8/blobs", "in/stance/8"},
	}

	for _, test := range tests {
		t.Run(test.resourceName, func(t *testing.T) {
			result := getInstanceName(test.resourceName)
			if result != test.expected {
				t.Errorf("getInstanceName(%q) = %q; want %q", test.resourceName, result, test.expected)
			}
		})
	}
}
