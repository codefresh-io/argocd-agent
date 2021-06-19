package acceptance

import "testing"

func TestCreateTests(t *testing.T) {
	_ = New()

	if len(tests) != 4 {
		t.Fatalf("Wrong amount of created tests")
	}

	if runner == nil {
		t.Error("runner should be initialized")
	}
}
