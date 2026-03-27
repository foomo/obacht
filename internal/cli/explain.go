package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/franklinkim/bouncer/internal/reporter"
	"github.com/franklinkim/bouncer/pkg/schema"
)

var explainCmd = &cobra.Command{
	Use:   "explain <rule-id>",
	Short: "Show detailed information about a rule",
	Args:  cobra.ExactArgs(1),
	RunE:  runExplain,
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

func runExplain(cmd *cobra.Command, args []string) error {
	ruleID := args[0]

	// Load built-in rules from embedded policies.
	ruleFiles, err := loadEmbeddedRuleFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading embedded rules: %v\n", err)
		os.Exit(Error)
	}

	// Optionally load and merge external rules from --rules-dir.
	if rulesDir != "" {
		extRuleFiles, err := loadExternalRuleFiles(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "loading external rules: %v\n", err)
			os.Exit(Error)
		}

		ruleFiles = mergeRuleFiles(ruleFiles, extRuleFiles)
	}

	// Find the rule matching the given ID (case-insensitive).
	var found *schema.Rule

	for _, rf := range ruleFiles {
		for i := range rf.Rules {
			if strings.EqualFold(rf.Rules[i].ID, ruleID) {
				found = &rf.Rules[i]
				break
			}
		}

		if found != nil {
			break
		}
	}

	if found == nil {
		fmt.Fprintf(os.Stderr, "rule %q not found\n", ruleID)
		os.Exit(Error)
	}

	printRule(found)

	return nil
}

func printRule(r *schema.Rule) {
	boldStyle := lipgloss.NewStyle().Bold(true)
	severityStyle := reporter.SeverityColorStyle(r.Severity)

	fmt.Printf("%s %s\n", boldStyle.Render("Rule:"), r.ID)
	fmt.Printf("%s %s\n", boldStyle.Render("Title:"), r.Title)
	fmt.Printf("%s %s\n", boldStyle.Render("Severity:"), severityStyle.Render(string(r.Severity)))
	fmt.Printf("%s %s\n", boldStyle.Render("Category:"), r.Category)

	if r.Description != "" {
		fmt.Printf("\n%s\n", boldStyle.Render("Description:"))

		for _, line := range wrapLines(r.Description) {
			fmt.Printf("  %s\n", line)
		}
	}

	if r.Remediation != "" {
		fmt.Printf("\n%s\n", boldStyle.Render("Remediation:"))

		for _, line := range wrapLines(r.Remediation) {
			fmt.Printf("  %s\n", line)
		}
	}
}

// wrapLines splits text into trimmed, non-empty lines.
func wrapLines(text string) []string {
	var lines []string

	for line := range strings.SplitSeq(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}

	return lines
}
