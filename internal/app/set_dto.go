package app

type SetDTO struct {
	TargetRaw     string
	PatchesRaw    []string
	NamespaceFlag string // from CLI flag (can be "")
	Reason        string // from --reason
	TTL           string // from --ttl (e.g. "10m"), not parsed yet
	DryRun        bool
}
