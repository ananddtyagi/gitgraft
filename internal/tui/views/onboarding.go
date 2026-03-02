package views

import (
	"github.com/anandtyagi/gitgraft/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OnboardingStep represents the current step in onboarding
type OnboardingStep int

const (
	StepWelcome OnboardingStep = iota
	StepAlias
	StepCustomAlias
	StepComplete
)

// OnboardingModel handles the first-run experience
type OnboardingModel struct {
	step       OnboardingStep
	width      int
	height     int
	cursor     int
	aliasInput textinput.Model
	alias      string
	done       bool
}

// NewOnboarding creates a new onboarding model
func NewOnboarding() OnboardingModel {
	ti := textinput.New()
	ti.Placeholder = "alias"
	ti.CharLimit = 10
	ti.Width = 20

	return OnboardingModel{
		step:       StepWelcome,
		cursor:     0,
		aliasInput: ti,
	}
}

// SetSize sets the view dimensions
func (m OnboardingModel) SetSize(width, height int) OnboardingModel {
	m.width = width
	m.height = height
	return m
}

// Done returns true when onboarding is complete
func (m OnboardingModel) Done() bool {
	return m.done
}

// Alias returns the configured alias
func (m OnboardingModel) Alias() string {
	return m.alias
}

// Update handles messages
func (m OnboardingModel) Update(msg tea.Msg) (OnboardingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case StepWelcome:
			switch msg.String() {
			case "enter", " ":
				m.step = StepAlias
			case "esc", "q":
				m.done = true
			}

		case StepAlias:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < 2 {
					m.cursor++
				}
			case "enter":
				switch m.cursor {
				case 0: // Yes (gg)
					m.alias = "gg"
					m.step = StepComplete
				case 1: // Custom
					m.step = StepCustomAlias
					m.aliasInput.Focus()
					return m, textinput.Blink
				case 2: // No
					m.alias = ""
					m.step = StepComplete
				}
			case "esc":
				m.step = StepWelcome
			}

		case StepCustomAlias:
			switch msg.String() {
			case "enter":
				if m.aliasInput.Value() != "" {
					m.alias = m.aliasInput.Value()
					m.step = StepComplete
				}
			case "esc":
				m.step = StepAlias
				m.aliasInput.Blur()
			default:
				var cmd tea.Cmd
				m.aliasInput, cmd = m.aliasInput.Update(msg)
				return m, cmd
			}

		case StepComplete:
			switch msg.String() {
			case "enter", " ":
				m.done = true
			}
		}
	}

	return m, nil
}

// View renders the onboarding screen
func (m OnboardingModel) View() string {
	var content string

	switch m.step {
	case StepWelcome:
		content = m.renderWelcome()
	case StepAlias:
		content = m.renderAliasChoice()
	case StepCustomAlias:
		content = m.renderCustomAlias()
	case StepComplete:
		content = m.renderComplete()
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m OnboardingModel) renderWelcome() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Logo(),
		"",
		"",
		styles.TitleStyle.Render("Welcome to Git-Graft!"),
		"",
		styles.MutedText.Render("A chill TUI for your git workflow"),
		"",
		"",
		styles.HelpStyle.Render("Press Enter to continue"),
	)
}

func (m OnboardingModel) renderAliasChoice() string {
	title := styles.TitleStyle.Render("Set up a shell alias?")
	subtitle := styles.MutedText.Render("Quick access from your terminal")

	options := []string{
		"Yes, use 'gg' (recommended)",
		"Custom alias",
		"No thanks",
	}

	var optionLines string
	for i, opt := range options {
		if i == m.cursor {
			optionLines += styles.SelectedItemStyle.Render("→ " + opt) + "\n"
		} else {
			optionLines += styles.MenuItemStyle.Render("  " + opt) + "\n"
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Logo(),
		"",
		title,
		subtitle,
		"",
		optionLines,
		"",
		styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc back"),
	)

	return styles.PanelStyle.
		Width(50).
		Render(content)
}

func (m OnboardingModel) renderCustomAlias() string {
	title := styles.TitleStyle.Render("Enter custom alias")

	inputStyle := styles.FocusedInputStyle.Width(20)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Logo(),
		"",
		title,
		"",
		inputStyle.Render(m.aliasInput.View()),
		"",
		styles.HelpStyle.Render("Enter confirm • Esc back"),
	)

	return styles.PanelStyle.
		Width(50).
		Render(content)
}

func (m OnboardingModel) renderComplete() string {
	var aliasMsg string
	if m.alias != "" {
		aliasMsg = styles.SuccessStyle.Render("Alias '" + m.alias + "' will be added to your shell config")
	} else {
		aliasMsg = styles.MutedText.Render("No alias configured")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Logo(),
		"",
		styles.TitleStyle.Render("You're all set!"),
		"",
		aliasMsg,
		"",
		styles.HelpStyle.Render("Press Enter to start using Git-Graft"),
	)

	return styles.PanelStyle.
		Width(50).
		Render(content)
}
