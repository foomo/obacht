package cli

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/foomo/obacht/pkg/engine"
	"github.com/foomo/obacht/pkg/schema"
)

const (
	categoryWidth = 12
	barWidth      = 40
)

const heroLogo = ` ┌─┐┌┐ ┌─┐┌─┐┬ ┬┌┬┐
 │ │├┴┐├─┤│  ├─┤ │
 └─┘└─┘┴ ┴└─┘┴ ┴ ┴ `

const heroTagline = "developer environment scanner"

// categoryState tracks the status of a single category during the scan.
type categoryState struct {
	category string
	status   string // "pending", "running", "done"
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
	categories  []categoryState
	catIndex    map[string]int
	bar         progress.Model
	spin        spinner.Model
	groupsDone  int
	groupsTotal int
	result      *schema.ScanResult
	err         error
	done        bool
	ctx         context.Context //nolint:containedctx
	rules       []schema.RulesFile
	program     *tea.Program
}

func newScanModel(ctx context.Context, ruleFiles []schema.RulesFile) *scanModel {
	var cats []categoryState

	catIndex := make(map[string]int)

	for _, rf := range ruleFiles {
		for _, r := range rf.Rules {
			if _, exists := catIndex[r.Category]; !exists {
				catIndex[r.Category] = len(cats)
				cats = append(cats, categoryState{
					category: r.Category,
					status:   "pending",
				})
			}
		}
	}

	bar := progress.New(
		progress.WithDefaultBlend(),
		progress.WithWidth(barWidth),
		progress.WithoutPercentage(),
	)

	sp := spinner.New(spinner.WithSpinner(spinner.Dot))
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	return &scanModel{
		categories: cats,
		catIndex:   catIndex,
		bar:        bar,
		spin:       sp,
		ctx:        ctx,
		rules:      ruleFiles,
	}
}

// SetProgram stores a reference to the tea.Program so the engine callback can send messages.
func (m *scanModel) SetProgram(p *tea.Program) {
	m.program = p
}

func (m *scanModel) Init() tea.Cmd {
	return tea.Batch(m.runScan(), m.spin.Tick)
}

func (m *scanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case progressMsg:
		evt := engine.ProgressEvent(msg)

		if m.groupsTotal == 0 {
			m.groupsTotal = evt.GroupTotal
		}

		idx, ok := m.catIndex[evt.Category]
		if !ok {
			break
		}

		switch evt.Kind {
		case engine.EventGroupStart:
			m.categories[idx].status = "running"

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
			m.groupsDone++

			pct := 0.0
			if m.groupsTotal > 0 {
				pct = float64(m.groupsDone) / float64(m.groupsTotal)
			}

			return m, m.bar.SetPercent(pct)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)

		return m, cmd

	case progress.FrameMsg:
		newBar, cmd := m.bar.Update(msg)
		m.bar = newBar

		if m.done && m.err == nil && !m.bar.IsAnimating() {
			return m, tea.Quit
		}

		return m, cmd

	case scanDoneMsg:
		m.done = true
		m.result = msg.result
		m.err = msg.err

		if msg.err != nil {
			return m, tea.Quit
		}

		cmd := m.bar.SetPercent(1.0)
		if !m.bar.IsAnimating() {
			return m, tea.Quit
		}

		return m, cmd
	}

	return m, nil
}

func (m *scanModel) View() tea.View {
	var b strings.Builder

	heroStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	taglineStyle := lipgloss.NewStyle().Faint(true).PaddingLeft(1)
	labelStyle := lipgloss.NewStyle().Width(categoryWidth)
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	dimStyle := lipgloss.NewStyle().Faint(true)
	countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	b.WriteString(heroStyle.Render(heroLogo) + "\n")
	b.WriteString(taglineStyle.Render(heroTagline) + "\n\n")

	count := fmt.Sprintf("%d/%d", m.groupsDone, m.groupsTotal)
	if m.groupsTotal == 0 {
		count = "starting…"
	}

	fmt.Fprintf(&b, " %s  %s\n\n", m.bar.View(), countStyle.Render(count))

	for _, cat := range m.categories {
		label := labelStyle.Render(cat.category)

		switch cat.status {
		case "pending":
			fmt.Fprintf(&b, "  %s %s\n", dimStyle.Render("·"), dimStyle.Render(label))

		case "running":
			fmt.Fprintf(&b, "  %s %s\n", m.spin.View(), label)

		case "done":
			icon := greenStyle.Render("✓")
			if cat.failed > 0 || cat.errors > 0 {
				icon = redStyle.Render("✗")
			}

			summary := formatCategorySummary(cat)
			fmt.Fprintf(&b, "  %s %s %s\n", icon, label, dimStyle.Render(summary))
		}
	}

	if m.done {
		b.WriteString("\n")
	}

	return tea.NewView(b.String())
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
