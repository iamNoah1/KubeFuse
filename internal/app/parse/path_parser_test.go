package parse

import "testing"

func TestPathParser_ThreeValidSegmentsShouldReturnCorrectSlice(t *testing.T) {
	testString := "spec.template.metadata"

	slice, err := ParsePath(testString)

	if err != nil {
		t.Errorf("Error while parsing")
	}

	if len(slice) != 3 {
		t.Errorf("Unexpected amount of elements")
	}

	if slice[0] != "spec" || slice[1] != "template" || slice[2] != "metadata" {
		t.Errorf("Unexpected elements")
	}
}

func TestPathParser_ThreeSegmentsMiddleOneEmptyShouldReturnError(t *testing.T) {
	testString := "spec..metadata"

	_, err := ParsePath(testString)

	if err != nil {
		if err.Error() != "Empty segment not allowed" {
			t.Errorf("Expected error: %s, got: %s", "Empty segment not allowed", err.Error())
		}
	} else {
		t.Errorf("Expected error")
	}
}

func TestPathParser_FirstSegmentEmptyShouldReturnError(t *testing.T) {
	testString := ".metadata"

	_, err := ParsePath(testString)

	if err != nil {
		if err.Error() != "Empty segment not allowed" {
			t.Errorf("Expected error: %s, got: %s", "Empty segment not allowed", err.Error())
		}
	} else {
		t.Errorf("Expected error")
	}
}

func TestPathParser_LastSegmentEmptyShouldReturnError(t *testing.T) {
	testString := "metadata."

	_, err := ParsePath(testString)

	if err != nil {
		if err.Error() != "Empty segment not allowed" {
			t.Errorf("Expected error: %s, got: %s", "Empty segment not allowed", err.Error())
		}
	} else {
		t.Errorf("Expected error")
	}
}

func TestPathParser_EmptyStringShouldReturnError(t *testing.T) {
	testString := ""

	_, err := ParsePath(testString)

	if err != nil {
		if err.Error() != "Empty path is not allowed" {
			t.Errorf("Expected error: %s, got: %s", "Empty path is not allowed", err.Error())
		}
	} else {
		t.Errorf("Expected error")
	}
}
