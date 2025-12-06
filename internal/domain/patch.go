package domain

type Patch struct {
	Path  []string
	Value Value
}

func NewPatch(path []string, value Value) Patch {
	return Patch{
		Path:  path,
		Value: value,
	}
}
