package app

import "kubefuse/internal/domain"

type SetDTO struct {
	Kind          string          // from target parser
	Name          string          // from target parser
	NamespaceFlag string          // from CLI flag (can be "")
	Assignments   []AssignmentDTO // parsed path/value pairs
	Reason        string          // from --reason
	TTLRaw        string          // from --ttl (e.g. "10m"), not parsed yet
}

type AssignmentDTO struct {
	PathSegments []string     // from path parser (["spec","replicas"])
	Literal      domain.Value // from literal parser (int/bool/string/null)
}
