package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

	"github.com/foomo/obacht/internal/runner"
	"github.com/foomo/obacht/pkg/schema"
	"github.com/foomo/obacht/rules"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check obacht dependencies and configuration",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

// bumblebeeVersion returns the first whitespace-separated token from
// `bumblebee version`, or "(version unknown)" on any failure.
func bumblebeeVersion(ctx context.Context, path string) string {
	out, err := exec.CommandContext(ctx, path, "version").Output()
	if err != nil {
		return "(version unknown)"
	}

	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return "(version unknown)"
	}

	return fields[0]
}

func runDoctor(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	boldStyle := lipgloss.NewStyle().Bold(true)
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	fmt.Println(boldStyle.Render("obacht Doctor"))
	fmt.Println(boldStyle.Render("=============="))
	fmt.Println()

	// --- OPA Engine ---
	fmt.Println(boldStyle.Render("OPA Engine"))
	fmt.Printf("  Status: %s embedded\n", greenStyle.Render("\u2713"))
	fmt.Println()

	// --- Policies ---
	fmt.Println(boldStyle.Render("Policies"))

	var (
		ruleFiles []schema.RulesFile
		err       error
	)

	if rulesDir != "" {
		ruleFiles, err = loadExternalRuleFiles(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error loading external rules: %v\n", err)
		}
	} else {
		ruleFiles, err = loadEmbeddedRuleFiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error loading rules: %v\n", err)
		}
	}

	var (
		totalRules int
		ids        []string
	)

	for _, rf := range ruleFiles {
		totalRules += len(rf.Rules)
		for _, r := range rf.Rules {
			ids = append(ids, r.ID)
		}
	}

	fmt.Printf("  Rules:  %d loaded\n", totalRules)
	fmt.Printf("  Files:  %d rule files\n", len(ruleFiles))

	if len(ids) > 0 {
		fmt.Printf("  IDs:    %s\n", strings.Join(ids, ", "))
	}

	fmt.Println()

	// --- Dependencies ---
	fmt.Println(boldStyle.Render("Dependencies"))

	if _, err := exec.LookPath("jq"); err == nil {
		fmt.Printf("  %s jq found on PATH\n", greenStyle.Render("\u2713"))
	} else {
		fmt.Printf("  %s jq not found on PATH (required by rule scripts; install via `brew install jq`)\n", redStyle.Render("\u2717"))
	}

	if path, err := exec.LookPath("bumblebee"); err == nil {
		version := bumblebeeVersion(ctx, path)
		fmt.Printf("  %s bumblebee %s\n", greenStyle.Render("\u2713"), version)
		fmt.Println("    bumblebee rules (BUM*) enabled")
	} else {
		fmt.Printf("  %s bumblebee not found on PATH (bumblebee rules BUM* will skip)\n", redStyle.Render("\u2717"))
		fmt.Println("    install: go install github.com/perplexityai/bumblebee/cmd/bumblebee@latest")
	}

	fmt.Println()

	// --- Input Scripts ---
	fmt.Println(boldStyle.Render("Input Scripts"))

	var fsys fs.FS
	if rulesDir != "" {
		fsys = os.DirFS(rulesDir)
	} else {
		fsys = rules.Embedded
	}

	for _, rf := range ruleFiles {
		for _, r := range rf.Rules {
			if r.Input == "" {
				continue
			}

			scriptPath := fmt.Sprintf("inputs/%s/%s.sh", r.Category, r.ID)
			result := runner.RunInputForRule(ctx, fsys, scriptPath, r.Input, r.ID)

			var icon string

			switch result.Status {
			case runner.StatusOK:
				icon = greenStyle.Render("\u2713")
				fmt.Printf("  %s %s ok\n", icon, r.ID)
			case runner.StatusSkipped:
				icon = greenStyle.Render("\u2713")
				fmt.Printf("  %s %s skip: %s\n", icon, r.ID, result.SkipReason)
			case runner.StatusError:
				icon = redStyle.Render("\u2717")
				fmt.Printf("  %s %s error: %v\n", icon, r.ID, result.Error)
			}
		}
	}

	fmt.Println()

	// --- System ---
	fmt.Println(boldStyle.Render("System"))
	fmt.Printf("  OS:      %s\n", runtime.GOOS)
	fmt.Printf("  Arch:    %s\n", runtime.GOARCH)
	fmt.Printf("  Shell:   %s\n", os.Getenv("SHELL"))
	fmt.Printf("  Go:      %s\n", runtime.Version())

	return nil
}
