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
	case projectsPage:
		body = m.viewProjects()
	case issuesPage:
		body = m.viewIssues()
	case inputTimeEntryPage:
		body = m.viewInputTimeEntry()
	case timeEntriesPage:
		body = m.viewTimeEntries()
	case errPage:
		body = m.viewError()
	}

	tail = m.help.View(m.key)

	if body == "" {
		return "Cannot detect current page :("
	}

	return strings.Join([]string{header, body, tail}, "\n")
}

func (m model) viewProjects() string {
	var view strings.Builder

	view.WriteString(titleStyle.Render("Bergen Projects") + "\n")

	for ind, p := range m.projects {
		cursor := " "
		name := p.Name
		if m.cursor == ind {
			cursor = cursorStyle.Render(">")
			name = currentLineStyle.Render(name)
		}

		view.WriteString(fmt.Sprintf("%s %s\n", cursor, name))
	}

	return textStyle.Render(view.String())
}

func (m model) viewIssues() string {
	var view strings.Builder

	view.WriteString(
		titleStyle.Render(fmt.Sprintf(
			"Issues (%v) project's #%v", m.issues.TotalCount, m.issues.ProjectID),
		) + "\n",
	)

	view.WriteString(fmt.Sprintf(
		"Show from %v to %v issues. Total - %v\n",
		m.issues.Offset+1,
		m.issues.Offset+m.issues.Limit,
		m.issues.TotalCount,
	))

	view.WriteString(
		filterStyle.Render(
			fmt.Sprintf("Issues for me: %v", m.filters.forMe),
		) + "\n\n",
	)

	if len(m.issues.Issues) == 0 {
		view.WriteString("None suitable issues\n")
	} else {
		for ind, i := range m.issues.Issues {
			cursor := " "
			subject := i.Subject
			if m.cursor == ind {
				cursor = cursorStyle.Render(">")
				subject = currentLineStyle.Render(subject)
			}

			view.WriteString(fmt.Sprintf("%s %s\n", cursor, subject))
		}
	}

	return textStyle.Render(view.String())
}

func (m model) viewTimeEntries() string {
	var view strings.Builder

	view.WriteString(
		titleStyle.Render(
			fmt.Sprintf("%s Time Entries", m.redmineClient.User.Lastname),
		) + "\n",
	)

	view.WriteString(
		fmt.Sprintf(
			"Show from %v to %v issues. Total - %v\n",
			m.timeEntries.Offset+1,
			m.timeEntries.Offset+m.issues.Limit,
			m.timeEntries.TotalCount,
		),
	)

	for ind, te := range m.timeEntries.TimeEntries {
		cursor := " "
		spentOn := te.SpentOn
		comment := te.Comments
		hours := fmt.Sprintf("%v", te.Hours)
		issueID := fmt.Sprintf("%v", te.Issue)
		if m.cursor == ind {
			cursor = cursorStyle.Render(">")
			spentOn = currentLineStyle.Render(spentOn)
			comment = currentLineStyle.Render(comment)
			hours = currentLineStyle.Render(hours)
			issueID = currentLineStyle.Render(issueID)
		}

		view.WriteString(fmt.Sprintf("%s %s %s %s %s\n", cursor, spentOn, issueID, hours, comment))
	}

	return textStyle.Render(view.String())
}

func (m model) viewInputTimeEntry() string {
	var view strings.Builder

	if m.status != "" {
		view.WriteString(
			statusStyle.Render(m.status) + "\n",
		)
	}

	view.WriteString(
		fmt.Sprint(
			"Text comment to time entry:\n",
			m.inputs[0].View(),
			"\nText date:\n",
			m.inputs[1].View(),
			"\nText work hours:\n",
			m.inputs[2].View(),
			"\n",
		),
	)

	return textStyle.Render(view.String())
}

func (m model) viewError() string {
	return errorStyle.Render("\n\nError! - " + fmt.Sprint(m.err))
}
