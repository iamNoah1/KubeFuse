package v1

import "time"

type PatchOptions struct {
	FieldManager string
	DryRun       []string
}

const DryRunAll = "All"

type GetOptions struct{}

// Duration stringer wrapper
func (d Duration) String() string { return time.Duration(d).String() }

type Duration time.Duration
