package util

import "testing"

var _ = func() bool {
	testing.Init()
	return true
}()

func TestContains(t *testing.T) {
	arr := []string{"1", "2"}
	result := Contains(arr, "2")
	if !result {
		t.Error("Element should be found")
	}
}

func TestContainsFalse(t *testing.T) {
	arr := []string{"1", "2"}
	result := Contains(arr, "3")
	if result {
		t.Error("Element should be not found")
	}
}
