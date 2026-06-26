package tui

import "github.com/charmbracelet/lipgloss"

// ANSI 256 colours for broad terminal compatibility.
var (
	green  = lipgloss.Color("34")  // good: <50ms
	yellow = lipgloss.Color("214") // ok: 50-150ms
	red    = lipgloss.Color("196") // bad: >150ms
	dim    = lipgloss.Color("240") // labels, borders
	white  = lipgloss.Color("255")
)

var (

	// Border style for the outer box.
	boxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(dim)

	// Large RTT display in the center.
	rttStyle = lipgloss.NewStyle().
			Bold(true)

	// Stat labels (min, max, avg, loss).
	labelStyle = lipgloss.NewStyle().
			Foreground(dim)

	// Stat values.
	valueStyle = lipgloss.NewStyle().
			Foreground(white)

	// The sparkline itself.
	sparkStyle = lipgloss.NewStyle().
			Foreground(green)

	// Target in the header.
	headerStyle = lipgloss.NewStyle().
			Foreground(white).
			Bold(true)
)

// ColorForRTT returns the appropriate colour based on ping latency thresholds.
func ColorForRTT(ms float64) lipgloss.Color {

	switch {
	case ms < 50:
		return green
	case ms <= 150:
		return yellow
	default:
		return red
	}
}
