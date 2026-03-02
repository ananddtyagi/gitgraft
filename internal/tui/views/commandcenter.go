package views

import (
	"strings"

	"github.com/anandtyagi/gitgraft/internal/git"
	"github.com/anandtyagi/gitgraft/internal/tui/styles"
	"github.com/anandtyagi/gitgraft/internal/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a menu option
type MenuItem struct {
	Title       string
	Description string
	State       types.AppState
	Icon        string
}

// CommandCenterModel is the main menu view
type CommandCenterModel struct {
	width       int
	height      int
	items       []MenuItem
	filtered    []MenuItem
	cursor      int
	search      textinput.Model
	searching   bool
	message     string
	navigateTo  *types.AppState
	gitClient   *git.Client
}

// NewCommandCenter creates a new command center
func NewCommandCenter(client *git.Client) CommandCenterModel {
	search := textinput.New()
	search.Placeholder = "Search commands..."
	search.CharLimit = 50
	search.Width = 40

	items := []MenuItem{
		{
			Title:       "New Branch",
			Description: "Create a new branch from an existing one",
			State:       types.StateNewBranch,
			Icon:        "+",
		},
		{
			Title:       "Switch Branch",
			Description: "Switch to a different branch",
			State:       types.StateSwitchBranch,
			Icon:        "⇄",
		},
		{
			Title:       "Commit",
			Description: "Stage files and create a commit",
			State:       types.StateCommit,
			Icon:        "✓",
		},
	}

	return CommandCenterModel{
		items:     items,
		filtered:  items,
		search:    search,
		gitClient: client,
	}
}

// SetSize sets the view dimensions
func (m CommandCenterModel) SetSize(width, height int) CommandCenterModel {
	m.width = width
	m.height = height
	m.search.Width = min(40, width-10)
	return m
}

// SetMessage sets a status message to display
func (m CommandCenterModel) SetMessage(msg string) CommandCenterModel {
	m.message = msg
	return m
}

// NavigateTo returns the requested navigation state
func (m CommandCenterModel) NavigateTo() *types.AppState {
	nav := m.navigateTo
	m.navigateTo = nil
	return nav
}

// Update handles messages
func (m CommandCenterModel) Update(msg tea.Msg) (CommandCenterModel, tea.Cmd) {
	m.navigateTo = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "enter":
				m.searching = false
				m.search.Blur()
				if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
					state := m.filtered[m.cursor].State
					m.navigateTo = &state
				}
			case "esc":
				m.searching = false
				m.search.Blur()
				m.search.SetValue("")
				m.filtered = m.items
				m.cursor = 0
			case "up", "ctrl+p":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "ctrl+n":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
				}
			default:
				var cmd tea.Cmd
				m.search, cmd = m.search.Update(msg)
				m.filterItems()
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
					state := m.filtered[m.cursor].State
					m.navigateTo = &state
				}
			case "/":
				m.searching = true
				m.search.Focus()
				return m, textinput.Blink
			case "q":
				return m, tea.Quit
			}
		}
	}

	// Clear message after any keypress
	m.message = ""

	return m, nil
}

func (m *CommandCenterModel) filterItems() {
	query := strings.ToLower(m.search.Value())
	if query == "" {
		m.filtered = m.items
		m.cursor = 0
		return
	}

	m.filtered = []MenuItem{}
	for _, item := range m.items {
		if strings.Contains(strings.ToLower(item.Title), query) ||
			strings.Contains(strings.ToLower(item.Description), query) {
			m.filtered = append(m.filtered, item)
		}
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

// View renders the command center
func (m CommandCenterModel) View() string {
	// Header with logo and branch info
	header := m.renderHeader()

	// Search bar
	searchBar := m.renderSearchBar()

	// Menu items
	menu := m.renderMenu()

	// Status message
	var statusBar string
	if m.message != "" {
		statusBar = "\n" + styles.SuccessStyle.Render("✓ "+m.message)
	}

	// Help
	help := m.renderHelp()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		searchBar,
		"",
		menu,
		statusBar,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m CommandCenterModel) renderHeader() string {
	logo := styles.Logo()

	var branchInfo string
	if m.gitClient != nil {
		if branch, err := m.gitClient.GetCurrentBranch(); err == nil {
			branchInfo = styles.MutedText.Render("on ") +
				styles.CurrentBranchStyle.Render(branch)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		branchInfo,
	)
}

func (m CommandCenterModel) renderSearchBar() string {
	searchStyle := styles.InputStyle
	if m.searching {
		searchStyle = styles.FocusedInputStyle
	}

	return lipgloss.NewStyle().
		Width(m.search.Width + 4).
		Render(searchStyle.Render(m.search.View()))
}

func (m CommandCenterModel) renderMenu() string {
	var items []string

	for i, item := range m.filtered {
		var itemStyle lipgloss.Style
		if i == m.cursor {
			itemStyle = styles.SelectedItemStyle
		} else {
			itemStyle = styles.MenuItemStyle
		}

		icon := styles.MutedText.Render("[" + item.Icon + "]")
		title := item.Title
		desc := styles.SubtitleStyle.Render(item.Description)

		if i == m.cursor {
			icon = styles.PastelBlue.Render("[" + item.Icon + "]")
			title = styles.PastelBlue.Bold(true).Render(item.Title)
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			icon,
			" ",
			lipgloss.JoinVertical(
				lipgloss.Left,
				itemStyle.Render(title),
				"    "+desc,
			),
		)

		items = append(items, line)
	}

	if len(items) == 0 {
		return styles.MutedText.Render("No matching commands")
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func (m CommandCenterModel) renderHelp() string {
	if m.searching {
		return styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel")
	}
	return styles.HelpStyle.Render("↑/↓ navigate • Enter select • / search • q quit")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
