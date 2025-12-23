package domain

import "time"

type PatchIntent struct {
	Resource ResourceRef
	Patches  []Patch
	Reason   string
	TTL      time.Duration
	DryRun   bool
}

type PatchIntentOption func(*PatchIntent)

func NewPatchIntent(resource ResourceRef, patches []Patch, opts ...PatchIntentOption) PatchIntent {
	intent := PatchIntent{
		Resource: resource,
		Patches:  patches,
	}

	for _, opt := range opts {
		opt(&intent)
	}

	return intent
}

func WithReason(reason string) PatchIntentOption {
	return func(intent *PatchIntent) {
		intent.Reason = reason
	}
}

func WithTTL(ttl time.Duration) PatchIntentOption {
	return func(intent *PatchIntent) {
		intent.TTL = ttl
	}
}

func WithDryRun(dryRun bool) PatchIntentOption {
	return func(intent *PatchIntent) {
		intent.DryRun = dryRun
	}
}
