package parse

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

	if literal == "" {
		return domain.Value{}, fmt.Errorf("empty value not allowed")
	}

	switch {
	case literal == "true":
		return domain.NewBool(true), nil

	case literal == "false":
		return domain.NewBool(false), nil

	case literal == "null":
		// TODO: Add support for null/nil values if needed
		return domain.Value{}, fmt.Errorf("null values not yet supported")

	case isIntLiteral(literal):
		if i, err := strconv.ParseInt(literal, 10, 64); err == nil {
			return domain.NewInt(i), nil
		}
		return domain.Value{}, fmt.Errorf("integer out of range: %q", literal)

	default:
		// Treat everything else as a string
		return domain.NewString(literal), nil
	}
}

func isIntLiteral(s string) bool {
	return intPattern.MatchString(s)
}
