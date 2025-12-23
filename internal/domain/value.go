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

func NewNull() Value {
	return Value{
		KindNull,
		nil,
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

func (v Value) ToInterface() any {
	return v.raw
}
