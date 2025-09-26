package domain

type Patch struct {
	Path  []string
	Value Value
}

func NewAssigment(path []string, value Value) Patch {
	return Patch{
		path,
		value,
	}
}
