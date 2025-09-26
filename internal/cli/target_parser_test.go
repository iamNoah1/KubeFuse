package cli

import "testing"

func TestTargetParser_TwoSegementsShouldReturnCorrectObject(t *testing.T) {
	testString := "deploy/web"

	target, err := ParseTarget(testString)

	if err != nil {
		t.Errorf("Error while parsing")
	}

	if target[0] != "deploy" || target[1] != "web" {
		t.Errorf("Expected kind: %s, but got: %s and name: %s, but got: %s", "deploy", target[0], "web", target[1])
	}
}

func TestTargetParser_ThreeSegmentsShouldReturnError(t *testing.T) {
	testString := "deploy/web/pod"

	_, err := ParseTarget(testString)

	if err != nil {
		if err.Error() != "More than 2 target segments are not allowed" {
			t.Errorf("Expected error: %s, got: %s", "More than 2 target segments are not allowed", err.Error())
		}
	}
}

func TestTargetParser_TwoSegmentsFirstEmptyShouldReturnError(t *testing.T) {
	testString := "/web"

	_, err := ParseTarget(testString)

	if err != nil {
		if err.Error() != "Empty target segments are not allowed" {
			t.Errorf("Expected error: %s, got: %s", "Empty target segments are not allowed", err.Error())
		}
	}
}

func TestTargetParser_TwoSegmentsSecondEmptyShouldReturnError(t *testing.T) {
	testString := "deploy/"

	_, err := ParseTarget(testString)

	if err != nil {
		if err.Error() != "Empty target segments are not allowed" {
			t.Errorf("Expected error: %s, got: %s", "Empty target segments are not allowed", err.Error())
		}
	}
}

func TestTargetParser_EmptyStringShouldReturnError(t *testing.T) {
	testString := ""

	_, err := ParseTarget(testString)

	if err != nil {
		if err.Error() != "Empty target is not allowed" {
			t.Errorf("Expected error: %s, got: %s", "Empty target is not allowed", err.Error())
		}
	}
}
