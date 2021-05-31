package acceptance

import "testing"

func NewTest(t *testing.T) {
	_ = New()

	if len(tests) != 3 {
		t.Fatalf("Wrong amount of created tests")
	}

	if runner == nil {
		t.Error("runner should be initialized")
	}
}
