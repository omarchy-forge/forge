package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/omarchy-forge/forge/internal/scaffold"
)

type agentRuntime interface {
	DefaultAgent(context.Context) (string, error)
	CommitBaseline(string, string) (string, error)
	Launch(string, io.Reader, io.Writer, io.Writer) error
}

type execAgentRuntime struct{}

func (execAgentRuntime) DefaultAgent(ctx context.Context) (string, error) {
	command := exec.CommandContext(ctx, "omarchy-default-agent")
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (execAgentRuntime) CommitBaseline(directory, name string) (string, error) {
	commands := []struct {
		label string
		args  []string
	}{
		{label: "add", args: []string{"-C", directory, "add", "--all"}},
		{label: "commit", args: []string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-C", directory, "commit", "-m", "chore: scaffold " + name + " with Omarchy Forge"}},
	}
	for _, operation := range commands {
		command := exec.Command("git", operation.args...)
		if output, err := command.CombinedOutput(); err != nil {
			return "", fmt.Errorf("git %s: %w: %s", operation.label, err, strings.TrimSpace(string(output)))
		}
	}
	command := exec.Command("git", "-C", directory, "rev-parse", "--short", "HEAD")
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("read baseline commit: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (execAgentRuntime) Launch(directory string, stdin io.Reader, stdout, stderr io.Writer) error {
	prompt, err := os.ReadFile(filepath.Join(directory, "AGENT_PROMPT.md"))
	if err != nil {
		return fmt.Errorf("read initial agent prompt: %w", err)
	}
	command := exec.Command("omarchy", "agent", "prompt", string(prompt))
	command.Dir = directory
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("launch configured Omarchy agent: %w", err)
	}
	return nil
}

func configuredAgent(runtime agentRuntime) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	agent, err := runtime.DefaultAgent(ctx)
	if err != nil {
		return "", fmt.Errorf("read the configured Omarchy agent: %w", err)
	}
	if agent == "" {
		return "", fmt.Errorf("no Omarchy coding agent is configured; choose one with: omarchy default agent <name>")
	}
	return agent, nil
}

func promptAgentPlan(reader *bufio.Reader, stdout io.Writer, options *scaffold.Options, agent string) bool {
	for {
		options.Name = prompt(reader, stdout, "Plugin name", options.Name)
		options.ID = prompt(reader, stdout, "Plugin ID (for example my.clock)", options.ID)
		options.Description = prompt(reader, stdout, "What outcome should this plugin provide? (for example: monitor local development projects and surface problems)", options.Description)
		options.Author = prompt(reader, stdout, "Author", options.Author)
		mode := strings.ToLower(prompt(reader, stdout, "How should Forge define the plugin? Type references or questionnaire", "references"))
		switch mode {
		case "references", "reference", "r":
			options.ReferenceDriven = true
			if !promptReferencePlan(reader, stdout, options) {
				fmt.Fprintln(stdout, "Reference setup was not completed; nothing was written or launched.")
				return false
			}
		case "questionnaire", "questions", "q":
			options.ReferenceDriven = false
			promptQuestionnairePlan(reader, stdout, options)
		default:
			fmt.Fprintln(stdout, "Please answer references or questionnaire.")
			continue
		}
		options.Section = prompt(reader, stdout, "Default section", options.Section)
		if !options.ReferenceDriven && options.References == nil {
			references, err := promptReferences(reader, stdout, *options, false)
			if err != nil {
				fmt.Fprintf(stdout, "Cannot prepare references: %v\n", err)
				continue
			}
			options.References = references
		}

		fmt.Fprintln(stdout, "\nReview the agent project:")
		fmt.Fprintf(stdout, "  Agent: %s\n  Directory: %s\n  Name: %s\n  Plugin ID: %s\n  Description: %s\n", agent, options.Directory, options.Name, options.ID, options.Description)
		fmt.Fprintf(stdout, "  Bar: %s\n  Click: %s\n  Dashboard: %s\n  Actions: %s\n  Data: %s\n  Commands: %s\n  Network: %s\n  Persistence: %s\n  Failures: %s\n  Section: %s\n", options.BarSummary, options.ClickBehavior, options.PopoutSummary, options.UserActions, options.DataSources, options.LocalCommands, options.NetworkAccess, options.Persistence, options.FailureBehavior, options.Section)
		if len(options.References) == 0 {
			fmt.Fprintln(stdout, "  Reference files: none")
		} else {
			fmt.Fprintln(stdout, "  Reference files:")
			for _, reference := range options.References {
				fmt.Fprintf(stdout, "    %s -> %s (%s, %d bytes, sha256 %s)\n", reference.SourcePath, reference.ProjectPath, reference.Kind, reference.Size, reference.SHA256[:12])
			}
		}
		fmt.Fprintln(stdout, "  Safety: Omarchy launches its configured agent with unattended permissions. Generated instructions constrain behavior but are not a sandbox.")
		fmt.Fprintln(stdout, "  Privacy: Forge does not upload references, but the configured coding agent/provider may receive their contents. Do not attach secrets or personal data.")
		answer := strings.ToLower(prompt(reader, stdout, "Is everything correct? Create the project, commit its baseline, and launch "+agent+"? Type yes, edit, or cancel", "yes"))
		switch answer {
		case "yes", "y":
			return true
		case "edit", "e":
			fmt.Fprintln(stdout, "Revisiting the answers. Press Enter to keep each current value.")
			changeReferences := strings.ToLower(prompt(reader, stdout, "Replace the selected reference files? Type yes or no", "no"))
			if changeReferences == "yes" || changeReferences == "y" {
				options.References = nil
			}
		case "cancel", "c", "no", "n":
			fmt.Fprintln(stdout, "Cancelled; nothing was written or launched.")
			return false
		default:
			fmt.Fprintln(stdout, "Please answer yes, edit, or cancel.")
		}
	}
}

func promptQuestionnairePlan(reader *bufio.Reader, stdout io.Writer, options *scaffold.Options) {
	fmt.Fprintln(stdout, "Design tip: keep the bar entry simple—usually one icon plus a short status. Put detailed information and controls in the popout dashboard.")
	options.BarSummary = prompt(reader, stdout, "What minimal bar entry should appear? (for example: one icon plus healthy/warning status)", defaultValue(options.BarSummary, "One recognizable icon with a concise status indicator."))
	options.ClickBehavior = prompt(reader, stdout, "What should happen when the user clicks it?", defaultValue(options.ClickBehavior, "Left-click toggles the popout; right-click refreshes."))
	options.PopoutSummary = prompt(reader, stdout, "List the exact dashboard cards or details to show (for example: Git branch, test status, disk space, running servers)", options.PopoutSummary)
	options.UserActions = prompt(reader, stdout, "List the exact buttons or actions users need (for example: refresh, run static tests, open project; or none)", defaultValue(options.UserActions, "Refresh the displayed data."))
	options.DataSources = prompt(reader, stdout, "For each dashboard item, what local file, command output, API, or in-memory value supplies it?", options.DataSources)
	options.LocalCommands = prompt(reader, stdout, "List each exact local command/process and its purpose; answer none only when all requested data is truly in memory", defaultValue(options.LocalCommands, "Not required — none."))
	options.NetworkAccess = prompt(reader, stdout, "What network access does it need?", defaultValue(options.NetworkAccess, "Not required — none."))
	options.Persistence = prompt(reader, stdout, "What data should it save, and where?", defaultValue(options.Persistence, "Not required — none."))
	options.FailureBehavior = prompt(reader, stdout, "How should failures appear?", defaultValue(options.FailureBehavior, "Show a clear error state and preserve the last known safe data."))
}

func promptReferencePlan(reader *bufio.Reader, stdout io.Writer, options *scaffold.Options) bool {
	options.BarSummary = "Derive every visible bar item and status from the confirmed reference files; keep the bar entry minimal unless a reference explicitly requires otherwise."
	options.ClickBehavior = "Derive every click, keyboard, and interaction behavior from the confirmed reference files."
	options.PopoutSummary = "Implement every panel, section, displayed value, status, and control visible or described in the confirmed reference files."
	options.UserActions = "Implement every user action visible or described in the confirmed reference files; do not reduce controls to decoration."
	options.DataSources = "Derive the exact source for every displayed value from the confirmed reference files and the approved access boundaries below."
	options.LocalCommands = defaultValue(options.LocalCommands, "Not permitted until the user confirms the local-access boundary below.")
	options.NetworkAccess = defaultValue(options.NetworkAccess, "Not required — none.")
	options.Persistence = defaultValue(options.Persistence, "Not required — none.")
	options.FailureBehavior = defaultValue(options.FailureBehavior, "Show a clear error state and preserve the last known safe data.")
	if options.References == nil {
		references, err := promptReferences(reader, stdout, *options, true)
		if err != nil {
			fmt.Fprintf(stdout, "Cannot prepare references: %v\n", err)
			return false
		}
		options.References = references
	}
	if len(options.References) == 0 {
		fmt.Fprintln(stdout, "Reference mode requires at least one text or image file. Add a reference or choose questionnaire.")
		options.References = nil
		return false
	}
	localAccess := strings.ToLower(prompt(reader, stdout, "May the plugin read local project files and run bounded local commands required by the references? Type yes or no", "yes"))
	if localAccess == "yes" || localAccess == "y" {
		options.LocalCommands = "Permitted only for functionality required by the confirmed references. Use fixed programs, array-form arguments, least privilege, and timeouts no longer than 10 seconds; document every exact command."
	} else {
		options.LocalCommands = "Not permitted — do not read project files or run local commands or processes."
	}
	options.NetworkAccess = prompt(reader, stdout, "What network access may the plugin use?", defaultValue(options.NetworkAccess, "Not required — none."))
	options.Persistence = prompt(reader, stdout, "What data may the plugin save, and where?", defaultValue(options.Persistence, "Not required — none."))
	options.FailureBehavior = prompt(reader, stdout, "How should failures appear?", defaultValue(options.FailureBehavior, "Show a clear error state and preserve the last known safe data."))
	return true
}

func promptReferences(reader *bufio.Reader, stdout io.Writer, options scaffold.Options, required bool) ([]scaffold.Reference, error) {
	directory, err := scaffold.CreateReferenceDirectory(options)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(stdout, "\nOptional references directory created: %s\n", directory)
	fmt.Fprintln(stdout, "Drop up to 10 .txt, .md, .png, .jpg, .jpeg, .webp, or .svg files there.")
	fmt.Fprintln(stdout, "Forge will not upload them, but the configured agent/provider may receive their contents. Do not include secrets or personal data.")
	prompt(reader, stdout, "Press Enter after the file copies have finished", "")
	for {
		references, prepareErr := scaffold.PrepareReferencesInDirectory(directory)
		if prepareErr == nil && len(references) > 0 {
			fmt.Fprintf(stdout, "Found %d reference file(s) in %s.\n", len(references), directory)
			return references, nil
		}
		if prepareErr == nil && !required {
			return nil, nil
		}
		if prepareErr != nil {
			fmt.Fprintf(stdout, "Forge could not use the files currently in %s: %v\n", directory, prepareErr)
		} else {
			fmt.Fprintf(stdout, "Forge found no supported reference files in %s.\n", directory)
		}
		fmt.Fprint(stdout, "Add or finish copying files, then press Enter to scan again; type cancel to stop: ")
		answer, readErr := reader.ReadString('\n')
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer == "cancel" || answer == "c" || (readErr != nil && answer == "") {
			fmt.Fprintf(stdout, "Reference selection stopped; any files already added remain in %s.\n", directory)
			return nil, nil
		}
	}
}

func defaultValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
