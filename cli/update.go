package cli

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

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
		case ERROR:
			return m.updateError(msg)
		}
	case errMsg:
		return m.updateError(msg)
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
			return m.errorCreate(err)
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

		hours, err := strconv.Atoi(m.inputs[1].Value()) // convert hours string to int
		if err != nil {
			return m.errorCreate(err)
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
			return m.errorCreate(err)
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

func (m model) updateError(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			m.page = PROJECTS
			return m, nil
		}
	}

	return m, nil
}

func (m model) errorCreate(err error) (model, tea.Cmd) {
	m.err = err
	m.page = ERROR
	return m, func() tea.Msg { return errMsg(err) }
}
