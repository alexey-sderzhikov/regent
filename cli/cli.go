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
	var f *os.File
	f, err = tea.LogToFile("logs", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(m, tea.WithInputTTY())
	if err := p.Start(); err != nil {
		fmt.Printf("Alias, there's been an error: %v", err)
		os.Exit(1)
		return nil
	} else {
		return err
	}
}
