package cli

import (
	"fmt"
	"kubefuse/internal/domain"
	"regexp"
	"strconv"
	"strings"
)

var intPattern = regexp.MustCompile(`^-?\d+$`)

func ParseLiteral(literal string) (domain.Value, error) {
	literal = strings.TrimSpace(literal)

	switch {
	case literal == "true":
		return domain.NewBool(true), nil

	case literal == "false":
		return domain.NewBool(false), nil

	case isIntLiteral(literal):
		if i, err := strconv.ParseInt(literal, 10, 64); err == nil {
			return domain.NewInt(i), nil
		}
		return domain.Value{}, fmt.Errorf("integer out of range: %q", literal)

	default:
		return domain.Value{}, fmt.Errorf("invalid literal: %q", literal)
	}
}

func isIntLiteral(s string) bool {
	return intPattern.MatchString(s)
}
