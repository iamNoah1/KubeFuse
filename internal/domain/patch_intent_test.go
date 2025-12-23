package domain

import (
	"testing"
	"time"
)

func TestNewPatchIntentAppliesOptions(t *testing.T) {
	resource := ResourceRef{Kind: "deploy", Name: "web", Namespace: "prod"}
	patches := []Patch{NewPatch([]string{"spec", "replicas"}, NewInt(2))}

	intent := NewPatchIntent(resource, patches,
		WithReason("hotfix"),
		WithTTL(10*time.Minute),
		WithDryRun(true),
	)

	if intent.Resource != resource {
		t.Fatalf("expected resource %v, got %v", resource, intent.Resource)
	}

	if len(intent.Patches) != 1 {
		t.Fatalf("expected 1 patch, got %d", len(intent.Patches))
	}

	if intent.Reason != "hotfix" {
		t.Fatalf("expected reason 'hotfix', got %q", intent.Reason)
	}

	if intent.TTL != 10*time.Minute {
		t.Fatalf("expected ttl 10m, got %s", intent.TTL)
	}

	if !intent.DryRun {
		t.Fatalf("expected dry-run to be true")
	}
}
