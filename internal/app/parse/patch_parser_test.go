package parse

import (
	"testing"
)

// Tests for ParsePath
func TestParsePath_ValidPath(t *testing.T) {
	path, err := ParsePath("spec.replicas")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(path) != 2 {
		t.Errorf("Expected 2 segments, got: %d", len(path))
	}

	if path[0] != "spec" || path[1] != "replicas" {
		t.Errorf("Expected ['spec', 'replicas'], got: %v", path)
	}
}

func TestParsePath_NestedPath(t *testing.T) {
	path, err := ParsePath("spec.template.spec.containers")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(path) != 4 {
		t.Errorf("Expected 4 segments, got: %d", len(path))
	}

	expected := []string{"spec", "template", "spec", "containers"}
	for i, segment := range path {
		if segment != expected[i] {
			t.Errorf("Expected segment %d to be '%s', got '%s'", i, expected[i], segment)
		}
	}
}

func TestParsePath_SingleSegment(t *testing.T) {
	path, err := ParsePath("replicas")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(path) != 1 {
		t.Errorf("Expected 1 segment, got: %d", len(path))
	}

	if path[0] != "replicas" {
		t.Errorf("Expected 'replicas', got: %s", path[0])
	}
}

func TestParsePath_EmptyString(t *testing.T) {
	_, err := ParsePath("")

	if err == nil {
		t.Error("Expected error for empty path, but got none")
	}

	if err.Error() != "Empty path is not allowed" {
		t.Errorf("Expected error: 'Empty path is not allowed', got: %s", err.Error())
	}
}

func TestParsePath_EmptySegment(t *testing.T) {
	_, err := ParsePath("spec..replicas")

	if err == nil {
		t.Error("Expected error for empty segment, but got none")
	}

	if err.Error() != "Empty segment not allowed" {
		t.Errorf("Expected error: 'Empty segment not allowed', got: %s", err.Error())
	}
}

// Tests for ParsePatchString
func TestParsePatchString_ValidPatch(t *testing.T) {
	path, value, err := ParsePatchString("spec.replicas=3")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(path) != 2 {
		t.Errorf("Expected 2 path segments, got: %d", len(path))
	}

	if path[0] != "spec" || path[1] != "replicas" {
		t.Errorf("Expected ['spec', 'replicas'], got: %v", path)
	}

	if value != "3" {
		t.Errorf("Expected value '3', got: %s", value)
	}
}

func TestParsePatchString_BooleanValue(t *testing.T) {
	path, value, err := ParsePatchString("spec.paused=true")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(path) != 2 {
		t.Errorf("Expected 2 path segments, got: %d", len(path))
	}

	if value != "true" {
		t.Errorf("Expected value 'true', got: %s", value)
	}
}

func TestParsePatchString_WithWhitespace(t *testing.T) {
	path, value, err := ParsePatchString("  spec.replicas  =  3  ")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if path[0] != "spec" || path[1] != "replicas" {
		t.Errorf("Expected ['spec', 'replicas'], got: %v", path)
	}

	if value != "3" {
		t.Errorf("Expected value '3', got: %s", value)
	}
}

func TestParsePatchString_EmptyPatch(t *testing.T) {
	_, _, err := ParsePatchString("")

	if err == nil {
		t.Error("Expected error for empty patch, but got none")
	}

	if err.Error() != "Empty patch is not allowed" {
		t.Errorf("Expected error: 'Empty patch is not allowed', got: %s", err.Error())
	}
}

func TestParsePatchString_NoEqualsSign(t *testing.T) {
	_, _, err := ParsePatchString("spec.replicas")

	if err == nil {
		t.Error("Expected error for missing equals sign, but got none")
	}

	if err.Error() != "Patch must be in format path=value" {
		t.Errorf("Expected error: 'Patch must be in format path=value', got: %s", err.Error())
	}
}

func TestParsePatchString_EmptyPath(t *testing.T) {
	_, _, err := ParsePatchString("=3")

	if err == nil {
		t.Error("Expected error for empty path, but got none")
	}

	if err.Error() != "Empty path is not allowed" {
		t.Errorf("Expected error: 'Empty path is not allowed', got: %s", err.Error())
	}
}

func TestParsePatchString_EmptyValue(t *testing.T) {
	_, _, err := ParsePatchString("spec.replicas=")

	if err == nil {
		t.Error("Expected error for empty value, but got none")
	}

	if err.Error() != "Empty value is not allowed" {
		t.Errorf("Expected error: 'Empty value is not allowed', got: %s", err.Error())
	}
}

// Tests for ParsePatch (returns domain object)
func TestParsePatch_ValidIntPatch(t *testing.T) {
	patch, err := ParsePatch("spec.replicas=3")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(patch.Path) != 2 {
		t.Errorf("Expected 2 path segments, got: %d", len(patch.Path))
	}

	if patch.Path[0] != "spec" || patch.Path[1] != "replicas" {
		t.Errorf("Expected ['spec', 'replicas'], got: %v", patch.Path)
	}

	if patch.Value.ToInterface() != int64(3) {
		t.Errorf("Expected value 3, got: %v", patch.Value.ToInterface())
	}
}

func TestParsePatch_ValidBoolPatch(t *testing.T) {
	patch, err := ParsePatch("spec.paused=true")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if patch.Value.ToInterface() != true {
		t.Errorf("Expected value true, got: %v", patch.Value.ToInterface())
	}
}

func TestParsePatch_StringValue(t *testing.T) {
	patch, err := ParsePatch("spec.image=nginx:latest")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if patch.Value.ToInterface() != "nginx:latest" {
		t.Errorf("Expected value 'nginx:latest', got: %v", patch.Value.ToInterface())
	}
}

// Tests for ParsePatches (multiple patches)
func TestParsePatches_MultiplePatchesValid(t *testing.T) {
	patchesRaw := []string{
		"spec.replicas=3",
		"spec.paused=true",
	}

	patches, err := ParsePatches(patchesRaw)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(patches) != 2 {
		t.Errorf("Expected 2 patches, got: %d", len(patches))
	}

	// Check first patch
	if patches[0].Path[0] != "spec" || patches[0].Path[1] != "replicas" {
		t.Errorf("Expected first patch path ['spec', 'replicas'], got: %v", patches[0].Path)
	}

	if patches[0].Value.ToInterface() != int64(3) {
		t.Errorf("Expected first patch value 3, got: %v", patches[0].Value.ToInterface())
	}

	// Check second patch
	if patches[1].Value.ToInterface() != true {
		t.Errorf("Expected second patch value true, got: %v", patches[1].Value.ToInterface())
	}
}

func TestParsePatches_EmptySlice(t *testing.T) {
	patchesRaw := []string{}

	_, err := ParsePatches(patchesRaw)

	if err == nil {
		t.Error("Expected error for empty patches slice, but got none")
	}

	if err.Error() != "No patches provided" {
		t.Errorf("Expected error: 'No patches provided', got: %s", err.Error())
	}
}

func TestParsePatches_OneInvalidPatch(t *testing.T) {
	patchesRaw := []string{
		"spec.replicas=3",
		"invalid",
		"spec.paused=true",
	}

	_, err := ParsePatches(patchesRaw)

	if err == nil {
		t.Error("Expected error for invalid patch, but got none")
	}
}
