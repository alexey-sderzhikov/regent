package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var border = lipgloss.Border{
	Top:         "-",
	Bottom:      "-",
	Left:        "|",
	Right:       "|",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}

var (
	titleStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	cursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	currentLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true).MarginLeft(2)
	crumbsStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	statusStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	filterStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#00a86b"))
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	textStyle        = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
)

func (m model) View() string {
	var header, body, tail string

	header = crumbsStyle.Render(m.crumbs.printStack())

	switch m.crumbs.getCurrentPage() {
	case PROJECTS:
		body = m.viewProjects()
	case ISSUES:
		body = m.viewIssues()
	case INPUT_TIME_ENTRY:
		body = m.viewInputTimeEntry()
	case TIME_ENTRIES:
		body = m.viewTimeEntries()
	case ERROR:
		body = m.viewError()
	}

	tail = m.help.View(m.key)

	if body == "" {
		return "Cannot detect current page :("
	}

	return strings.Join([]string{header, body, tail}, "\n")
}

func (m model) viewProjects() string {
	s := titleStyle.Render("Bergen Projects")

	s += "\n"

	var mainText string
	for ind, p := range m.projects {
		cursor := " "
		name := p.Name
		if m.cursor == ind {
			cursor = cursorStyle.Render(">")
			name = currentLineStyle.Render(name)
		}

		mainText += fmt.Sprintf("%s %s\n", cursor, name)
	}

	s += textStyle.Render(mainText)

	return s
}

func (m model) viewIssues() string {
	s := titleStyle.Render(fmt.Sprintf(
		"Issues (%v) project's #%v", m.issues.Total_count, m.issues.Project_id),
	)

	s += "\n"

	var mainText string
	mainText += fmt.Sprintf(
		"Show from %v to %v issues. Total - %v\n",
		m.issues.Offset+1,
		m.issues.Offset+m.issues.Limit,
		m.issues.Total_count,
	)
	mainText += filterStyle.Render(fmt.Sprintf("Issues for me: %v", m.filters.for_me))
	mainText += "\n\n"
	if len(m.issues.Issues) == 0 {
		mainText += "None suitable issues\n"
	} else {
		for ind, i := range m.issues.Issues {
			cursor := " "
			subject := i.Subject
			if m.cursor == ind {
				cursor = cursorStyle.Render(">")
				subject = currentLineStyle.Render(subject)
			}

			mainText += fmt.Sprintf("%s %s\n", cursor, subject)
		}
	}

	s += textStyle.Render(mainText)

	return s
}
func (m model) viewTimeEntries() string {
	s := titleStyle.Render(fmt.Sprintf("%s Time Entries", m.redmineClient.User.Lastname))

	s += "\n"

	var mainText string
	mainText += fmt.Sprintf(
		"Show from %v to %v issues. Total - %v\n",
		m.timeEntries.Offset+1,
		m.timeEntries.Offset+m.issues.Limit,
		m.timeEntries.Total_count,
	)
	for ind, te := range m.timeEntries.Time_entries {
		cursor := " "
		spent_on := te.Spent_on
		comment := te.Comments
		hours := fmt.Sprintf("%v", te.Hours)
		issueId := fmt.Sprintf("%v", te.Issue.Id)
		if m.cursor == ind {
			cursor = cursorStyle.Render(">")
			spent_on = currentLineStyle.Render(spent_on)
			comment = currentLineStyle.Render(comment)
			hours = currentLineStyle.Render(hours)
			issueId = currentLineStyle.Render(issueId)
		}

		mainText += fmt.Sprintf("%s %s %s %s %s\n", cursor, spent_on, issueId, hours, comment)
	}

	s += textStyle.Render(mainText)

	return s
}

func (m model) viewInputTimeEntry() string {
	var mainText, s string
	mainText += fmt.Sprint(
		"Text comment to time entry:\n",
		m.inputs[0].View(),
		"\nText date:\n",
		m.inputs[1].View(),
		"\nText work hours:\n",
		m.inputs[2].View(),
		"\n",
	)

	if m.status != "" {
		s += statusStyle.Render(m.status)
		s += "\n"
	}

	s += textStyle.Render(mainText)

	return s
}

func (m model) viewError() string {

	s := "\n\n"
	s += errorStyle.Render("Error! - " + fmt.Sprint(m.err))

	return s
}
