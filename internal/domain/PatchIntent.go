package domain

type PatchIntent struct {
	resource ResourceRef
	patches  []Patch
}
