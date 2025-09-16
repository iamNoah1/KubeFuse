package domain

type ValueKind int

const (
	KindNull ValueKind = iota
	KindBool
	KindInt
	KindString
)

type Value struct {
	kind ValueKind
	raw  any
}

func (v Value) Kind() ValueKind {
	return v.kind
}

func NewBool(b bool) Value {
	return Value{
		KindBool,
		b,
	}
}

func NewInt(i int64) Value {
	return Value{
		KindInt,
		i,
	}
}

func NewString(s string) Value {
	return Value{
		KindString,
		s,
	}
}

func (v Value) AsBool() (bool, bool) {
	if v.kind != KindBool {
		return false, false
	}

	b, ok := v.raw.(bool)
	return b, ok
}

func (v Value) AsInt() (int64, bool) {
	if v.kind != KindInt {
		return 0, false
	}

	i, ok := v.raw.(int64)
	return i, ok
}

func (v Value) AsString() (string, bool) {
	if v.kind != KindString {
		return "", false
	}

	s, ok := v.raw.(string)
	return s, ok
}
