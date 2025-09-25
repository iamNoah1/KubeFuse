package domain

type ResourceRef struct {
	Kind      string
	Name      string
	Namespace string
}

func NewResourceRef(kind string, name string) *ResourceRef {
	return &ResourceRef{
		Kind: kind,
		Name: name,
	}
}
