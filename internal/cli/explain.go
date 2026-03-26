package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

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
	rules, err := loadEmbeddedRules()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading embedded rules: %v\n", err)
		os.Exit(Error)
	}

	// Optionally load and merge external rules from --rules-dir.
	if rulesDir != "" {
		extRules, _, err := loadExternalRules(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "loading external rules: %v\n", err)
			os.Exit(Error)
		}
		rules = mergeRules(rules, extRules)
	}

	// Find the rule matching the given ID (case-insensitive).
	var found *schema.Rule
	for i := range rules {
		if strings.EqualFold(rules[i].ID, ruleID) {
			found = &rules[i]
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
	severityStyle := severityColorStyle(r.Severity)

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

func severityColorStyle(s schema.Severity) lipgloss.Style {
	switch s {
	case schema.SeverityCritical:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	case schema.SeverityHigh:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	case schema.SeverityWarn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	case schema.SeverityInfo:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	default:
		return lipgloss.NewStyle()
	}
}

// wrapLines splits text into trimmed, non-empty lines.
func wrapLines(text string) []string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}
