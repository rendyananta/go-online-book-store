package log

import "testing"

func TestSetUp(t *testing.T) {
	// make sure the package does not panic

	SetUp(Config{})
}
