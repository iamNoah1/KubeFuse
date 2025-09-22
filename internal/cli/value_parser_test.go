package cli

import (
	"kubefuse/internal/domain"
	"testing"
)

func TestParser_Bool(t *testing.T) {
	domainValue, err := ParseLiteral("true")

	if err != nil {
		t.Errorf("Error while parsing")
	}

	if domainValue.Kind() != domain.KindBool {
		t.Errorf("Wrong type")
	}

	domainValue, err = ParseLiteral("false")

	if err != nil {
		t.Errorf("Error while parsing")
	}

	if domainValue.Kind() != domain.KindBool {
		t.Errorf("Wrong type")
	}

	domainValue, err = ParseLiteral("ture")

	if err == nil {
		t.Errorf("Should not be able to parse")
	}
}

func TestParser_Int(t *testing.T) {
	domainValue, err := ParseLiteral("2")

	if err != nil {
		t.Errorf("Error while parsing")
	}

	if domainValue.Kind() != domain.KindInt {
		t.Errorf("Wrong type")
	}

	domainValue, err = ParseLiteral("2,5")

	if err == nil {
		t.Errorf("Should not be able to parse")
	}

	domainValue, err = ParseLiteral("2.5")

	if err == nil {
		t.Errorf("Should not be able to parse")
	}
}
