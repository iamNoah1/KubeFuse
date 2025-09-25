package cli

import (
	"errors"
	"kubefuse/internal/domain"
	"strings"
)

func ParseTarget(target string) (*domain.ResourceRef, error) {
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

	return domain.NewResourceRef(split[0], split[1]), nil
}
