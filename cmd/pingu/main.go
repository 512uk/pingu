package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/512uk/pingu/internal/ping"
	"github.com/512uk/pingu/internal/tui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	target := "1.1.1.1"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	// Context lets us cleanly shut down the pinger goroutine when the TUI exits.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pinger := ping.New(target, 1*time.Second)
	ch := pinger.Start(ctx)

	model := tui.New(target, ch)
	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err := p.Run()
	return err
}
