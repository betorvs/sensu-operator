package utils

import "testing"

func TestContains(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	testString := "test"
	testResult := Contains(testSlice, testString)
	if !testResult {
		t.Fatalf("Invalid 1.1 TestContains %v", testResult)
	}

}

func TestRemove(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	testString := "test"
	testResult := Remove(testSlice, testString)
	if Contains(testResult, testString) {
		t.Fatalf("Invalid 2.1 TestRemove %v", testResult)
	}
}
