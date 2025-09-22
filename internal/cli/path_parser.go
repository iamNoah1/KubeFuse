package cli

import (
	"errors"
	"strings"
)

func ParsePath(s string) ([]string, error) {
	if s == "" {
		return nil, errors.New("Empty path is not allowed")
	}

	split := strings.Split(s, ".")

	for _, elem := range split {
		if elem == "" {
			return nil, errors.New("Empty segment not allowed")
		}
	}

	return strings.Split(s, "."), nil
}
