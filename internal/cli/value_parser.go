package cli

import (
	"fmt"
	"kubefuse/internal/domain"
	"regexp"
	"strconv"
	"strings"
)

var intPattern = regexp.MustCompile(`^-?\d+$`)

func ParseLiteral(s string) (domain.Value, error) {
	s = strings.TrimSpace(s)

	switch {
	case s == "true":
		return domain.NewBool(true), nil

	case s == "false":
		return domain.NewBool(false), nil

	case isIntLiteral(s):
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return domain.NewInt(i), nil
		}
		return domain.Value{}, fmt.Errorf("integer out of range: %q", s)

	default:
		return domain.Value{}, fmt.Errorf("invalid literal: %q", s)
	}
}

func isIntLiteral(s string) bool {
	return intPattern.MatchString(s)
}
