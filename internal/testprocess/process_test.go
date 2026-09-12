package testprocess

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name                    string
		timeout                 time.Duration
		wantError, wantDeadline bool
	}{
		{"success", 5 * time.Second, false, false},
		{"failure", 5 * time.Second, true, false},
		{"panic", 5 * time.Second, true, false},
		{"hang", 100 * time.Millisecond, true, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := run("TestSubprocessFixture", tt.timeout, tt.name)
			if (err != nil) != tt.wantError {
				t.Fatalf("run = %v, want error %v", err, tt.wantError)
			}
			if errors.Is(err, context.DeadlineExceeded) != tt.wantDeadline {
				t.Errorf("run = %v, want deadline %v", err, tt.wantDeadline)
			}
		})
	}
}

func TestSubprocessFixture(t *testing.T) {
	switch os.Getenv("PORTUGOL_TEST_MODE") {
	case "failure":
		t.Fatal("intentional child failure")
	case "panic":
		panic("intentional child panic")
	case "hang":
		<-time.After(time.Minute)
	}
}

func TestRunIsolates(t *testing.T) {
	Run(t, func() {
		if os.Getenv("PORTUGOL_TEST_CHILD") != t.Name() {
			t.Fatal("body did not run in a child process")
		}
	})
}
