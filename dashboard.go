package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	backends    []*Backend
	totalReqs   int64
	rateLimited int64
	input       textinput.Model
	addingMode  bool
}

func (m model) Init() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.addingMode {
			if msg.String() == "enter" {
				addBackend(m.input.Value())
				m.input.Reset()
				m.addingMode = false
			} else {
				m.input, _ = m.input.Update(msg)
			}
		} else {
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "a":
				m.addingMode = true
				m.input.Focus()
			}
		}
	}
	m.backends = listBackends()
	return m, m.Init()
}
func (m model) View() string {
	var s strings.Builder
	s.WriteString("=== GoRelay Dashboard ===\n\n")
	s.WriteString(fmt.Sprintf("%-20s  %-8s  %-10s  %s\n", "Backend", "Status", "Requests", "Latency"))
	s.WriteString(strings.Repeat("-", 55) + "\n")
	for _, b := range m.backends {
		status := "✓ UP"
		if !b.Healthy {
			status = "✗ DOWN"
		}
		fmt.Fprintf(&s, "%-20s  %-8s  %-10d  %v\n",
			b.Address, status, b.Requests.Load(), b.Latency)
	}
	fmt.Fprintf(&s, "\nTotal Requests: %d  |  Rate Limited: %d\n",
		totalRequests.Load(), totalRateLimited.Load())
	if m.addingMode {
		s.WriteString("\nEnter backend address: " + m.input.View())
	} else {
		s.WriteString("\n[Q] Quit  [A] Add Backend\n")
	}
	return s.String()
}
