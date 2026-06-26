package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/512uk/pingu/internal/ping"
)

const maxHistory = 512

// PingResultMsg wraps a ping.Result as a Bubble Tea message.
type PingResultMsg ping.Result

// Model is the Bubble Tea model for the ping monitor TUI.
type Model struct {
	target  string
	ch      <-chan ping.Result
	history []ping.Result
	width   int
	height  int
}

// New creates a TUI model wired to the given ping result channel.
func New(target string, ch <-chan ping.Result) Model {
	return Model{
		target: target,
		ch:     ch,
	}
}

// waitForPing bridges the ping channel into Bubble Tea's message loop.
// This is the standard pattern from Bubble Tea's realtime examples —
// return a Cmd that blocks on the channel, then delivers the value as a Msg.
func waitForPing(ch <-chan ping.Result) tea.Cmd {
	return func() tea.Msg {
		r, ok := <-ch
		if !ok {
			return nil // channel closed
		}
		return PingResultMsg(r)
	}
}

func (m Model) Init() tea.Cmd {
	return waitForPing(m.ch)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case PingResultMsg:
		m.history = append(m.history, ping.Result(msg))
		if len(m.history) > maxHistory {
			m.history = m.history[len(m.history)-maxHistory:]
		}
		// Re-subscribe to the channel for the next result.
		return m, waitForPing(m.ch)
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "waiting for terminal size..."
	}

	compact := m.width < 25 || m.height < 18

	if len(m.history) == 0 {
		content := "waiting for first ping..."
		box := boxStyle.
			Width(m.width - 2).
			Height(m.height - 2).
			Align(lipgloss.Center, lipgloss.Center).
			Render(content)
		return box
	}

	last := m.history[len(m.history)-1]
	ms := float64(last.RTT) / float64(time.Millisecond)
	col := ColorForRTT(ms)

	// Build the RTT display.
	var rttText string
	if last.Err != nil {
		rttText = rttStyle.Foreground(red).Render("timeout")
	} else {
		rttText = rttStyle.Foreground(col).Render(fmt.Sprintf("%dms", last.RTT.Milliseconds()))
	}

	// Build sparkline.
	sparkWidth := m.width - 4 // padding for box borders
	if sparkWidth < 1 {
		sparkWidth = 1
	}
	durations := make([]time.Duration, len(m.history))
	for i, r := range m.history {
		durations[i] = r.RTT
	}
	spark := sparkStyle.Foreground(col).Render(Render(durations, sparkWidth))

	// Compute stats.
	minRTT, maxRTT, avgRTT, lossPct := stats(m.history)

	if compact {
		return m.compactView(rttText, spark, minRTT, maxRTT, avgRTT, lossPct)
	}
	return m.expandedView(rttText, spark, minRTT, maxRTT, avgRTT, lossPct)
}

func (m Model) compactView(rtt, spark string, min, max, avg time.Duration, loss float64) string {
	minC := ColorForRTT(float64(min) / float64(time.Millisecond))
	maxC := ColorForRTT(float64(max) / float64(time.Millisecond))
	avgC := ColorForRTT(float64(avg) / float64(time.Millisecond))
	lossC := ColorForRTT(loss * 3) // scale so any loss looks bad

	var b strings.Builder
	b.WriteString(rtt + "\n")
	b.WriteString("\n")
	b.WriteString(spark + "\n")
	b.WriteString("\n")
	b.WriteString(
		labelStyle.Render("mn:") + valueStyle.Foreground(minC).Render(fmt.Sprintf("%-4d", min.Milliseconds())) +
			" " +
			labelStyle.Render("mx:") + valueStyle.Foreground(maxC).Render(fmt.Sprintf("%d", max.Milliseconds())) + "\n")
	b.WriteString(
		labelStyle.Render("av:") + valueStyle.Foreground(avgC).Render(fmt.Sprintf("%-4d", avg.Milliseconds())) +
			" " +
			labelStyle.Render("ls:") + valueStyle.Foreground(lossC).Render(fmt.Sprintf("%.0f%%", loss)))

	content := b.String()

	box := boxStyle.
		Width(m.width - 2).
		Height(m.height - 2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	// Inject target into top border.
	return injectHeader(box, m.target, m.width)
}

func (m Model) expandedView(rtt, spark string, min, max, avg time.Duration, loss float64) string {
	minC := ColorForRTT(float64(min) / float64(time.Millisecond))
	maxC := ColorForRTT(float64(max) / float64(time.Millisecond))
	avgC := ColorForRTT(float64(avg) / float64(time.Millisecond))
	lossC := ColorForRTT(loss * 3)

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(rtt + "\n")
	b.WriteString("\n")
	b.WriteString(spark + "\n")
	b.WriteString("\n")
	b.WriteString(
		labelStyle.Render("min  ") + valueStyle.Foreground(minC).Render(fmt.Sprintf("%-8s", fmt.Sprintf("%dms", min.Milliseconds()))) +
			labelStyle.Render("max  ") + valueStyle.Foreground(maxC).Render(fmt.Sprintf("%dms", max.Milliseconds())) + "\n")
	b.WriteString(
		labelStyle.Render("avg  ") + valueStyle.Foreground(avgC).Render(fmt.Sprintf("%-8s", fmt.Sprintf("%dms", avg.Milliseconds()))) +
			labelStyle.Render("loss ") + valueStyle.Foreground(lossC).Render(fmt.Sprintf("%.0f%%", loss)))

	content := b.String()

	header := "pingu ~ " + m.target
	box := boxStyle.
		Width(m.width - 2).
		Height(m.height - 2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	return injectHeader(box, header, m.width)
}

// injectHeader overwrites part of the top border with a label.
func injectHeader(box, label string, width int) string {
	lines := strings.Split(box, "\n")
	if len(lines) == 0 {
		return box
	}

	styled := " " + headerStyle.Render(label) + " "
	top := lines[0]
	topRunes := []rune(top)

	// Place the label after the corner + one border char.
	// We need to be careful with ANSI escape sequences in the border,
	// so we use a simple prefix/suffix approach.
	if len(topRunes) > 4 {
		// Find a safe insertion point (after "╭─").
		prefix := string(topRunes[:2])
		suffix := string(topRunes[2:])
		// Trim enough border chars to make room for the label.
		labelLen := len([]rune(label)) + 2 // +2 for spaces
		if labelLen < len([]rune(suffix))-1 {
			suffix = string([]rune(suffix)[labelLen:])
		}
		lines[0] = prefix + styled + suffix
	}

	return strings.Join(lines, "\n")
}

// stats computes min, max, avg RTT and loss percentage over a history slice.
func stats(history []ping.Result) (min, max, avg time.Duration, lossPct float64) {
	if len(history) == 0 {
		return 0, 0, 0, 0
	}

	var sum time.Duration
	var count int
	var losses int

	for _, r := range history {
		if r.Err != nil || r.RTT <= 0 {
			losses++
			continue
		}
		if count == 0 || r.RTT < min {
			min = r.RTT
		}
		if r.RTT > max {
			max = r.RTT
		}
		sum += r.RTT
		count++
	}

	if count > 0 {
		avg = sum / time.Duration(count)
	}
	lossPct = float64(losses) / float64(len(history)) * 100

	return min, max, avg, lossPct
}
