package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestUrlToPath(t *testing.T) {
	tests := []struct {
		url  string
		path []string
	}{{
		"example.com/foo",
		[]string{"example.com", "foo"},
	}, {
		"https://example.com/foo",
		[]string{"example.com", "foo"},
	}, {
		"git@example.com/foo",
		[]string{"example.com", "foo"},
	}, {
		"https://git@example.com/foo",
		[]string{"example.com", "foo"},
	}, {
		"example.com/foo/bar",
		[]string{"example.com", "foo", "bar"},
	}, {
		"example.com/~foo/bar",
		[]string{"example.com", "~foo", "bar"},
	}, {
		"example.com/foo/bar.git",
		[]string{"example.com", "foo", "bar"},
	}}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			path, err := urlToPath(tt.url)
			if err != nil {
				t.Fatal(err)
			}
			assert.DeepEqual(t, path, tt.path)
		})
	}
}
