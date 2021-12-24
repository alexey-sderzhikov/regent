package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alexey-sderzhikov/regent/restapi"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	PROJECTS         = "projects"
	ISSUES           = "issues"
	INPUT_TIME_ENTRY = "input_time_entry"
)

type errMsg error

type model struct {
	redmineClient *restapi.RmClient
	projects      []restapi.Project
	issues        []restapi.Issue
	inputs        []textinput.Model
	focusIndex    int
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

	inputs := make([]textinput.Model, 3)
	inputs[0] = initialCommentInput()
	inputs[1] = initialDateInput()
	inputs[2] = initialHoursInput()

	return model{
		redmineClient: rc,
		projects:      projects.Projects,
		page:          PROJECTS,
		objectCount:   len(projects.Projects),
		inputs:        inputs,
		focusIndex:    0,
	}, nil
}

func initialHoursInput() textinput.Model {
	ti := textinput.NewModel()
	ti.Placeholder = "Work hours"
	ti.CharLimit = 5
	ti.Width = 10

	return ti
}

func initialDateInput() textinput.Model {
	ti := textinput.NewModel()
	ti.Placeholder = "Date format - 2020-12-25"
	ti.CharLimit = 12
	ti.Width = 12

	return ti
}

func initialCommentInput() textinput.Model {
	ti := textinput.NewModel()
	ti.Placeholder = "Some comment"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 20

	return ti
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
	case INPUT_TIME_ENTRY:
		return m.viewInputTimeEntry()
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

func (m model) viewInputTimeEntry() string {
	return fmt.Sprint(
		"\nText comment to time entry:\n",
		m.inputs[0].View(),
		"\nText date:\n",
		m.inputs[1].View(),
		"\nText work hours:\n",
		m.inputs[2].View(),
		"\n",
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.page {
		case PROJECTS:
			return m.updateProjects(msg)
		case ISSUES:
			return m.updateIssues(msg)
		case INPUT_TIME_ENTRY:
			return m.updateInputTimeEntry(msg)
		}
	}

	return m, nil
}
func (m model) navigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down":
		if m.cursor < m.objectCount-1 {
			m.cursor++
		}
	}

	return m, nil
}
func (m model) updateProjects(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	default:
		return m.navigation(msg)
	}

	return m, nil
}

func (m model) updateIssues(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		m.page = INPUT_TIME_ENTRY
	default:
		return m.navigation(msg)
	}

	return m, nil
}

func (m model) updateInputTimeEntry(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "up":
		if m.focusIndex > 0 {
			m.inputs[m.focusIndex].Blur()
			m.focusIndex--
			m.inputs[m.focusIndex].Focus()
		}
	case "down":
		if m.focusIndex < len(m.inputs)-1 {
			m.inputs[m.focusIndex].Blur()
			m.focusIndex++
			m.inputs[m.focusIndex].Focus()
		}
	case "enter":
		issue := m.issues[m.cursor]

		hours, err := strconv.Atoi(m.inputs[2].Value()) // convert hours string to int
		if err != nil {
			return m, nil
		}

		// err = m.redmineClient.CreateTimeEntry(
		// 	issue.Id,
		// 	m.inputs[2].Value(), // input date
		// 	m.inputs[0].Value(), // input comment
		// 	hours,
		// )
		fmt.Print(issue.Id,
			m.inputs[2].Value(), // input date
			m.inputs[0].Value(), // input comment
			hours)
		if err != nil {
			fmt.Print(err)
			return m, tea.Quit
		}

		fmt.Print(m.inputs[m.focusIndex].Value())
	case "esc":
		return m, tea.Quit
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for ind := 0; ind <= len(m.inputs)-1; ind++ {
		m.inputs[ind], cmd = m.inputs[ind].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
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
