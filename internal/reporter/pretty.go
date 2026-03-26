package reporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/franklinkim/bouncer/pkg/schema"
)

// severityOrder maps severity to a sort rank (lower = more severe).
var severityOrder = map[schema.Severity]int{
	schema.SeverityCritical: 0,
	schema.SeverityHigh:     1,
	schema.SeverityWarn:     2,
	schema.SeverityInfo:     3,
}

// PrettyReporter renders a human-friendly, coloured report using lipgloss.
type PrettyReporter struct{}

// NewPrettyReporter creates a new PrettyReporter.
func NewPrettyReporter() *PrettyReporter {
	return &PrettyReporter{}
}

// Report writes the pretty-printed scan results to w.
func (p *PrettyReporter) Report(w io.Writer, result *schema.ScanResult) error {
	// Styles.
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	boldStyle := lipgloss.NewStyle().Bold(true)

	// Group results by category.
	groups := make(map[string][]schema.CheckResult)
	var categoryOrder []string
	for _, cr := range result.Results {
		cat := cr.Category
		if _, exists := groups[cat]; !exists {
			categoryOrder = append(categoryOrder, cat)
		}
		groups[cat] = append(groups[cat], cr)
	}
	sort.Strings(categoryOrder)

	// Sort each group by severity (critical first).
	for _, checks := range groups {
		sort.Slice(checks, func(i, j int) bool {
			return severityOrder[checks[i].Severity] < severityOrder[checks[j].Severity]
		})
	}

	// Render categories.
	for _, cat := range categoryOrder {
		checks := groups[cat]

		// Category header.
		fmt.Fprintln(w, boldStyle.Render(cat))

		for _, cr := range checks {
			var icon string
			switch cr.Status {
			case schema.StatusPass:
				icon = greenStyle.Render("\u2713")
			case schema.StatusFail:
				icon = redStyle.Render("\u2717")
			case schema.StatusSkip:
				icon = yellowStyle.Render("-")
			case schema.StatusError:
				icon = redStyle.Render("!")
			}

			fmt.Fprintf(w, "  %s %s: %s\n", icon, cr.RuleID, cr.Title)

			// For failures and errors, show evidence and remediation.
			if cr.Status == schema.StatusFail || cr.Status == schema.StatusError {
				if cr.Evidence != "" {
					fmt.Fprintf(w, "      Evidence: %s\n", cr.Evidence)
				}
				if cr.Remediation != "" {
					fmt.Fprintf(w, "      Fix: %s\n", cr.Remediation)
				}
			}
		}
		fmt.Fprintln(w)
	}

	// Separator.
	fmt.Fprintln(w, strings.Repeat("-", 60))

	// Summary line.
	s := result.Summary
	fmt.Fprintf(w, "Summary: %d failed, %d passed, %d skipped (%d critical, %d high, %d warn, %d info)\n",
		s.Failed, s.Passed, s.Skipped,
		s.Critical, s.High, s.Warn, s.Info,
	)

	return nil
}
