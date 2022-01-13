package cli

import (
	"fmt"

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

var titleStyle = lipgloss.NewStyle().
	BorderStyle(border).
	Bold(true).
	PaddingLeft(3).
	PaddingRight(3).
	Foreground(lipgloss.Color("#7D56F4"))

var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
var currentLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true).MarginLeft(2)
var crumbsStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
var statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
var textStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
var helperText = "\n'Esc' - close; 'Ctrl+q' - go to previos page\n"

func (m model) View() string {
	switch m.crumbs.getCurrentPage() {
	case PROJECTS:
		return m.viewProjects()
	case ISSUES:
		return m.viewIssues()
	case INPUT_TIME_ENTRY:
		return m.viewInputTimeEntry()
	case TIME_ENTRIES:
		return m.viewTimeEntries()
	case ERROR:
		return m.viewError()
	}

	return "Cannot detect current page :("
}

func (m model) viewProjects() string {
	s := titleStyle.Render("Bergen Projects")
	s += "\n"
	s += crumbsStyle.Render(
		m.crumbs.printStack(),
	)

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
	s += helperText

	return s
}

func (m model) viewIssues() string {
	s := titleStyle.Render(fmt.Sprintf(
		"Issues (%v) project's - %q", m.issues.Total_count, m.issues.Issues[0].Project.Name),
	)
	s += "\n"
	s += crumbsStyle.Render(
		m.crumbs.printStack(),
	)
	s += "\n"

	var mainText string
	mainText += fmt.Sprintf(
		"Issues from %v to %v. Total issues - %v\n\n",
		m.issues.Offset,
		m.issues.Offset+m.issues.Limit,
		m.issues.Total_count,
	)
	for ind, i := range m.issues.Issues {
		cursor := " "
		subject := i.Subject
		if m.cursor == ind {
			cursor = cursorStyle.Render(">")
			subject = currentLineStyle.Render(subject)
		}

		mainText += fmt.Sprintf("%s %s\n", cursor, subject)
	}

	s += textStyle.Render(mainText)
	s += helperText

	return s
}
func (m model) viewTimeEntries() string {
	s := titleStyle.Render(fmt.Sprintf("%s Time Entries", m.redmineClient.User.Lastname))
	s += "\n"
	s += crumbsStyle.Render(
		m.crumbs.printStack(),
	)

	s += "\n"

	for ind, te := range m.timeEntries {
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

		s += fmt.Sprintf("%s %s %s %s %s\n", cursor, spent_on, issueId, hours, comment)
	}

	s += helperText

	return s
}

func (m model) viewInputTimeEntry() string {
	s := crumbsStyle.Render(
		m.crumbs.printStack(),
	)

	s += "\n"

	var mainText string
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
	s += helperText

	return s
}

func (m model) viewError() string {
	s := crumbsStyle.Render(
		m.crumbs.printStack(),
	)

	s += "\n\n"
	s += errorStyle.Render("Error! - " + fmt.Sprint(m.err))
	s += "\n"

	s += helperText
	return s
}
