package cli

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.crumbs.getCurrentPage() {
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
	switch msg.Type {
	case tea.KeyEscape:
		return m, tea.Quit
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < m.objectCount-1 {
			m.cursor++
		}
	}

	return m, nil
}

func (m model) updateProjects(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter: // go to project issues
		if len(m.issues) == 0 {
			project := m.projects[m.cursor]

			issues, err := m.redmineClient.GetIssues(project.Id)
			if err != nil {
				return m.errorCreate(err)
			}

			m.issues = issues.Issues
			m.objectCount = len(m.issues)
		}

		m.crumbs = m.crumbs.addPage(ISSUES)
		m.cursor = 0
	default:
		return m.navigation(msg)
	}

	return m, nil
}

func (m model) updateIssues(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter: // go to creation new time entry for issue
		m.crumbs = m.crumbs.addPage(INPUT_TIME_ENTRY)
	case tea.KeyBackspace: // go to previos page
		m.cursor = 0
		m.crumbs = m.crumbs.popPage() // go to previos page
	default:
		return m.navigation(msg)
	}

	return m, nil
}

func (m model) updateInputTimeEntry(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyUp: // go to upstair input field
		if m.focusIndex > 0 {
			m.inputs[m.focusIndex].Blur()
			m.focusIndex--
			m.inputs[m.focusIndex].Focus()
		}
	case tea.KeyDown: //go to downstair input field
		if m.focusIndex < len(m.inputs)-1 {
			m.inputs[m.focusIndex].Blur()
			m.focusIndex++
			m.inputs[m.focusIndex].Focus()
		}
	case tea.KeyEnter: // create time entire
		issue := m.issues[m.cursor]

		hours, err := strconv.Atoi(m.inputs[2].Value()) // convert hours string to int
		if err != nil {
			return m.errorCreate(err)
		}

		err = m.redmineClient.CreateTimeEntry(
			issue.Id,
			m.inputs[1].Value(), // input date
			m.inputs[0].Value(), // input comment
			hours,
		)
		if err != nil {
			return m.errorCreate(err)
		}
	case tea.KeyBackspace: // go to previos page
		m.cursor = 0
		m.crumbs = m.crumbs.popPage()
	case tea.KeyEsc: // escape programm
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
		switch msg.Type {
		case tea.KeyEnter:
			m.cursor = 0
			m.crumbs = m.crumbs.popPage()
			return m, nil
		}
	}

	return m, nil
}

func (m model) errorCreate(err error) (model, tea.Cmd) {
	m.err = err
	m.crumbs = m.crumbs.addPage(ERROR)
	return m, func() tea.Msg { return errMsg(err) }
}
