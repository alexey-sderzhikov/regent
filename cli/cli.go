package cli

import (
	"fmt"
	"os"

	"github.com/alexey-sderzhikov/regent/restapi"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	projects restapi.ProjectList
	cursor   int
	selected map[int]struct{}
}

func initialModel() (model, error) {
	projects, err := restapi.GetProjects()
	if err != nil {
		return model{}, err
	}

	chs := make([]string, len(projects.Projects))
	for ind, proj := range projects.Projects {
		chs[ind] = proj.Name
	}

	return model{
		choices:  chs,
		projects: *projects,
		selected: make(map[int]struct{}),
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	s := "\nBergen Projects\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			project := m.projects.Projects[m.cursor]
			is, err := restapi.GetIssues(project.Id)

			m.choices = make([]string, len(is.Issues))
			for ind, issue := range is.Issues {
				m.choices[ind] = issue.Subject
			}

			if err != nil {
				fmt.Print(err)
			}
			fmt.Print(is)
		}
	}

	return m, nil
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
