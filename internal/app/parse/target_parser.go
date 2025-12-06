package parse

import (
	"errors"
	"kubefuse/internal/domain"
	"strings"
)

// ParseTarget parses raw target string and validates format
func ParseTarget(target string) ([]string, error) {
	if target == "" {
		return nil, errors.New("Empty target is not allowed")
	}

	split := strings.Split(target, "/")

	if len(split) != 2 {
		return nil, errors.New("More than 2 target segments are not allowed")
	}

	if split[0] == "" || split[1] == "" {
		return nil, errors.New("Empty target segments are not allowed")
	}

	return split, nil
}

// ParseTargetToResourceRef parses raw target string into a domain.ResourceRef
func ParseTargetToResourceRef(target string, namespace string) (*domain.ResourceRef, error) {
	parts, err := ParseTarget(target)
	if err != nil {
		return nil, err
	}

	ref := domain.NewResourceRef(parts[0], parts[1])
	if namespace != "" {
		ref.Namespace = namespace
	}

	return ref, nil
}
