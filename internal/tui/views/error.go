package views

import (
	"strings"

	"github.com/anandtyagi/gitgraft/internal/git"
	"github.com/anandtyagi/gitgraft/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrorModel handles error display with recovery actions
type ErrorModel struct {
	width          int
	height         int
	err            *git.GitError
	cursor         int
	selectedAction *git.Action
	showOutput     bool
}

// NewError creates a new error view
func NewError(err *git.GitError) ErrorModel {
	return ErrorModel{
		err:    err,
		cursor: 0,
	}
}

// SetSize sets the view dimensions
func (m ErrorModel) SetSize(width, height int) ErrorModel {
	m.width = width
	m.height = height
	return m
}

// SelectedAction returns the selected action
func (m ErrorModel) SelectedAction() *git.Action {
	return m.selectedAction
}

// Update handles messages
func (m ErrorModel) Update(msg tea.Msg) (ErrorModel, tea.Cmd) {
	m.selectedAction = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showOutput {
			// Output view mode
			switch msg.String() {
			case "esc", "q", "enter":
				m.showOutput = false
			}
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.err != nil && m.cursor < len(m.err.Actions)-1 {
				m.cursor++
			}
		case "enter":
			if m.err != nil && len(m.err.Actions) > 0 && m.cursor < len(m.err.Actions) {
				action := m.err.Actions[m.cursor]

				// Handle "View output" specially
				if action.ActionType == git.ActionRetry && strings.Contains(action.Label, "View") {
					m.showOutput = true
					return m, nil
				}

				m.selectedAction = &action
			}
		case "o":
			// Quick toggle output view
			if m.err != nil && m.err.Output != "" {
				m.showOutput = true
			}
		case "esc":
			// Select cancel action if available
			if m.err != nil {
				for _, action := range m.err.Actions {
					if action.ActionType == git.ActionCancel {
						m.selectedAction = &action
						break
					}
				}
			}
		}
	}

	return m, nil
}

// View renders the error view
func (m ErrorModel) View() string {
	if m.err == nil {
		return ""
	}

	if m.showOutput {
		return m.renderOutput()
	}

	return m.renderError()
}

func (m ErrorModel) renderError() string {
	// Error icon and title
	icon := styles.ErrorStyle.Render("✗")
	title := styles.ErrorStyle.Bold(true).Render("Error: " + m.err.Op)

	// Error message
	message := styles.PrimaryText.Render(m.err.Message)

	// Actions
	var actionLines []string
	for i, action := range m.err.Actions {
		var line string
		if i == m.cursor {
			line = styles.SelectedItemStyle.Render("→ " + action.Label)
		} else {
			line = styles.MenuItemStyle.Render("  " + action.Label)
		}

		// Add description on next line
		desc := styles.SubtleText.Render("    " + action.Description)
		actionLines = append(actionLines, line, desc)
	}
	actions := strings.Join(actionLines, "\n")

	// Output hint
	var outputHint string
	if m.err.Output != "" {
		outputHint = styles.SubtleText.Render("\nPress 'o' to view full output")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, icon, " ", title),
		"",
		message,
		outputHint,
		"",
		styles.LabelStyle.Render("What would you like to do?"),
		"",
		actions,
	)

	help := styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel")

	panel := styles.PanelStyle.
		Width(min(70, m.width-4)).
		BorderForeground(styles.ErrorText).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			panel,
			"",
			help,
		),
	)
}

func (m ErrorModel) renderOutput() string {
	title := styles.TitleStyle.Render("Command Output")

	// Format output with some wrapping
	output := m.err.Output
	if len(output) > 2000 {
		output = output[:2000] + "\n... (truncated)"
	}

	// Style the output
	outputStyled := styles.MutedText.Render(output)

	maxHeight := m.height - 10
	lines := strings.Split(outputStyled, "\n")
	if len(lines) > maxHeight {
		lines = lines[:maxHeight]
		lines = append(lines, styles.SubtleText.Render("... (scroll to see more)"))
	}
	outputStyled = strings.Join(lines, "\n")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		outputStyled,
	)

	help := styles.HelpStyle.Render("Press Esc, Enter, or 'q' to go back")

	panel := styles.PanelStyle.
		Width(min(80, m.width-4)).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			panel,
			"",
			help,
		),
	)
}
