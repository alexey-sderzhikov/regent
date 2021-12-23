package cli

import (
	"fmt"
	"os"

	"github.com/alexey-sderzhikov/regent/restapi"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	PROJECTS = "projects"
	ISSUES   = "issues"
)

type model struct {
	redmineClient *restapi.RmClient
	projects      []restapi.Project
	issues        []restapi.Issue
	objectCount   int
	cursor        int
	page          string
}

func initialModel() (model, error) {
	rc, err := restapi.NewRm("", "")
	if err != nil {
		return model{}, err
	}

	projects, err := rc.GetProjects()
	if err != nil {
		return model{}, err
	}

	chs := make([]string, len(projects.Projects))
	for ind, proj := range projects.Projects {
		chs[ind] = proj.Name
	}

	return model{
		redmineClient: rc,
		projects:      projects.Projects,
		page:          PROJECTS,
		objectCount:   len(projects.Projects),
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	switch m.page {
	case PROJECTS:
		return m.viewProjects()
	case ISSUES:
		return m.viewIssues()
	}

	return "Cannot detect current page :("
}

func (m model) viewProjects() string {
	s := "\nBergen Projects\n"

	for ind, p := range m.projects {
		cursor := " "
		if m.cursor == ind {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, p.Name)
	}

	s += "\nPress q to quit.\n"

	return s
}

func (m model) viewIssues() string {
	s := fmt.Sprintf("\nIssues project's - %q\n", m.issues[0].Project.Name)

	for ind, i := range m.issues {
		cursor := " "
		if m.cursor == ind {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, i.Subject)
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
			if m.cursor < m.objectCount-1 {
				m.cursor++
			}

		case "enter", " ":
			switch m.page {
			case PROJECTS:
				return m.updateProjects(msg)
			case ISSUES:
				return m.updateIssues(msg)
			}
		}
	}

	return m, nil
}

func (m model) updateProjects(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			project := m.projects[m.cursor]

			issues, err := m.redmineClient.GetIssues(project.Id)
			if err != nil {
				fmt.Print(err)
				return m, tea.Quit
			}

			m.issues = issues.Issues
			m.page = ISSUES
			m.cursor = 0
			m.objectCount = len(m.issues)
		}
	}

	return m, nil
}

func (m model) updateIssues(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			issue := m.issues[m.cursor]

			err := m.redmineClient.CreateTimeEntry(issue.Id, "2021-12-23")
			if err != nil {
				fmt.Print(err)
				return m, tea.Quit
			}
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
