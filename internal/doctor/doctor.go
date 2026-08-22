// Package doctor performs read-only diagnostics against the local toolchain.
package doctor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/omarchy-forge/forge/internal/checks"
)

type Status string

const (
	Pass Status = "pass"
	Warn Status = "warn"
	Fail Status = "fail"
)

type Diagnostic struct {
	RuleID      string
	Status      Status
	Message     string
	Remediation string
}
type Report struct {
	Diagnostics []Diagnostic
	Project     checks.Report
}

type Runner interface {
	LookPath(string) (string, error)
	Run(context.Context, string, ...string) ([]byte, error)
}
type ExecRunner struct{}

func (ExecRunner) LookPath(name string) (string, error) { return exec.LookPath(name) }
func (ExecRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, name, args...)
	return command.CombinedOutput()
}

func Run(target string, runner Runner) Report {
	result := Report{Diagnostics: []Diagnostic{}, Project: checks.Run(target)}
	add := func(id string, status Status, message, fix string) {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{id, status, message, fix})
	}
	if _, err := runner.LookPath("omarchy"); err != nil {
		add("OD100", Fail, "Omarchy command was not found", "Install Omarchy or run doctor on an Omarchy system.")
	} else {
		add("OD100", Pass, "Omarchy command is available", "")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		output, err := runner.Run(ctx, "pacman", "-Q", "omarchy")
		cancel()
		if err != nil {
			add("OD101", Warn, "Omarchy was found but its package version could not be determined", "Confirm the installation with the system package manager.")
		} else {
			add("OD101", Pass, "Detected "+strings.TrimSpace(string(output)), "")
		}
	}
	if _, err := runner.LookPath("quickshell"); err != nil {
		add("OD102", Fail, "Quickshell was not found", "Install the Quickshell version required by Omarchy.")
	} else {
		add("OD102", Pass, "Quickshell is available", "")
	}
	if _, err := runner.LookPath("omarchy-shell"); err != nil {
		add("OD103", Fail, "omarchy-shell was not found", "Repair the Omarchy shell installation.")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		output, err := runner.Run(ctx, "omarchy-shell", "shell", "ping")
		cancel()
		if err != nil || !bytes.Equal(bytes.TrimSpace(output), []byte("ok")) {
			add("OD103", Fail, "Omarchy Shell IPC did not answer ping", "Start or restart Omarchy Shell and inspect its logs.")
		} else {
			add("OD103", Pass, "Omarchy Shell IPC answered ping", "")
		}
	}
	if _, err := runner.LookPath("qmllint"); err != nil {
		add("OD104", Warn, "qmllint is unavailable; QML linting cannot be performed", "Install Qt QML tooling for additional local diagnostics.")
	} else {
		add("OD104", Pass, "qmllint is available", "")
	}
	if _, err := runner.LookPath("omarchy"); err != nil {
		add("OD105", Fail, "official Omarchy plugin validator was not found", "Update or repair Omarchy before validating plugins.")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		output, err := runner.Run(ctx, "omarchy", "plugin", "validate", target)
		cancel()
		if err != nil {
			message := "Official validator rejected the project"
			if detail := strings.TrimSpace(string(output)); detail != "" {
				message += ": " + detail
			}
			add("OD105", Fail, message, "Resolve the official validator error before installation.")
		} else {
			add("OD105", Pass, "Official Omarchy plugin validation passed", "")
		}
	}
	return result
}

func WriteText(w interface{ Write([]byte) (int, error) }, report Report) error {
	if _, err := fmt.Fprintln(w, "Omarchy Forge doctor"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	for _, d := range report.Diagnostics {
		if _, err := fmt.Fprintf(w, "%-4s  %s  %s\n", strings.ToUpper(string(d.Status)), d.RuleID, d.Message); err != nil {
			return err
		}
		if d.Remediation != "" {
			if _, err := fmt.Fprintf(w, "      Fix: %s\n", d.Remediation); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintln(w, "\nProject checks:"); err != nil {
		return err
	}
	return checks.WriteText(w, report.Project)
}

func (r Report) Failed() bool {
	for _, d := range r.Diagnostics {
		if d.Status == Fail {
			return true
		}
	}
	return r.Project.Summary.Errors > 0
}
