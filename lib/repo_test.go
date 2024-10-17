package repo

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
	}, {
		"git@example.com/foo/bar.git",
		[]string{"example.com", "foo", "bar"},
	}, {
		"git@example.com:foo/bar.git",
		[]string{"example.com", "foo", "bar"},
	}, {
		"git:p4ssw0rd@example.com/foo/bar.git",
		[]string{"example.com", "foo", "bar"},
	}, {
		"git@example.com:foo?bar=baz",
		[]string{"example.com", "foo"},
	}, {
		"git://git@example.com/foo?bar=baz",
		[]string{"example.com", "foo"},
	}, {
		"https://git:p4ssw0rd@example.com/foo/bar.git?biz=baz",
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

func TestSplitUrl(t *testing.T) {
	tests := []struct {
		url   string
		parts []string
	}{{
		"example.com/foo",
		[]string{"", "example.com/foo", ""},
	}, {
		"https://example.com/foo",
		[]string{"https://", "example.com/foo", ""},
	}, {
		"git@example.com/foo",
		[]string{"git@", "example.com/foo", ""},
	}, {
		"https://git@example.com/foo",
		[]string{"https://git@", "example.com/foo", ""},
	}, {
		"example.com/foo/bar",
		[]string{"", "example.com/foo/bar", ""},
	}, {
		"example.com/~foo/bar",
		[]string{"", "example.com/~foo/bar", ""},
	}, {
		"example.com/foo/bar.git",
		[]string{"", "example.com/foo/bar", ".git"},
	}, {
		"git@example.com/foo/bar.git",
		[]string{"git@", "example.com/foo/bar", ".git"},
	}, {
		"git@example.com:foo/bar.git",
		[]string{"git@", "example.com/foo/bar", ".git"},
	}, {
		"git:p4ssw0rd@example.com/foo/bar.git",
		[]string{"git:p4ssw0rd@", "example.com/foo/bar", ".git"},
	}, {
		"git@example.com:foo?bar=baz",
		[]string{"git@", "example.com/foo", "?bar=baz"},
	}, {
		"git://git@example.com/foo?bar=baz",
		[]string{"git://git@", "example.com/foo", "?bar=baz"},
	}, {
		"https://git:p4ssw0rd@example.com/foo/bar.git?biz=baz",
		[]string{"https://git:p4ssw0rd@", "example.com/foo/bar",
			".git?biz=baz"},
	}}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			prefix, path, suffix := splitUrl(tt.url)
			assert.DeepEqual(t, []string{prefix, path, suffix}, tt.parts)
		})
	}
}

func TestMergeUrl(t *testing.T) {
	urlpairs := [][]string{
		{"example.com/foo", "https://example.com/foo"},
		{"https://example.com/foo", "https://example.com/foo"},
		{"git@example.com:foo", "git@example.com:foo"},
		{"https://git@example.com/foo", "https://git@example.com/foo"},
		{"example.com/foo/bar", "https://example.com/foo/bar"},
		{"example.com/~foo/bar", "https://example.com/~foo/bar"},
		{"example.com/foo/bar.git", "https://example.com/foo/bar.git"},
		{"git@example.com/foo/bar.git", "git@example.com:foo/bar.git"},
		{"git@example.com:foo/bar.git", "git@example.com:foo/bar.git"},
		{
			"git:p4ssw0rd@example.com/foo/bar.git",
			"git:p4ssw0rd@example.com:foo/bar.git",
		},
		{"git@example.com:foo?bar=baz", "git@example.com:foo?bar=baz"},
		{
			"git://git@example.com/foo?bar=baz",
			"git://git@example.com/foo?bar=baz",
		},
		{
			"https://git:p4ssw0rd@example.com/foo/bar.git?biz=baz",
			"https://git:p4ssw0rd@example.com/foo/bar.git?biz=baz",
		},
	}
	for _, urlpair := range urlpairs {
		t.Run(urlpair[0], func(t *testing.T) {
			assert.Equal(t, urlpair[1], mergeUrl(splitUrl(urlpair[0])))
		})
	}
}
