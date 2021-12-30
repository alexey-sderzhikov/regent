package cli

import (
	"fmt"

	"github.com/alexey-sderzhikov/regent/restapi"
	"github.com/charmbracelet/bubbles/textinput"
)

const (
	COMMENT_FIELD = "comment"
	DATE_FIELD    = "date"
	HOURS_FIELD   = "hours"
)

type errMsg error

type model struct {
	redmineClient *restapi.RmClient
	projects      []restapi.Project
	issues        []restapi.Issue
	timeEntries   []restapi.TimeEntryResponse
	inputs        []textinput.Model
	focusIndex    int
	objectCount   int
	cursor        int
	crumbs        pagesStack
	status        string
	err           error
}

type pagesStack []string

func (p pagesStack) addPage(page string) pagesStack {
	return append(p, page)
}

func (p pagesStack) popPage() (pagesStack, error) {
	if len(p) <= 1 {
		return p, fmt.Errorf("")
	}
	return p[:len(p)-1]
}

func (p pagesStack) printStack() string {
	var res string = "/"
	for _, page := range p {
		res += page + "/"
	}
	return res
}

func (p pagesStack) getCurrentPage() string {
	return p[len(p)-1]
}

func initialModel() (model, error) {
	m := model{}

	rc, err := restapi.NewRm("", "")
	if err != nil {
		return model{}, err
	}
	m.redmineClient = rc

	projects, err := rc.GetProjects()
	if err != nil {
		return model{}, err
	}
	m.projects = projects.Projects
	m.objectCount = len(m.projects)

	m.crumbs = pagesStack{PROJECTS}

	m.inputs = make([]textinput.Model, 3)
	m.inputs[0] = initialCommentInput()
	m.inputs[1] = initialDateInput()
	m.inputs[2] = initialHoursInput()

	m.focusIndex = 0

	return m, nil
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
