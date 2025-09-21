package domain

import "testing"

func TestValue_BoolConstructorAndAccessor(t *testing.T) {
	v := NewBool(true)

	if v.Kind() != KindBool {
		t.Errorf("Kind() = %v, want %v", v.Kind(), KindBool)
	}

	raw, ok := v.ToInterface().(bool)
	if !ok || raw != true {
		t.Errorf("raw bool = (%v,%v), want (true,true)", raw, ok)
	}
}

func TestValue_IntegerConstructorAndAccessor(t *testing.T) {
	v := NewInt(3)

	if v.Kind() != KindInt {
		t.Errorf("Kind() = %v, want %v", v.Kind(), KindInt)
	}

	raw, ok := v.ToInterface().(int64) // int64 per constructor
	if !ok || raw != 3 {
		t.Errorf("raw int64 = (%v,%v), want (3,true)", raw, ok)
	}
}

func TestValue_StringConstructorAndAccessor(t *testing.T) {
	v := NewString("string")

	if v.Kind() != KindString {
		t.Errorf("Kind() = %v, want %v", v.Kind(), KindString)
	}

	raw, ok := v.ToInterface().(string)
	if !ok || raw != "string" {
		t.Errorf("raw string = (%v,%v), want (\"string\",true)", raw, ok)
	}
}
