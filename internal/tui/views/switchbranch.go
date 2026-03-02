package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/anandtyagi/gitgraft/internal/git"
	"github.com/anandtyagi/gitgraft/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SwitchBranchModel handles branch switching with preview
type SwitchBranchModel struct {
	width            int
	height           int
	branches         []git.Branch
	filtered         []git.Branch
	commits          map[string][]git.CommitInfo
	cursor           int
	filter           textinput.Model
	filtering        bool
	currentBranch    string
	selectedBranch   string
	branchForCommits string
	shouldGoBack     bool
	gitClient        *git.Client
}

// NewSwitchBranch creates a new branch switch view
func NewSwitchBranch(client *git.Client) SwitchBranchModel {
	filter := textinput.New()
	filter.Placeholder = "Filter branches..."
	filter.CharLimit = 50
	filter.Width = 30

	return SwitchBranchModel{
		filter:    filter,
		commits:   make(map[string][]git.CommitInfo),
		gitClient: client,
	}
}

// SetSize sets the view dimensions
func (m SwitchBranchModel) SetSize(width, height int) SwitchBranchModel {
	m.width = width
	m.height = height
	m.filter.Width = min(30, width/3)
	return m
}

// SetBranches sets the available branches
func (m SwitchBranchModel) SetBranches(branches []git.Branch, current string) SwitchBranchModel {
	m.branches = branches
	m.filtered = branches
	m.currentBranch = current
	m.cursor = 0
	return m
}

// SetCommits sets commits for a branch
func (m SwitchBranchModel) SetCommits(branch string, commits []git.CommitInfo) SwitchBranchModel {
	m.commits[branch] = commits
	m.branchForCommits = ""
	return m
}

// SelectedBranch returns the branch to switch to
func (m SwitchBranchModel) SelectedBranch() string {
	result := m.selectedBranch
	m.selectedBranch = ""
	return result
}

// BranchForCommits returns the branch that needs commits loaded
func (m SwitchBranchModel) BranchForCommits() string {
	return m.branchForCommits
}

// ShouldGoBack returns true if should navigate back
func (m SwitchBranchModel) ShouldGoBack() bool {
	return m.shouldGoBack
}

// Reset resets the view state
func (m SwitchBranchModel) Reset() SwitchBranchModel {
	m.selectedBranch = ""
	m.shouldGoBack = false
	m.filter.SetValue("")
	m.filtered = m.branches
	m.cursor = 0
	return m
}

// Update handles messages
func (m SwitchBranchModel) Update(msg tea.Msg) (SwitchBranchModel, tea.Cmd) {
	m.selectedBranch = ""
	m.shouldGoBack = false
	m.branchForCommits = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filtering {
			switch msg.String() {
			case "enter":
				m.filtering = false
				m.filter.Blur()
				if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
					branch := m.filtered[m.cursor]
					if branch.Name != m.currentBranch {
						m.selectedBranch = branch.Name
					}
				}
			case "esc":
				m.filtering = false
				m.filter.Blur()
				m.filter.SetValue("")
				m.filtered = m.branches
				m.cursor = 0
			case "up", "ctrl+p":
				if m.cursor > 0 {
					m.cursor--
					m.loadCommitsForCursor()
				}
			case "down", "ctrl+n":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
					m.loadCommitsForCursor()
				}
			default:
				var cmd tea.Cmd
				m.filter, cmd = m.filter.Update(msg)
				m.filterBranches()
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m.loadCommitsForCursor()
				}
			case "down", "j":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
					m.loadCommitsForCursor()
				}
			case "enter":
				if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
					branch := m.filtered[m.cursor]
					if branch.Name != m.currentBranch {
						m.selectedBranch = branch.Name
					}
				}
			case "/":
				m.filtering = true
				m.filter.Focus()
				return m, textinput.Blink
			case "esc":
				m.shouldGoBack = true
			}
		}
	}

	return m, nil
}

