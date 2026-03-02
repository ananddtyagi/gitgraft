package views

import (
	"strings"

	"github.com/anandtyagi/gitgraft/internal/git"
	"github.com/anandtyagi/gitgraft/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NewBranchModel handles branch creation
type NewBranchModel struct {
	width        int
	height       int
	fromInput    textinput.Model
	nameInput    textinput.Model
	focusIndex   int
	gitClient    *git.Client
	branches     []string
	shouldCreate bool
	shouldGoBack bool
}

// NewNewBranch creates a new branch creation view
func NewNewBranch(client *git.Client) NewBranchModel {
	fromInput := textinput.New()
	fromInput.Placeholder = "base branch"
	fromInput.CharLimit = 100
	fromInput.Width = 40

	nameInput := textinput.New()
	nameInput.Placeholder = "new-branch-name"
	nameInput.CharLimit = 100
	nameInput.Width = 40
	nameInput.Focus()

	// Set default base branch
	if client != nil {
		fromInput.SetValue(client.GetDefaultBranch())
	}

	return NewBranchModel{
		fromInput:  fromInput,
		nameInput:  nameInput,
		focusIndex: 0, // Start on name input
		gitClient:  client,
	}
}

// SetSize sets the view dimensions
func (m NewBranchModel) SetSize(width, height int) NewBranchModel {
	m.width = width
	m.height = height
	inputWidth := min(40, width-20)
	m.fromInput.Width = inputWidth
	m.nameInput.Width = inputWidth
	return m
}

// SetBranches sets available branches for autocomplete
func (m NewBranchModel) SetBranches(branches []git.Branch) NewBranchModel {
	m.branches = make([]string, len(branches))
	for i, b := range branches {
		m.branches[i] = b.Name
	}
	return m
}

// ShouldCreate returns true if branch should be created
func (m NewBranchModel) ShouldCreate() bool {
	return m.shouldCreate
}

// ShouldGoBack returns true if should navigate back
func (m NewBranchModel) ShouldGoBack() bool {
	return m.shouldGoBack
}

// GetBranchInfo returns the branch name and base
func (m NewBranchModel) GetBranchInfo() (name, base string) {
	return m.nameInput.Value(), m.fromInput.Value()
}

// Reset resets the view state
func (m NewBranchModel) Reset() NewBranchModel {
	m.shouldCreate = false
	m.shouldGoBack = false
	m.nameInput.SetValue("")
	m.focusIndex = 0
	m.nameInput.Focus()
	m.fromInput.Blur()
	return m
}

// Update handles messages
func (m NewBranchModel) Update(msg tea.Msg) (NewBranchModel, tea.Cmd) {
	m.shouldCreate = false
	m.shouldGoBack = false

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			// Toggle focus
			m.focusIndex = (m.focusIndex + 1) % 2
			if m.focusIndex == 0 {
				m.nameInput.Focus()
				m.fromInput.Blur()
			} else {
				m.fromInput.Focus()
				m.nameInput.Blur()
			}
			return m, textinput.Blink

		case "enter":
			name := strings.TrimSpace(m.nameInput.Value())
			if name != "" {
				m.shouldCreate = true
			}

		case "esc":
			m.shouldGoBack = true

		default:
			var cmd tea.Cmd
			if m.focusIndex == 0 {
				m.nameInput, cmd = m.nameInput.Update(msg)
			} else {
				m.fromInput, cmd = m.fromInput.Update(msg)
			}
			return m, cmd
		}
	}

	return m, nil
}

// View renders the branch creation view
func (m NewBranchModel) View() string {
	title := styles.TitleStyle.Render("Create New Branch")

	// Branch name input
	nameLabel := styles.LabelStyle.Render("Branch name")
	nameStyle := styles.InputStyle
	if m.focusIndex == 0 {
		nameLabel = styles.FocusedLabelStyle.Render("Branch name")
		nameStyle = styles.FocusedInputStyle
	}
	nameInput := nameStyle.Width(m.nameInput.Width + 2).Render(m.nameInput.View())

	// From branch input
	fromLabel := styles.LabelStyle.Render("From branch")
	fromStyle := styles.InputStyle
	if m.focusIndex == 1 {
		fromLabel = styles.FocusedLabelStyle.Render("From branch")
		fromStyle = styles.FocusedInputStyle
	}
	fromInput := fromStyle.Width(m.fromInput.Width + 2).Render(m.fromInput.View())

	// Current branch hint
	var currentHint string
	if m.gitClient != nil {
		if current, err := m.gitClient.GetCurrentBranch(); err == nil {
			currentHint = styles.MutedText.Render("Currently on: ") +
				styles.CurrentBranchStyle.Render(current)
		}
	}

	// Validation hint
	var hint string
	name := strings.TrimSpace(m.nameInput.Value())
	if name == "" {
		hint = styles.MutedText.Render("Enter a branch name")
	} else if strings.Contains(name, " ") {
		hint = styles.WarningStyle.Render("Branch names cannot contain spaces")
	} else {
		hint = styles.SuccessStyle.Render("Ready to create branch")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		currentHint,
		"",
		nameLabel,
		nameInput,
		"",
		fromLabel,
		fromInput,
		"",
		hint,
	)

	help := styles.HelpStyle.Render("Tab switch fields • Enter create • Esc back")

	panel := styles.PanelStyle.
		Width(min(60, m.width-4)).
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
