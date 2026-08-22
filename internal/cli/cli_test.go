package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"-h"}} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(args, &stdout, &stderr, BuildInfo{})
			if code != 0 {
				t.Fatalf("Run() code = %d, want 0", code)
			}
			if !strings.Contains(stdout.String(), "Usage:") {
				t.Errorf("stdout missing usage: %q", stdout.String())
			}
			if stderr.Len() != 0 {
				t.Errorf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	info := BuildInfo{Version: "v1.2.3", Commit: "abc123", BuildDate: "2026-08-22"}
	code := Run([]string{"version"}, &stdout, &stderr, info)

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	want := "omaforge v1.2.3\ncommit: abc123\nbuilt: 2026-08-22\nmanifest schemas: 1\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunDevelopmentVersion(t *testing.T) {
	var stdout bytes.Buffer
	code := Run([]string{"version"}, &stdout, ioDiscard{}, BuildInfo{})

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	if !strings.HasPrefix(stdout.String(), "omaforge dev\n") {
		t.Errorf("stdout = %q, want development version", stdout.String())
	}
}

func TestRunErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "unknown command", args: []string{"doctor"}, want: `unknown command "doctor"`},
		{name: "version arguments", args: []string{"version", "extra"}, want: "does not accept arguments"},
		{name: "help arguments", args: []string{"--help", "extra"}, want: `unknown command "--help"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(test.args, &stdout, &stderr, BuildInfo{})
			if code == 0 {
				t.Fatal("Run() code = 0, want nonzero")
			}
			if stdout.Len() != 0 {
				t.Errorf("stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), test.want) {
				t.Errorf("stderr = %q, want substring %q", stderr.String(), test.want)
			}
		})
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
