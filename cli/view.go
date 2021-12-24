package cli

import "fmt"

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
	s := "\nBergen Projects\n"

	for ind, p := range m.projects {
		cursor := " "
		if m.cursor == ind {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, p.Name)
	}

	s += "\nPress q to quit.\n"

	return s
}

func (m model) viewIssues() string {
	s := fmt.Sprintf("\nIssues project's - %q\n", m.issues[0].Project.Name)

	for ind, i := range m.issues {
		cursor := " "
		if m.cursor == ind {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, i.Subject)
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
