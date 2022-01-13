package cli

import (
	"strconv"
	"time"

	"github.com/alexey-sderzhikov/regent/restapi"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	PROJECTS         = "projects"
	ISSUES           = "issues"
	INPUT_TIME_ENTRY = "input_time_entry"
	ERROR            = "error"
	TIME_ENTRIES     = "time_entries"
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
		case TIME_ENTRIES:
			return m.updateTimeEntries(msg)
		case ERROR:
			return m.updateError(msg)
		}
	case errMsg:
		return m.updateError(msg)
	}

	return m, nil
}

// navigation logic base for most pages
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
	case tea.KeyCtrlQ: // go to previos page
		m.status = ""
		m.cursor = 0

		var err error
		m.crumbs, err = m.crumbs.popPage()
		if err != nil {
			return m.errorCreate(err)
		}
	}

	return m, nil
}

// update logic if key tap on "projects" page
func (m model) updateProjects(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter: // go to project issues
		var err error
		projectId := m.projects[m.cursor].Id

		params := make(restapi.Params, 0)
		params["project_id"] = projectId

		if m.filters.for_me {
			params["assigned_to_id"] = "me"
		}

		m.issues, err = m.redmineClient.GetIssues(params)
		if err != nil {
			return m.errorCreate(err)
		}

		m.objectCount = len(m.issues.Issues)

		m.cursor = 0
		m.crumbs = m.crumbs.addPage(ISSUES)
	default:
		return m.navigation(msg)
	}

	return m, nil
}

// TODO refactoring pagination switching
// update logic if key tap on "issues" page
func (m model) updateIssues(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter: // go to creation new time entry for issue
		m.inputs[1].SetValue(time.Now().Format("2006-01-02")) // set today date
		m.crumbs = m.crumbs.addPage(INPUT_TIME_ENTRY)
	case tea.KeyCtrlQ: // go to previos page
		m.cursor = 0
		m.crumbs, _ = m.crumbs.popPage()

		projects, err := m.redmineClient.GetProjects()
		if err != nil {
			return m.errorCreate(err)
		}

		m.projects = projects.Projects
		m.objectCount = len(m.projects)
	case tea.KeyCtrlA:
		params := make(restapi.Params, 0)

		params["limit"] = 10
		params["user_id"] = m.redmineClient.User.Id

		te, err := m.redmineClient.GetTimeEntryList(params)
		if err != nil {
			m.errorCreate(err)
		}

		m.timeEntries = te.Time_entries
		m.objectCount = len(m.timeEntries)

		m.cursor = 0
		m.crumbs = m.crumbs.addPage(TIME_ENTRIES)
	case tea.KeyRight: // go to next set of issues
		if m.issues.Offset+m.issues.Limit < m.issues.Total_count {
			var err error

			params := make(restapi.Params, 0)
			params["project_id"] = m.issues.Project_id
			params["limit"] = m.issues.Limit
			params["offset"] = m.issues.Offset + m.issues.Limit

			if m.filters.for_me {
				params["assigned_to_id"] = "me"
			}

			m.issues, err = m.redmineClient.GetIssues(params)
			if err != nil {
				return m.errorCreate(err)
			}

			m.objectCount = len(m.issues.Issues)
			m.cursor = 0
		}
	case tea.KeyLeft: // go to previos set of issues
		if m.issues.Offset-m.issues.Limit >= 0 {
			var err error

			params := make(restapi.Params, 0)
			params["project_id"] = m.issues.Project_id
			params["limit"] = m.issues.Limit
			params["offset"] = m.issues.Offset - m.issues.Limit

			if m.filters.for_me {
				params["assigned_to_id"] = "me"
			}

			m.issues, err = m.redmineClient.GetIssues(params)
			if err != nil {
				return m.errorCreate(err)
			}

			m.objectCount = len(m.issues.Issues)
			m.cursor = 0
		}
	case tea.KeyCtrlT:
		var err error
		m.filters.for_me = !m.filters.for_me

		params := make(restapi.Params, 0)
		params["project_id"] = m.issues.Project_id
		params["limit"] = m.issues.Limit

		if m.filters.for_me {
			params["assigned_to_id"] = "me"
		}

		m.issues, err = m.redmineClient.GetIssues(params)
		if err != nil {
			return m.errorCreate(err)
		}

		m.objectCount = len(m.issues.Issues)
		m.cursor = 0
	default:
		return m.navigation(msg)
	}

	return m, nil
}

// update logic if key tap on "time entries" page
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
		issue := m.issues.Issues[m.cursor]

		hours, err := strconv.ParseFloat(m.inputs[2].Value(), 32) // convert input hours string to float32
		if err != nil {
			return m.errorCreate(err)
		}

		status, err := m.redmineClient.CreateTimeEntry(
			issue.Id,
			m.inputs[1].Value(), // input date
			m.inputs[0].Value(), // input comment
			float32(hours),
		)
		if err != nil {
			return m.errorCreate(err)
		}
		m.status = status
	case tea.KeyCtrlQ: // go to previos page
		m.cursor = 0
		m.crumbs, _ = m.crumbs.popPage() // go to previos page
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

func (m model) updateTimeEntries(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	default:
		return m.navigation(msg)
	}
}

// update logic if tap key on "error" page
func (m model) updateError(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlQ:
			m.cursor = 0
			m.crumbs, _ = m.crumbs.popPage()
			return m, nil
		case tea.KeyEscape:
			return m, tea.Quit
		}
	}

	return m, nil
}

// creating error before view error page
func (m model) errorCreate(err error) (model, tea.Cmd) {
	m.err = err
	m.crumbs = m.crumbs.addPage(ERROR)
	return m, func() tea.Msg { return errMsg(err) }
}
