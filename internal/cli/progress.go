package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/franklinkim/bouncer/pkg/engine"
	"github.com/franklinkim/bouncer/pkg/schema"
)

const categoryWidth = 12 // left-column width for category labels
const barWidth = 30      // width of each progress bar

// categoryState tracks the status of a single category during the scan.
type categoryState struct {
	category string
	status   string // "pending", "running", "done"
	bar      progress.Model
	passed   int
	failed   int
	skipped  int
	errors   int
}

// progressMsg wraps an engine.ProgressEvent as a bubbletea message.
type progressMsg engine.ProgressEvent

// scanDoneMsg signals that the scan has finished.
type scanDoneMsg struct {
	result *schema.ScanResult
	err    error
}

// scanModel is the bubbletea model for the scan progress display.
type scanModel struct {
	categories []categoryState
	catIndex   map[string]int // category name -> index in categories slice
	result     *schema.ScanResult
	err        error
	done       bool
	ctx        context.Context
	rules      []schema.RulesFile
	program    *tea.Program
}

func newProgressBar() progress.Model {
	return progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(barWidth),
		progress.WithoutPercentage(),
	)
}

func newScanModel(ctx context.Context, ruleFiles []schema.RulesFile) *scanModel {
	// Collect unique categories in order from rule files.
	var cats []categoryState
	catIndex := make(map[string]int)

	for _, rf := range ruleFiles {
		for _, r := range rf.Rules {
			if _, exists := catIndex[r.Category]; !exists {
				catIndex[r.Category] = len(cats)
				cats = append(cats, categoryState{
					category: r.Category,
					status:   "pending",
					bar:      newProgressBar(),
				})
			}
		}
	}

	return &scanModel{
		categories: cats,
		catIndex:   catIndex,
		ctx:        ctx,
		rules:      ruleFiles,
	}
}

// SetProgram stores a reference to the tea.Program so the engine callback can send messages.
func (m *scanModel) SetProgram(p *tea.Program) {
	m.program = p
}

func (m *scanModel) Init() tea.Cmd {
	return m.runScan()
}

func (m *scanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case progressMsg:
		evt := engine.ProgressEvent(msg)
		idx, ok := m.catIndex[evt.Category]
		if !ok {
			break
		}

		switch evt.Kind {
		case engine.EventGroupStart:
			m.categories[idx].status = "running"
			// Animate the bar to 100%.
			cmd := m.categories[idx].bar.SetPercent(1.0)
			return m, cmd

		case engine.EventGroupDone:
			cat := &m.categories[idx]
			for _, r := range evt.Results {
				switch r.Status {
				case schema.StatusPass:
					cat.passed++
				case schema.StatusFail:
					cat.failed++
				case schema.StatusSkip:
					cat.skipped++
				case schema.StatusError:
					cat.errors++
				}
			}
			cat.status = "done"
		}

	case scanDoneMsg:
		m.done = true
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit

	case progress.FrameMsg:
		// Forward frame messages to all running progress bars.
		var cmds []tea.Cmd
		for i := range m.categories {
			if m.categories[i].status == "running" {
				model, cmd := m.categories[i].bar.Update(msg)
				m.categories[i].bar = model.(progress.Model)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *scanModel) View() string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().Width(categoryWidth)
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	dimStyle := lipgloss.NewStyle().Faint(true)

	for _, cat := range m.categories {
		label := labelStyle.Render(cat.category)

		switch cat.status {
		case "pending":
			bar := cat.bar.ViewAs(0)
			fmt.Fprintf(&b, "  %s %s\n", label, dimStyle.Render(bar))

		case "running":
			bar := cat.bar.View()
			fmt.Fprintf(&b, "  %s %s\n", label, bar)

		case "done":
			bar := cat.bar.ViewAs(1.0)
			summary := formatCategorySummary(cat)

			icon := greenStyle.Render("✓")
			if cat.failed > 0 || cat.errors > 0 {
				icon = redStyle.Render("✗")
			}

			fmt.Fprintf(&b, "  %s %s %s %s\n", label, bar, icon, dimStyle.Render(summary))
		}
	}

	if m.done {
		b.WriteString("\n")
	}

	return b.String()
}

// runScan returns a tea.Cmd that runs the engine evaluation in a goroutine.
func (m *scanModel) runScan() tea.Cmd {
	return func() tea.Msg {
		callback := func(evt engine.ProgressEvent) {
			if m.program != nil {
				m.program.Send(progressMsg(evt))
			}
		}

		result, err := engine.Evaluate(m.ctx, m.rules, callback)

		return scanDoneMsg{result: result, err: err}
	}
}

// formatCategorySummary builds a short summary string like "3 passed, 1 failed".
func formatCategorySummary(cat categoryState) string {
	var parts []string

	if cat.passed > 0 {
		parts = append(parts, fmt.Sprintf("%d passed", cat.passed))
	}

	if cat.failed > 0 {
		parts = append(parts, fmt.Sprintf("%d failed", cat.failed))
	}

	if cat.skipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", cat.skipped))
	}

	if cat.errors > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", cat.errors))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ")
}
