package cli

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

	"github.com/foomo/obacht/internal/runner"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check obacht dependencies and configuration",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
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

	ruleFiles, err := loadEmbeddedRuleFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error loading rules: %v\n", err)
	}

	if rulesDir != "" {
		extRuleFiles, err := loadExternalRuleFiles(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error loading external rules: %v\n", err)
		} else {
			ruleFiles = mergeRuleFiles(ruleFiles, extRuleFiles)
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

	// --- Input Scripts ---
	fmt.Println(boldStyle.Render("Input Scripts"))

	// Deduplicate input scripts and test each one.
	seen := make(map[string]bool)

	for _, rf := range ruleFiles {
		inputs := []string{rf.Input}
		for _, r := range rf.Rules {
			if r.Input != "" {
				inputs = append(inputs, r.Input)
			}
		}

		for _, input := range inputs {
			if input == "" || seen[input] {
				continue
			}

			seen[input] = true

			result := runner.RunInput(ctx, input)

			var icon string

			switch result.Status {
			case runner.StatusOK:
				icon = greenStyle.Render("\u2713")
				fmt.Printf("  %s input script ok\n", icon)
			case runner.StatusError:
				icon = redStyle.Render("\u2717")
				fmt.Printf("  %s input script error: %v\n", icon, result.Error)
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
