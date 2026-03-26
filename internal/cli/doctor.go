package cli

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/franklinkim/bouncer/internal/collector"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check bouncer dependencies and configuration",
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
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	fmt.Println(boldStyle.Render("Bouncer Doctor"))
	fmt.Println(boldStyle.Render("=============="))
	fmt.Println()

	// --- OPA Engine ---
	fmt.Println(boldStyle.Render("OPA Engine"))
	fmt.Printf("  Status: %s embedded\n", greenStyle.Render("\u2713"))
	fmt.Println()

	// --- Policies ---
	fmt.Println(boldStyle.Render("Policies"))

	rules, err := loadEmbeddedRules()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error loading rules: %v\n", err)
	}

	regoFiles, err := loadEmbeddedRego()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error loading rego: %v\n", err)
	}

	if rulesDir != "" {
		extRules, extRego, err := loadExternalRules(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error loading external rules: %v\n", err)
		} else {
			rules = mergeRules(rules, extRules)

			regoFiles = append(regoFiles, extRego...)
		}
	}

	fmt.Printf("  Rules:  %d loaded\n", len(rules))
	fmt.Printf("  Rego:   %d files\n", len(regoFiles))

	if len(rules) > 0 {
		ids := make([]string, len(rules))
		for i, r := range rules {
			ids[i] = r.ID
		}

		fmt.Printf("  IDs:    %s\n", strings.Join(ids, ", "))
	}

	fmt.Println()

	// --- Collectors ---
	fmt.Println(boldStyle.Render("Collectors"))

	collectors := []collector.Collector{
		collector.NewSSHCollector(),
		collector.NewGitCollector(),
		collector.NewDockerCollector(),
		collector.NewKubeCollector(),
		collector.NewEnvCollector(),
		collector.NewShellCollector(),
		collector.NewToolsCollector(),
		collector.NewPathCollector(),
		collector.NewOSCollector(),
	}

	_, results, err := collector.CollectAll(ctx, collectors)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error running collectors: %v\n", err)
	} else {
		for _, r := range results {
			var icon string

			switch r.Status {
			case collector.StatusOK:
				icon = greenStyle.Render("\u2713")
			case collector.StatusSkipped:
				icon = yellowStyle.Render("-")
			case collector.StatusError:
				icon = redStyle.Render("\u2717")
			}

			statusStr := string(r.Status)
			if r.Status == collector.StatusError && r.Error != nil {
				statusStr = fmt.Sprintf("error: %v", r.Error)
			}

			fmt.Printf("  %s %-8s %s\n", icon, r.Name, statusStr)
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
