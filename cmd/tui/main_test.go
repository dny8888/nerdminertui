package main

import (
	"os"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain_StartsAndQuits_NoGoroutineLeak(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("github.com/charmbracelet/bubbletea.Tick.func1"))

	// Save original args and stdin
	oldArgs := os.Args
	oldStdin := os.Stdin
	defer func() {
		os.Args = oldArgs
		os.Stdin = oldStdin
	}()

	// Mock arguments to start with mock mining and no store
	os.Args = []string{"nerdtui", "--mock", "--no-store"}
	os.Setenv("NM_MOCK_MINING", "true")
	defer os.Unsetenv("NM_MOCK_MINING")

	// Mock stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r

	// Write 'q' after a small delay to let UI initialize
	go func() {
		time.Sleep(500 * time.Millisecond)
		w.Write([]byte("q"))
		w.Close()
	}()

	// Since main() does not call os.Exit on success, it will simply return when the program finishes.
	main()
}
