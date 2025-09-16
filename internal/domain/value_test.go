package domain

import (
	"testing"
)

func TestValue_BoolConstructorAndAccessor(t *testing.T) {
	valueToTest := NewBool(true)

	if valueToTest.Kind() != KindBool {
		t.Errorf("Kind() = %v, want %v", valueToTest.Kind(), KindBool)
	}

	if got, ok := valueToTest.AsBool(); got != true || ok != true {
		t.Errorf("AsBool() = (%v,%v), want (%v,%v)", got, ok, true, true)
	}
}

func TestValue_IntegerConstructorAndAccessor(t *testing.T) {
	valueToTest := NewInt(3)

	if valueToTest.Kind() != KindInt {
		t.Errorf("Kind() = %v, want %v", valueToTest.Kind(), KindInt)
	}

	if got, ok := valueToTest.AsInt(); got != 3 || ok != true {
		t.Errorf("AsInt() = (%v,%v), want (%v,%v)", got, ok, 3, true)
	}
}

func TestValue_StringConstructorAndAccessor(t *testing.T) {
	valueToTest := NewString("string")

	if valueToTest.Kind() != KindString {
		t.Errorf("Kind() = %v, want %v", valueToTest.Kind(), KindString)
	}

	if got, ok := valueToTest.AsString(); got != "string" || ok != true {
		t.Errorf("AsString() = (%v,%v), want (%v,%v)", got, ok, "string", true)
	}
}
