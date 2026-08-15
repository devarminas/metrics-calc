package main

import (
	"errors"
	"strings"
	"testing"
)

type failingReader struct{ done bool }

func (f *failingReader) Read(p []byte) (int, error) {
	if f.done {
		return 0, errors.New("device on fire")
	}
	f.done = true
	return copy(p, "3 4\n"), nil
}

func TestRun_validLine(t *testing.T) {
	var out, errOut strings.Builder

	if err := run(strings.NewReader("3 4\n"), &out, &errOut); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	if got, want := out.String(), "12.00 14.00\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if got := errOut.String(); got != "" {
		t.Errorf("stderr = %q, want empty", got)
	}
}

func TestRun_invalidLinesReportedAndSkipped(t *testing.T) {
	input := "3 4\nx 4\n2.5 0.5\n1 2 3\n5\n-3 4\n3 0\n10 10\n"
	var out, errOut strings.Builder

	if err := run(strings.NewReader(input), &out, &errOut); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	want := "12.00 14.00\n1.25 6.00\n100.00 40.00\n"
	if got := out.String(); got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}

	stderrLines := strings.Split(strings.TrimSuffix(errOut.String(), "\n"), "\n")
	if len(stderrLines) != 5 {
		t.Fatalf("stderr produced %d lines, want 5:\n%s", len(stderrLines), errOut.String())
	}
	for i, prefix := range []string{"line 2:", "line 4:", "line 5:", "line 6:", "line 7:"} {
		if !strings.HasPrefix(stderrLines[i], prefix) {
			t.Errorf("stderr line %d = %q, want prefix %q", i+1, stderrLines[i], prefix)
		}
	}
}

func TestRun_toleratesSurroundingWhitespace(t *testing.T) {
	var out, errOut strings.Builder

	if err := run(strings.NewReader("  3 \t  4  \n"), &out, &errOut); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	if got, want := out.String(), "12.00 14.00\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if got := errOut.String(); got != "" {
		t.Errorf("stderr = %q, want empty", got)
	}
}

func TestRun_overlongLineDoesNotAbort(t *testing.T) {
	input := strings.Repeat("1 ", 50_000) + "\n3 4\n"
	var out, errOut strings.Builder

	if err := run(strings.NewReader(input), &out, &errOut); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}

	if got, want := out.String(), "12.00 14.00\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if got := strings.Count(errOut.String(), "\n"); got != 1 {
		t.Errorf("stderr produced %d lines, want 1:\n%s", got, errOut.String())
	}
	if !strings.HasPrefix(errOut.String(), "line 1:") {
		t.Errorf("stderr = %q, want a message for line 1", errOut.String())
	}
}

func TestRun_readErrorIsReturned(t *testing.T) {
	var out, errOut strings.Builder

	err := run(&failingReader{}, &out, &errOut)

	if err == nil {
		t.Fatal("run() error = nil, want a read error")
	}
	if !strings.Contains(err.Error(), "device on fire") {
		t.Errorf("run() error = %v, want it to wrap the reader's error", err)
	}
}
