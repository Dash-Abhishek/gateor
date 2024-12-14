package internal

import "testing"

func TestSingleton(t *testing.T) {
	mux1 := InitializeMux()
	mux2 := InitializeMux()

	if mux1 != mux2 {
		t.Error("Expected mux1 and mux2 to be the same")
	}
}
