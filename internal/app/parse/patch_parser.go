package parse

import (
	"errors"
	"github.com/iamNoah1/KubeFuse/internal/domain"
	"strings"
)

// ParsePath parses a dot-separated path string into segments
// Example: "spec.replicas" -> ["spec", "replicas"]
func ParsePath(pathStr string) ([]string, error) {
	if pathStr == "" {
		return nil, errors.New("Empty path is not allowed")
	}

	segments := strings.Split(pathStr, ".")

	for _, segment := range segments {
		if segment == "" {
			return nil, errors.New("Empty segment not allowed")
		}
	}

	return segments, nil
}

// ParsePatchString parses a single patch string (path=value) into path and value strings
// Example: "spec.replicas=3" -> (["spec", "replicas"], "3", nil)
func ParsePatchString(patchStr string) ([]string, string, error) {
	if patchStr == "" {
		return nil, "", errors.New("Empty patch is not allowed")
	}

	parts := strings.SplitN(patchStr, "=", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("Patch must be in format path=value")
	}

	pathStr := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	if pathStr == "" {
		return nil, "", errors.New("Empty path is not allowed")
	}

	if valueStr == "" {
		return nil, "", errors.New("Empty value is not allowed")
	}

	path, err := ParsePath(pathStr)
	if err != nil {
		return nil, "", err
	}

	return path, valueStr, nil
}

// ParsePatch parses a single patch string into a domain.Patch object
// Example: "spec.replicas=3" -> domain.Patch{Path: ["spec", "replicas"], Value: domain.Value{...}}
func ParsePatch(patchStr string) (domain.Patch, error) {
	path, valueStr, err := ParsePatchString(patchStr)
	if err != nil {
		return domain.Patch{}, err
	}

	value, err := ParseLiteral(valueStr)
	if err != nil {
		return domain.Patch{}, err
	}

	return domain.NewPatch(path, value), nil
}

// ParsePatches parses multiple patch strings into domain.Patch objects
// Example: ["spec.replicas=3", "spec.paused=true"] -> []domain.Patch{...}
func ParsePatches(patchesRaw []string) ([]domain.Patch, error) {
	if len(patchesRaw) == 0 {
		return nil, errors.New("No patches provided")
	}

	patches := make([]domain.Patch, 0, len(patchesRaw))

	for _, patchStr := range patchesRaw {
		patch, err := ParsePatch(patchStr)
		if err != nil {
			return nil, err
		}
		patches = append(patches, patch)
	}

	return patches, nil
}
