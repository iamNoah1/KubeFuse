package parse

import (
	"testing"
)

// Tests for ParseTarget (the low-level string parser)
func TestParseTarget_TwoSegmentsShouldReturnCorrectSlice(t *testing.T) {
	testString := "deploy/web"

	target, err := ParseTarget(testString)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if target[0] != "deploy" || target[1] != "web" {
		t.Errorf("Expected kind: %s, but got: %s and name: %s, but got: %s", "deploy", target[0], "web", target[1])
	}
}

func TestParseTarget_ThreeSegmentsShouldReturnError(t *testing.T) {
	testString := "deploy/web/pod"

	_, err := ParseTarget(testString)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if err.Error() != "More than 2 target segments are not allowed" {
		t.Errorf("Expected error: %s, got: %s", "More than 2 target segments are not allowed", err.Error())
	}
}

func TestParseTarget_TwoSegmentsFirstEmptyShouldReturnError(t *testing.T) {
	testString := "/web"

	_, err := ParseTarget(testString)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if err.Error() != "Empty target segments are not allowed" {
		t.Errorf("Expected error: %s, got: %s", "Empty target segments are not allowed", err.Error())
	}
}

func TestParseTarget_TwoSegmentsSecondEmptyShouldReturnError(t *testing.T) {
	testString := "deploy/"

	_, err := ParseTarget(testString)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if err.Error() != "Empty target segments are not allowed" {
		t.Errorf("Expected error: %s, got: %s", "Empty target segments are not allowed", err.Error())
	}
}

func TestParseTarget_EmptyStringShouldReturnError(t *testing.T) {
	testString := ""

	_, err := ParseTarget(testString)

	if err == nil {
		t.Error("Expected error but got none")
	}

	if err.Error() != "Empty target is not allowed" {
		t.Errorf("Expected error: %s, got: %s", "Empty target is not allowed", err.Error())
	}
}

// Tests for ParseTargetToResourceRef (returns domain object)
func TestParseTargetToResourceRef_ValidTargetWithoutNamespace(t *testing.T) {
	testString := "deployment/nginx"

	ref, err := ParseTargetToResourceRef(testString, "")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if ref.Kind != "deployment" {
		t.Errorf("Expected kind 'deployment', got: %s", ref.Kind)
	}

	if ref.Name != "nginx" {
		t.Errorf("Expected name 'nginx', got: %s", ref.Name)
	}

	if ref.Namespace != "" {
		t.Errorf("Expected empty namespace, got: %s", ref.Namespace)
	}
}

func TestParseTargetToResourceRef_ValidTargetWithNamespace(t *testing.T) {
	testString := "service/api"
	namespace := "production"

	ref, err := ParseTargetToResourceRef(testString, namespace)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if ref.Kind != "service" {
		t.Errorf("Expected kind 'service', got: %s", ref.Kind)
	}

	if ref.Name != "api" {
		t.Errorf("Expected name 'api', got: %s", ref.Name)
	}

	if ref.Namespace != "production" {
		t.Errorf("Expected namespace 'production', got: %s", ref.Namespace)
	}
}

func TestParseTargetToResourceRef_InvalidTargetReturnsError(t *testing.T) {
	testString := "invalid"

	_, err := ParseTargetToResourceRef(testString, "")

	if err == nil {
		t.Error("Expected error for invalid target format, but got none")
	}
}

func TestParseTargetToResourceRef_EmptyTargetReturnsError(t *testing.T) {
	testString := ""

	_, err := ParseTargetToResourceRef(testString, "default")

	if err == nil {
		t.Error("Expected error for empty target, but got none")
	}

	if err.Error() != "Empty target is not allowed" {
		t.Errorf("Expected error: %s, got: %s", "Empty target is not allowed", err.Error())
	}
}
