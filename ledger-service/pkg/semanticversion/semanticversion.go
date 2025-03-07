package semanticversion

import (
	"fmt"
	"strconv"
	"strings"
)

// SemanticVersion represents a semantic versioning structure.
type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

// String returns the string representation of a SemanticVersion.
func (sv SemanticVersion) String() string {
	if sv.Minor == 0 && sv.Patch == 0 {
		return fmt.Sprintf("%d", sv.Major)
	} else if sv.Patch == 0 {
		return fmt.Sprintf("%d.%d", sv.Major, sv.Minor)
	}
	return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
}

// Parse converts a version string (e.g., "1.2.3") into a SemanticVersion struct.
// If the version string is invalid, an error is returned.
func Parse(version string) (SemanticVersion, error) {
	version = strings.TrimSpace(version)
	parts := strings.Split(version, ".")

	if len(parts) < 1 || len(parts) > 3 {
		return SemanticVersion{}, fmt.Errorf("invalid version format: %s", version)
	}

	var sv SemanticVersion
	var err error

	// Parse major version
	if sv.Major, err = strconv.Atoi(parts[0]); err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid major version: %s", parts[0])
	}

	// Parse minor version (if present)
	if len(parts) > 1 {
		if sv.Minor, err = strconv.Atoi(parts[1]); err != nil {
			return SemanticVersion{}, fmt.Errorf("invalid minor version: %s", parts[1])
		}
	}

	// Parse patch version (if present)
	if len(parts) > 2 {
		if sv.Patch, err = strconv.Atoi(parts[2]); err != nil {
			return SemanticVersion{}, fmt.Errorf("invalid patch version: %s", parts[2])
		}
	}

	return sv, nil
}
