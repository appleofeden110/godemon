package shell

import (
	"testing"
)

// LINUX ONLY, use `pgrep godemon` on your machine first and change expVal 1 and 2
// Conclusion: Successful
func TestGetPIDs(t *testing.T) {
	ints, err := GetPIDs("godemon")
	if err != nil {
		t.Fatalf("Error in GetPIDs func: %v\n", err)
	}
	expVal_1 := 4292
	trueVal_1 := ints[0]

	if expVal_1 != trueVal_1 {
		t.Fatalf("Values don't match. expected: %v, got: %v", expVal_1, trueVal_1)
	}

	expVal_2 := 4345
	trueVal_2 := ints[1]

	if expVal_2 != trueVal_2 {
		t.Fatalf("Values don't match. expected: %v, got: %v", expVal_2, trueVal_2)
	}
}
