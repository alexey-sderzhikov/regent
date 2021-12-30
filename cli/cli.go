package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func Start() error {
	m, err := initialModel()
	if err != nil {
		return err
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Alias, there's been an error: %v", err)
		os.Exit(1)
		return nil
	} else {
		return err
	}
}
