package logger

import "testing"

func TestHelloName(t *testing.T) {
	LogLocation("test")
	if GetLogLocation() != "test" {
		t.Errorf(`GetLogLocation() = %q, want test`, GetLogLocation())
	}
}