func (m *SwitchBranchModel) filterBranches() {
	query := strings.ToLower(m.filter.Value())
	if query == "" {
		m.filtered = m.branches
		m.cursor = 0
		return
	}

	m.filtered = []git.Branch{}
	for _, branch := range m.branches {
		if strings.Contains(strings.ToLower(branch.Name), query) {
			m.filtered = append(m.filtered, branch)
		}
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m *SwitchBranchModel) loadCommitsForCursor() {
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		branch := m.filtered[m.cursor].Name
		if _, ok := m.commits[branch]; !ok {
			m.branchForCommits = branch
		}
	}
}

// View renders the branch switch view
func (m SwitchBranchModel) View() string {
	// Calculate panel widths (40% branches, 60% commits)
	leftWidth := int(float64(m.width) * 0.4)
	rightWidth := m.width - leftWidth - 4 // Account for borders

	// Left panel: branches
	leftPanel := m.renderBranchList(leftWidth)

	// Right panel: commits
	rightPanel := m.renderCommitHistory(rightWidth)

	// Combine panels
	splitView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	// Help
	var help string
	if m.filtering {
		help = styles.HelpStyle.Render("↑/↓ navigate • Enter switch • Esc cancel filter")
	} else {
		help = styles.HelpStyle.Render("↑/↓ navigate • Enter switch • / filter • Esc back")
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			splitView,
			"",
			help,
		),
	)
}

func (m SwitchBranchModel) renderBranchList(width int) string {
	title := styles.TitleStyle.Render("Branches")

	// Filter input
	filterStyle := styles.InputStyle
	if m.filtering {
		filterStyle = styles.FocusedInputStyle
	}
	filterInput := filterStyle.Width(width - 6).Render(m.filter.View())

	// Branch list
	maxItems := m.height - 12
	if maxItems < 5 {
		maxItems = 5
	}

	var items []string
	startIdx := 0
	if m.cursor >= maxItems {
		startIdx = m.cursor - maxItems + 1
	}

	for i := startIdx; i < len(m.filtered) && i < startIdx+maxItems; i++ {
		branch := m.filtered[i]
		var line string

		if branch.IsCurrent {
			line = styles.CurrentBranchStyle.Render("* " + branch.Name)
		} else if i == m.cursor {
			line = styles.SelectedItemStyle.Render("→ " + branch.Name)
		} else {
			line = styles.MenuItemStyle.Render("  " + branch.Name)
		}

		// Add relative time if available
		if branch.LastCommit != nil {
			relTime := formatRelativeTime(branch.LastCommit.Date)
			line += " " + styles.SubtleText.Render(relTime)
		}

		items = append(items, truncateStringWidth(line, width-4))
	}

	branchList := strings.Join(items, "\n")

	if len(m.filtered) == 0 {
		branchList = styles.MutedText.Render("No branches found")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		filterInput,
		"",
		branchList,
	)

	return styles.FocusedPanelStyle.
		Width(width).
		Height(m.height - 6).
		Render(content)
}

func (m SwitchBranchModel) renderCommitHistory(width int) string {
	title := styles.TitleStyle.Render("Recent Commits")

	var branchName string
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		branchName = m.filtered[m.cursor].Name
	}

	subtitle := styles.MutedText.Render("on " + branchName)

	commits, ok := m.commits[branchName]
	if !ok || len(commits) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
			"",
			styles.MutedText.Render("Loading..."),
		)

		return styles.PanelStyle.
			Width(width).
			Height(m.height - 6).
			Render(content)
	}

	maxItems := m.height - 12
	if maxItems < 5 {
		maxItems = 5
	}

	var items []string
	for i := 0; i < len(commits) && i < maxItems; i++ {
		commit := commits[i]

		hash := styles.CommitHashStyle.Render(commit.ShortHash)
		msg := styles.CommitMsgStyle.Render(truncateString(commit.Message, width-15))
		author := styles.SubtleText.Render(commit.Author)

		line := fmt.Sprintf("%s %s", hash, msg)
		items = append(items, line)
		items = append(items, "       "+author)
	}

	commitList := strings.Join(items, "\n")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		commitList,
	)

	return styles.PanelStyle.
		Width(width).
		Height(m.height - 6).
		Render(content)
}

func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	default:
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
}

func truncateStringWidth(s string, maxWidth int) string {
	// Simple truncation - doesn't account for ANSI codes
	if len(s) <= maxWidth {
		return s
	}
	return s[:maxWidth-3] + "..."
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
