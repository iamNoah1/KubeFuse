package schema

type GroupVersion struct {
	Group   string
	Version string
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

func (g GroupVersionKind) GroupKind() GroupKind {
	return GroupKind{Group: g.Group, Kind: g.Kind}
}

type GroupKind struct {
	Group string
	Kind  string
}

type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}
