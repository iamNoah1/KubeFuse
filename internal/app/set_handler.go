package app

import (
	"context"
	"fmt"
	"github.com/iamNoah1/KubeFuse/internal/app/parse"
	"github.com/iamNoah1/KubeFuse/internal/domain"
	"time"
)

func SetHandler(dto SetDTO) error {
	resourceRef, err := parse.ParseTargetToResourceRef(dto.TargetRaw, dto.NamespaceFlag)
	if err != nil {
		return fmt.Errorf("failed to parse target: %w", err)
	}

	patches, err := parse.ParsePatches(dto.PatchesRaw)
	if err != nil {
		return fmt.Errorf("failed to parse patches: %w", err)
	}

	ttl, err := time.ParseDuration(dto.TTL)
	if err != nil {
		return fmt.Errorf("invalid TTL value %q: %w", dto.TTL, err)
	}

	intent := domain.NewPatchIntent(*resourceRef, patches,
		domain.WithReason(dto.Reason),
		domain.WithTTL(ttl),
		domain.WithDryRun(dto.DryRun),
	)

	executor, err := NewPatchExecutor()
	if err != nil {
		return err
	}

	return executor.ExecutePatchIntent(context.Background(), intent)
}
