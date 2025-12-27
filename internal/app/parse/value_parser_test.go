package parse

import (
	"github.com/iamNoah1/KubeFuse/internal/domain"
	"testing"
)

func TestParseLiteral_BoolTrue(t *testing.T) {
	domainValue, err := ParseLiteral("true")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindBool {
		t.Errorf("Expected KindBool, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != true {
		t.Errorf("Expected true, got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_BoolFalse(t *testing.T) {
	domainValue, err := ParseLiteral("false")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindBool {
		t.Errorf("Expected KindBool, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != false {
		t.Errorf("Expected false, got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_Int(t *testing.T) {
	domainValue, err := ParseLiteral("2")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindInt {
		t.Errorf("Expected KindInt, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != int64(2) {
		t.Errorf("Expected 2, got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_NegativeInt(t *testing.T) {
	domainValue, err := ParseLiteral("-42")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindInt {
		t.Errorf("Expected KindInt, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != int64(-42) {
		t.Errorf("Expected -42, got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_String(t *testing.T) {
	domainValue, err := ParseLiteral("nginx:1.21")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindString {
		t.Errorf("Expected KindString, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != "nginx:1.21" {
		t.Errorf("Expected 'nginx:1.21', got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_StringWithSpaces(t *testing.T) {
	domainValue, err := ParseLiteral("hello world")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindString {
		t.Errorf("Expected KindString, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != "hello world" {
		t.Errorf("Expected 'hello world', got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_DecimalAsString(t *testing.T) {
	// Decimals are not integers, so they should be treated as strings
	domainValue, err := ParseLiteral("2.5")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindString {
		t.Errorf("Expected KindString, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != "2.5" {
		t.Errorf("Expected '2.5', got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_Null(t *testing.T) {
	domainValue, err := ParseLiteral("null")

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindNull {
		t.Fatalf("Expected KindNull, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != nil {
		t.Fatalf("Expected nil value, got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_TypoInBool(t *testing.T) {
	// "ture" is a typo, should be treated as a string
	domainValue, err := ParseLiteral("ture")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if domainValue.Kind() != domain.KindString {
		t.Errorf("Expected KindString, got: %v", domainValue.Kind())
	}

	if domainValue.ToInterface() != "ture" {
		t.Errorf("Expected 'ture', got: %v", domainValue.ToInterface())
	}
}

func TestParseLiteral_EmptyString(t *testing.T) {
	_, err := ParseLiteral("")

	if err == nil {
		t.Error("Expected error for empty string, but got none")
	}

	if err.Error() != "empty value not allowed" {
		t.Errorf("Expected error 'empty value not allowed', got: %s", err.Error())
	}
}

func TestParseLiteral_WhitespaceOnly(t *testing.T) {
	_, err := ParseLiteral("   ")

	if err == nil {
		t.Error("Expected error for whitespace-only string, but got none")
	}
}
