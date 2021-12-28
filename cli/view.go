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
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
var textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))

func (m model) View() string {
	switch m.page {
	case PROJECTS:
		return m.viewProjects()
	case ISSUES:
		return m.viewIssues()
	case INPUT_TIME_ENTRY:
		return m.viewInputTimeEntry()
	case ERROR:
		return m.viewError()
	}

	return "Cannot detect current page :("
}

func (m model) viewProjects() string {

	s := titleStyle.Render("Bergen Projects")
	s += "\n\n"

	for ind, p := range m.projects {
		cursor := " "
		name := p.Name
		if m.cursor == ind {
			cursor = cursorStyle.Render(">")
			name = textStyle.Render(name)
		}

		s += fmt.Sprintf("%s %s\n", cursor, name)
	}

	s += "\nPress q to quit.\n"

	return s
}

func (m model) viewIssues() string {
	s := titleStyle.Render(fmt.Sprintf(
		"\nIssues project's - %q\n", m.issues[0].Project.Name),
	)
	s += "\n\n"

	for ind, i := range m.issues {
		cursor := " "
		subject := i.Subject
		if m.cursor == ind {
			cursor = cursorStyle.Render(">")
			subject = textStyle.Render(subject)
		}

		s += fmt.Sprintf("%s %s\n", cursor, subject)
	}

	s += "\nPress q to quit.\n"

	return s
}

func (m model) viewInputTimeEntry() string {

	return fmt.Sprint(
		"\nText comment to time entry:\n",
		m.inputs[0].View(),
		"\nText date:\n",
		m.inputs[1].View(),
		"\nText work hours:\n",
		m.inputs[2].View(),
		"\n",
	)
}

func (m model) viewError() string {
	return fmt.Sprint(m.err)
}
