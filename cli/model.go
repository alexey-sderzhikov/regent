package cli

import (
	"fmt"
	"os"

	"github.com/alexey-sderzhikov/regent/restapi"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/joho/godotenv"
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
	issues        restapi.IssueList
	timeEntries   restapi.TimeEntryListResponse
	inputs        []textinput.Model
	focusIndex    int        // need for switch between input fields
	objectCount   int        // need View for correct switch between elements
	cursor        int        // current select line
	crumbs        pagesStack // bread crumbs
	filters       filterStruct
	status        string
	err           error
}

type filterStruct struct {
	for_me bool
}

type pagesStack []string

func (p pagesStack) addPage(page string) pagesStack {
	return append(p, page)
}

func (p pagesStack) popPage() (pagesStack, error) {
	if len(p) <= 1 {
		return p, fmt.Errorf("its last page in stack")
	}
	return p[:len(p)-1], nil
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

	err := godotenv.Load(".env")
	if err != nil {
		return model{}, fmt.Errorf("error occure during reading .env file\n%q", err)
	}

	apiKey := os.Getenv("USER_API_KEY")
	if apiKey == "" {
		return model{}, fmt.Errorf("api key in .env file is nil")
	}

	source := os.Getenv("SOURCE")
	if source == "" {
		return model{}, fmt.Errorf("source in .env file is nil")
	}

	rc, err := restapi.NewRm(source, apiKey)
	if err != nil {
		return model{}, fmt.Errorf("error occure during creating redmine client object\n%q", err)
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

	m.filters.for_me = false

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
	ti.CharLimit = 254
	ti.Width = 30

	return ti
}
