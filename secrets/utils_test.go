package secrets

import (
	"testing"
)

func TestFormatPath(t *testing.T) {
	tests := []struct {
		name     string
		engine   string
		root     string
		path     string
		expected string
	}{
		{
			name:     "KV v1, no leading slashes",
			engine:   "kv-v1",
			root:     "secret",
			path:     "foo/bar",
			expected: "secret/foo/bar",
		},
		{
			name:     "KV v1, leading slash in root",
			engine:   "kv-v1",
			root:     "/secret",
			path:     "foo/bar",
			expected: "secret/foo/bar",
		},
		{
			name:     "KV v1, leading slash in path",
			engine:   "kv-v1",
			root:     "secret",
			path:     "/foo/bar",
			expected: "secret/foo/bar",
		},
		{
			name:     "KV v1, leading slashes in both",
			engine:   "kv-v1",
			root:     "/secret",
			path:     "/foo/bar",
			expected: "secret/foo/bar",
		},
		{
			name:     "KV v2, no leading slashes",
			engine:   "kv-v2",
			root:     "secret",
			path:     "foo/bar",
			expected: "secret/data/foo/bar",
		},
		{
			name:     "KV v2, leading slash in root",
			engine:   "kv-v2",
			root:     "/secret",
			path:     "foo/bar",
			expected: "secret/data/foo/bar",
		},
		{
			name:     "KV v2, leading slash in path",
			engine:   "kv-v2",
			root:     "secret",
			path:     "/foo/bar",
			expected: "secret/data/foo/bar",
		},
		{
			name:     "KV v2, leading slashes in both",
			engine:   "kv-v2",
			root:     "/secret",
			path:     "/foo/bar",
			expected: "secret/data/foo/bar",
		},
		{
			name:     "Empty root",
			engine:   "kv-v1",
			root:     "",
			path:     "foo/bar",
			expected: "/foo/bar",
		},
		{
			name:     "Empty path",
			engine:   "kv-v1",
			root:     "secret",
			path:     "",
			expected: "secret/",
		},
		{
			name:     "Empty root and path",
			engine:   "kv-v1",
			root:     "",
			path:     "",
			expected: "/",
		},
		{
			name:     "Unicode in path",
			engine:   "kv-v1",
			root:     "секрет",
			path:     "путь/значение",
			expected: "секрет/путь/значение",
		},
		{
			name:     "Unicode in root and path, kv-v2",
			engine:   "kv-v2",
			root:     "😭",
			path:     "путь/значение",
			expected: "😭/data/путь/значение",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPath(tt.engine, tt.root, tt.path)
			if got != tt.expected {
				t.Errorf("formatPath(%q, %q, %q) = %q; want %q", tt.engine, tt.root, tt.path, got, tt.expected)
			}
		})
	}
}
