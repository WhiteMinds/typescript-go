package ls

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetPackagePathComponents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filePath string
		expected []string
	}{
		{
			name:     "not in node_modules",
			filePath: "/home/user/project/src/index.ts",
			expected: nil,
		},
		{
			name:     "regular package",
			filePath: "/project/node_modules/foo/index.d.ts",
			expected: []string{"", "project", "node_modules", "foo"},
		},
		{
			name:     "scoped package",
			filePath: "/project/node_modules/@types/react/index.d.ts",
			expected: []string{"", "project", "node_modules", "@types", "react"},
		},
		{
			name:     "nested node_modules",
			filePath: "/project/node_modules/foo/node_modules/bar/index.d.ts",
			expected: []string{"", "project", "node_modules", "foo", "node_modules", "bar"},
		},
		{
			name:     "scoped nested node_modules",
			filePath: "/project/node_modules/@scope/pkg/node_modules/@inner/dep/index.d.ts",
			expected: []string{"", "project", "node_modules", "@scope", "pkg", "node_modules", "@inner", "dep"},
		},
		{
			name:     "node_modules at root",
			filePath: "node_modules/foo/index.ts",
			expected: []string{"node_modules", "foo"},
		},
		{
			name:     "node_modules with no package after",
			filePath: "/project/node_modules",
			expected: []string{"", "project", "node_modules"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := getPackagePathComponents(tt.filePath)
			assert.DeepEqual(t, result, tt.expected)
		})
	}
}
