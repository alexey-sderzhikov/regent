package cli

import (
	"fmt"
	"os"

	"github.com/alexey-sderzhikov/regent/restapi"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	help          help.Model
	key           keyMap
	status        string
	err           error
}

type filterStruct struct {
	for_me bool
}

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Help       key.Binding
	Quit       key.Binding
	Back       key.Binding
	Select     key.Binding
	MyIssues   key.Binding
	AllEntries key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Select}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},                  // first column
		{k.Quit, k.Select},                               // second column
		{k.MyIssues, k.AllEntries, k.Help, k.AllEntries}, // third column
	}
}

var keys = keyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "previous elements"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→", "next elements"),
	),
	Help: key.NewBinding(
		key.WithKeys("CtrlH"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Back: key.NewBinding(
		key.WithKeys("CtrlQ"),
		key.WithHelp("ctrl+q", "go back"),
	),
	MyIssues: key.NewBinding(
		key.WithKeys("CtrlT"),
		key.WithHelp("ctrl+t", "show only my issues"),
	),
	AllEntries: key.NewBinding(
		key.WithKeys("CtrlA"),
		key.WithHelp("ctrl+a", "go to time entries"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	),
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

	// open config file .env and take redmine USER API KEY and URL
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

	// create redmine client he do all request to redmine server
	rc, err := restapi.NewRm(source, apiKey)
	if err != nil {
		return model{}, fmt.Errorf("error occure during creating redmine client object\n%q", err)
	}
	m.redmineClient = rc

	m.help = help.New()
	m.key = keys

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
