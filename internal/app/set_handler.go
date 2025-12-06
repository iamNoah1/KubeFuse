package app

import (
	"fmt"
	"kubefuse/internal/app/parse"
)

func SetHandler(dto SetDTO) error {
	// Parse target into domain.ResourceRef
	resourceRef, err := parse.ParseTargetToResourceRef(dto.TargetRaw, dto.NamespaceFlag)
	if err != nil {
		return fmt.Errorf("failed to parse target: %w", err)
	}

	fmt.Printf("Parsed resource: Kind=%s, Name=%s, Namespace=%s\n", resourceRef.Kind, resourceRef.Name, resourceRef.Namespace)

	// Parse patches into []domain.Patch
	patches, err := parse.ParsePatches(dto.PatchesRaw)
	if err != nil {
		return fmt.Errorf("failed to parse patches: %w", err)
	}

	fmt.Printf("Parsed %d patches\n", len(patches))
	for i, patch := range patches {
		fmt.Printf("  Patch %d: Path=%v, Value=%v\n", i+1, patch.Path, patch.Value.ToInterface())
	}

	// TODO: Create PatchIntent and execute
	// intent := domain.NewPatchIntent(resourceRef, patches)
	// return executePatchIntent(intent)

	return nil
}
