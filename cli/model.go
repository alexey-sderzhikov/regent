package cli

import (
	"github.com/alexey-sderzhikov/regent/restapi"
	"github.com/charmbracelet/bubbles/textinput"
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
	err           error
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
