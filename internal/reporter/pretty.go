package reporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/foomo/obacht/pkg/schema"
)

// severityOrder maps severity to a sort rank (lower = more severe).
var severityOrder = map[schema.Severity]int{
	schema.SeverityCritical: 0,
	schema.SeverityHigh:     1,
	schema.SeverityWarn:     2,
	schema.SeverityInfo:     3,
}

// PrettyReporter renders a human-friendly, coloured report using lipgloss.
type PrettyReporter struct {
	// ShowPassing controls whether checks with status=pass are rendered.
	// When false (default), passing checks are omitted from the per-check
	// listing; the summary line still reflects the full counts.
	ShowPassing bool
}

// NewPrettyReporter creates a new PrettyReporter.
func NewPrettyReporter() *PrettyReporter {
	return &PrettyReporter{}
}

// Report writes the pretty-printed scan results to w.
func (p *PrettyReporter) Report(w io.Writer, result *schema.ScanResult) error {
	// Styles.
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
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

		// Skip the category entirely if every check would be filtered out.
		if !p.ShowPassing {
			anyVisible := false

			for _, cr := range checks {
				if cr.Status != schema.StatusPass {
					anyVisible = true
					break
				}
			}

			if !anyVisible {
				continue
			}
		}

		// Category header.
		fmt.Fprintln(w, boldStyle.Render(cat))

		for _, cr := range checks {
			if cr.Status == schema.StatusPass && !p.ShowPassing {
				continue
			}

			var icon string

			switch cr.Status {
			case schema.StatusPass:
				icon = greenStyle.Render("\u2713")
			case schema.StatusFail:
				icon = SeverityColorStyle(cr.Severity).Render("\u2717")
			case schema.StatusSkip:
				icon = yellowStyle.Render("-")
			case schema.StatusError:
				icon = SeverityColorStyle(cr.Severity).Render("!")
			}

			if cr.Status == schema.StatusFail || cr.Status == schema.StatusError {
				sevStyle := SeverityColorStyle(cr.Severity)
				sevLabel := sevStyle.Render("[" + string(cr.Severity) + "]")
				fmt.Fprintf(w, "  %s %s %s: %s\n", icon, cr.RuleID, sevLabel, cr.Title)
			} else {
				fmt.Fprintf(w, "  %s %s: %s\n", icon, cr.RuleID, cr.Title)
			}

			// For failures and errors, show evidence and remediation.
			if cr.Status == schema.StatusFail || cr.Status == schema.StatusError {
				if cr.Evidence != "" {
					parts := splitEvidence(cr.Evidence)
					switch len(parts) {
					case 0:
						// All-whitespace evidence — render nothing.
					case 1:
						fmt.Fprintf(w, "      Evidence: %s\n", parts[0])
					default:
						fmt.Fprintln(w, "      Evidence:")

						for _, p := range parts {
							fmt.Fprintf(w, "        - %s\n", p)
						}
					}
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
	critStyle := SeverityColorStyle(schema.SeverityCritical)
	highStyle := SeverityColorStyle(schema.SeverityHigh)
	warnStyle := SeverityColorStyle(schema.SeverityWarn)
	infoStyle := SeverityColorStyle(schema.SeverityInfo)
	fmt.Fprintf(w, "Summary: %d failed, %d passed, %d skipped (%s, %s, %s, %s)\n",
		s.Failed, s.Passed, s.Skipped,
		critStyle.Render(fmt.Sprintf("%d critical", s.Critical)),
		highStyle.Render(fmt.Sprintf("%d high", s.High)),
		warnStyle.Render(fmt.Sprintf("%d warn", s.Warn)),
		infoStyle.Render(fmt.Sprintf("%d info", s.Info)),
	)

	return nil
}

// splitEvidence splits an aggregated evidence string into individual findings.
// The engine joins per-finding evidence with "; " (see pkg/engine/engine.go).
// Returns whitespace-trimmed, non-empty parts; nil for empty input.
func splitEvidence(s string) []string {
	if s == "" {
		return nil
	}

	raw := strings.Split(s, "; ")
	parts := make([]string, 0, len(raw))

	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}

	return parts
}

// SeverityColorStyle returns the lipgloss style for the given severity level.
func SeverityColorStyle(s schema.Severity) lipgloss.Style {
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
