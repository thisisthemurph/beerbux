package semanticversion

import "testing"

func TestString(t *testing.T) {
	testCases := []struct {
		name     string
		sv       SemanticVersion
		expected string
	}{
		{
			name:     "Major version only",
			sv:       SemanticVersion{Major: 1},
			expected: "1",
		},
		{
			name:     "Major version only",
			sv:       SemanticVersion{Major: 0},
			expected: "0",
		},
		{
			name:     "Major and minor versions",
			sv:       SemanticVersion{Major: 1, Minor: 2},
			expected: "1.2",
		},
		{
			name:     "Major and minor versions",
			sv:       SemanticVersion{Major: 0, Minor: 2},
			expected: "0.2",
		},
		{
			name:     "Major, minor, and patch versions",
			sv:       SemanticVersion{Major: 1, Minor: 2, Patch: 3},
			expected: "1.2.3",
		},
		{
			name:     "Major, minor, and patch versions",
			sv:       SemanticVersion{Major: 0, Minor: 2, Patch: 3},
			expected: "0.2.3",
		},
		{
			name:     "Major, minor, and patch versions",
			sv:       SemanticVersion{Major: 1, Minor: 0, Patch: 3},
			expected: "1.0.3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.sv.String()
			if got != tc.expected {
				t.Errorf("String() = %q; want %q", got, tc.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		version  string
		expected SemanticVersion
	}{
		{
			name:     "Major version only",
			version:  "1",
			expected: SemanticVersion{Major: 1},
		},
		{
			name:     "Major version only",
			version:  "0",
			expected: SemanticVersion{Major: 0},
		},
		{
			name:     "Major and minor versions",
			version:  "1.2",
			expected: SemanticVersion{Major: 1, Minor: 2},
		},
		{
			name:     "Major and minor versions",
			version:  "0.2",
			expected: SemanticVersion{Major: 0, Minor: 2},
		},
		{
			name:     "Major, minor, and patch versions",
			version:  "1.2.3",
			expected: SemanticVersion{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:     "Major, minor, and patch versions",
			version:  "0.2.3",
			expected: SemanticVersion{Major: 0, Minor: 2, Patch: 3},
		},
		{
			name:     "Major, minor, and patch versions",
			version:  "1.0.3",
			expected: SemanticVersion{Major: 1, Minor: 0, Patch: 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.version)
			if err != nil {
				t.Fatalf("Parse() returned an error: %v", err)
			}

			if got != tc.expected {
				t.Errorf("Parse() = %v; want %v", got, tc.expected)
			}
		})
	}
}
