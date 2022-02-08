package cli

import (
	"strconv"
	"time"

	"github.com/alexey-sderzhikov/regent/restapi"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	projectsPage       = "projects"
	issuesPage         = "issues"
	inputTimeEntryPage = "input_time_entry"
	errPage            = "error"
	timeEntriesPage    = "time_entries"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch m.crumbs.getCurrentPage() {
		case projectsPage:
			return m.projectsHandler(msg)
		case issuesPage:
			return m.issuesHandler(msg)
		case inputTimeEntryPage:
			return m.inputTimeEntryHandler(msg)
		case timeEntriesPage:
			return m.timeEntriesHandler(msg)
		case errPage:
			return m.errorHandler(msg)
		}
	case errMsg:
		return m.errorHandler(msg)
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
	case tea.KeyCtrlH: // extend or reduce size of helper
		m.help.ShowAll = !m.help.ShowAll
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
func (m model) projectsHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter: // go to project issues
		var err error
		projectID := m.projects[m.cursor].ID

		params := make(restapi.Params, 0)
		params["project_id"] = projectID

		if m.filters.forMe {
			params["assigned_to_id"] = "me"
		}

		m.issues, err = m.redmineClient.GetIssues(params)
		if err != nil {
			return m.errorCreate(err)
		}

		m.objectCount = len(m.issues.Issues)

		m.cursor = 0
		m.crumbs = m.crumbs.addPage(issuesPage)
	default:
		return m.navigation(msg)
	}

	return m, nil
}

// TODO refactoring pagination switching
// update logic if key tap on "issues" page
func (m model) issuesHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter: // go to creation new time entry for issue
		m.inputs[1].SetValue(time.Now().Format("2006-01-02")) // set today date
		m.inputs[2].SetValue("8")                             // set 8 hour
		m.crumbs = m.crumbs.addPage(inputTimeEntryPage)
	case tea.KeyCtrlQ: // go to previos page
		m.cursor = 0
		m.crumbs, _ = m.crumbs.popPage()

		projects, err := m.redmineClient.GetProjects()
		if err != nil {
			return m.errorCreate(err)
		}

		m.projects = projects.Projects
		m.objectCount = len(m.projects)
	case tea.KeyCtrlA: // show my time entries
		var err error
		params := make(restapi.Params, 0)

		params["user_id"] = m.redmineClient.User.ID

		m.timeEntries, err = m.redmineClient.GetTimeEntryList(params)
		if err != nil {
			m.errorCreate(err)
		}

		m.objectCount = len(m.timeEntries.TimeEntries)

		m.cursor = 0
		m.crumbs = m.crumbs.addPage(timeEntriesPage)
	case tea.KeyCtrlT: // filter -show only my issues
		var err error
		m.filters.forMe = !m.filters.forMe

		params := make(restapi.Params, 0)
		params["project_id"] = m.issues.ProjectID
		params["limit"] = m.issues.Limit

		if m.filters.forMe {
			params["assigned_to_id"] = "me"
		}

		m.issues, err = m.redmineClient.GetIssues(params)
		if err != nil {
			return m.errorCreate(err)
		}

		m.objectCount = len(m.issues.Issues)
		m.cursor = 0
	case tea.KeyRight, tea.KeyLeft: // go to next or previous set of issues
		var err error
		params := make(restapi.Params, 0)

		// if key right or left and have opportunity for pagination,
		// then change 'offset' parameter
		// else do nothing
		if msg.Type == tea.KeyLeft && m.issues.Offset-m.issues.Limit >= 0 {
			params["offset"] = m.issues.Offset - m.issues.Limit
		} else if msg.Type == tea.KeyRight && m.issues.Offset+m.issues.Limit < m.issues.TotalCount {
			params["offset"] = m.issues.Offset + m.issues.Limit
		} else {
			return m, nil
		}

		params["project_id"] = m.issues.ProjectID
		params["limit"] = m.issues.Limit

		if m.filters.forMe {
			params["assigned_to_id"] = "me"
		}

		m.issues, err = m.redmineClient.GetIssues(params)
		if err != nil {
			return m.errorCreate(err)
		}

		m.objectCount = len(m.issues.Issues)
		m.cursor = 0

		return m, nil
	default:
		return m.navigation(msg)
	}

	return m, nil
}

// update logic if key tap on "time entries" page
func (m model) inputTimeEntryHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// clear status if it not empty while any key
	if m.status != "" {
		m.status = ""
	}

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
		date := m.inputs[1].Value()    // input date
		comment := m.inputs[0].Value() // input comment

		hours, err := strconv.ParseFloat(m.inputs[2].Value(), 32) // convert input hours string to float32
		if err != nil {
			return m.errorCreate(err)
		}

		status, err := m.redmineClient.CreateTimeEntry(
			issue.ID,
			date,
			comment,
			float32(hours),
		)

		if err != nil {
			return m.errorCreate(err)
		}

		m.status = status + "time entry at date " + date
		// TODO: this case duplicate code in navigation func
	case tea.KeyCtrlQ: // go to the previous page
		m.status = ""
		m.cursor = 0

		var err error
		m.crumbs, err = m.crumbs.popPage()
		if err != nil {
			return m.errorCreate(err)
		}
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

func (m model) timeEntriesHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRight, tea.KeyLeft:
		var err error
		params := make(restapi.Params, 0)

		// if key right or left and have opportunity for pagination,
		// then change 'offset' parameter
		// else do nothing
		if msg.Type == tea.KeyLeft && m.timeEntries.Offset-m.timeEntries.Limit >= 0 {
			params["offset"] = m.timeEntries.Offset - m.timeEntries.Limit
		} else if msg.Type == tea.KeyRight && m.timeEntries.Offset+m.timeEntries.Limit < m.timeEntries.TotalCount {
			params["offset"] = m.timeEntries.Offset + m.timeEntries.Limit
		} else {
			return m, nil
		}

		params["limit"] = m.timeEntries.Limit

		m.timeEntries, err = m.redmineClient.GetTimeEntryList(params)
		if err != nil {
			return m.errorCreate(err)
		}

		m.objectCount = len(m.timeEntries.TimeEntries)
		m.cursor = 0

		return m, nil
	default:
		return m.navigation(msg)
	}
}

// update logic if tap key on "error" page
func (m model) errorHandler(msg tea.Msg) (model, tea.Cmd) {
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
	m.crumbs = m.crumbs.addPage(errPage)
	return m, func() tea.Msg { return errMsg(err) }
}
