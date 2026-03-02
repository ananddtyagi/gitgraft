package views

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anandtyagi/gitgraft/internal/git"
	"github.com/anandtyagi/gitgraft/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CommitSection represents which section is focused
type CommitSection int

const (
	SectionFiles CommitSection = iota
	SectionMessage
	SectionRegex
)

// CommitModel handles the commit view
type CommitModel struct {
	width        int
	height       int
	files        []git.FileEntry
	selected     map[string]bool
	cursor       int
	section      CommitSection
	message      textarea.Model
	regexInput   textinput.Model
	showRegex    bool
	shouldCommit bool
	shouldGoBack bool
	gitClient    *git.Client
}

// NewCommit creates a new commit view
func NewCommit(client *git.Client) CommitModel {
	ta := textarea.New()
	ta.Placeholder = "Commit message..."
	ta.SetWidth(50)
	ta.SetHeight(3)
	ta.CharLimit = 1000

	regex := textinput.New()
	regex.Placeholder = "Regex pattern (e.g. *.go)"
	regex.CharLimit = 100
	regex.Width = 40

	return CommitModel{
		selected:  make(map[string]bool),
		message:   ta,
		regexInput: regex,
		gitClient: client,
	}
}

// SetSize sets the view dimensions
func (m CommitModel) SetSize(width, height int) CommitModel {
	m.width = width
	m.height = height
	// Only set if initialized (message has a non-nil viewport)
	if width > 10 && m.selected != nil {
		m.message.SetWidth(min(60, width-10))
		m.regexInput.Width = min(40, width-10)
	}
	return m
}

// SetFiles sets the file list
func (m CommitModel) SetFiles(files []git.FileEntry) CommitModel {
	m.files = files
	// Auto-select already staged files
	for _, f := range files {
		if f.IsStaged {
			m.selected[f.Path] = true
		}
	}
	return m
}

// ShouldCommit returns true if commit should proceed
func (m CommitModel) ShouldCommit() bool {
	return m.shouldCommit
}

// ShouldGoBack returns true if should navigate back
func (m CommitModel) ShouldGoBack() bool {
	return m.shouldGoBack
}

// GetCommitInfo returns selected files and message
func (m CommitModel) GetCommitInfo() ([]string, string) {
	var files []string
	for path, selected := range m.selected {
		if selected {
			files = append(files, path)
		}
	}
	return files, strings.TrimSpace(m.message.Value())
}

// Reset resets the view state
func (m CommitModel) Reset() CommitModel {
	m.shouldCommit = false
	m.shouldGoBack = false
	m.selected = make(map[string]bool)
	m.message.Reset()
	m.cursor = 0
	m.section = SectionFiles
	m.showRegex = false
	return m
}

// Update handles messages
func (m CommitModel) Update(msg tea.Msg) (CommitModel, tea.Cmd) {
	m.shouldCommit = false
	m.shouldGoBack = false

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle regex input mode
		if m.showRegex {
			switch msg.String() {
			case "enter":
				// Apply regex selection
				pattern := m.regexInput.Value()
				if pattern != "" {
					m.applyRegexSelection(pattern)
				}
				m.showRegex = false
				m.regexInput.Blur()
				m.regexInput.SetValue("")
				m.section = SectionFiles
			case "esc":
				m.showRegex = false
				m.regexInput.Blur()
				m.regexInput.SetValue("")
				m.section = SectionFiles
			default:
				var cmd tea.Cmd
				m.regexInput, cmd = m.regexInput.Update(msg)
				return m, cmd
			}
			return m, nil
		}

		switch m.section {
		case SectionFiles:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.files)-1 {
					m.cursor++
				}
			case " ":
				// Toggle selection
				if len(m.files) > 0 && m.cursor < len(m.files) {
					path := m.files[m.cursor].Path
					m.selected[path] = !m.selected[path]
				}
			case "ctrl+a":
				// Select all
				allSelected := m.areAllSelected()
				for _, f := range m.files {
					m.selected[f.Path] = !allSelected
				}
			case "ctrl+r":
				// Open regex input
				m.showRegex = true
				m.section = SectionRegex
				m.regexInput.Focus()
				return m, textinput.Blink
			case "tab":
				m.section = SectionMessage
				m.message.Focus()
				return m, textarea.Blink
			case "enter":
				// Check if ready to commit
				if m.canCommit() {
					m.shouldCommit = true
				} else if len(m.selected) == 0 {
					// Jump to file selection
				} else {
					// Jump to message
					m.section = SectionMessage
					m.message.Focus()
					return m, textarea.Blink
				}
			case "esc":
				m.shouldGoBack = true
			}

		case SectionMessage:
			switch msg.String() {
			case "tab":
				m.section = SectionFiles
				m.message.Blur()
			case "ctrl+s", "ctrl+enter":
				if m.canCommit() {
					m.shouldCommit = true
				}
			case "esc":
				if m.message.Value() == "" {
					m.section = SectionFiles
					m.message.Blur()
				} else {
					m.shouldGoBack = true
				}
			default:
				var cmd tea.Cmd
				m.message, cmd = m.message.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m *CommitModel) applyRegexSelection(pattern string) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	for _, f := range m.files {
		if re.MatchString(f.Path) {
			m.selected[f.Path] = true
		}
	}
}

func (m CommitModel) areAllSelected() bool {
	if len(m.files) == 0 {
		return false
	}
	for _, f := range m.files {
		if !m.selected[f.Path] {
			return false
		}
	}
	return true
}

func (m CommitModel) canCommit() bool {
	hasSelected := false
	for _, selected := range m.selected {
		if selected {
			hasSelected = true
			break
		}
	}
	return hasSelected && strings.TrimSpace(m.message.Value()) != ""
}

func (m CommitModel) selectedCount() int {
	count := 0
	for _, selected := range m.selected {
		if selected {
			count++
		}
	}
	return count
}

// View renders the commit view
func (m CommitModel) View() string {
	title := styles.TitleStyle.Render("Create Commit")

	// Regex input overlay
	if m.showRegex {
		return m.renderRegexInput()
	}

	// File list
	fileList := m.renderFileList()

	// Commit message
	messageBox := m.renderMessageBox()

	// Status
	status := m.renderStatus()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		fileList,
		"",
		messageBox,
		"",
		status,
	)

	help := m.renderHelp()

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

func (m CommitModel) renderFileList() string {
	header := styles.LabelStyle.Render(fmt.Sprintf("Files (%d selected)", m.selectedCount()))
	if m.section == SectionFiles {
		header = styles.FocusedLabelStyle.Render(fmt.Sprintf("Files (%d selected)", m.selectedCount()))
	}

	if len(m.files) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			styles.MutedText.Render("No changes to commit"),
		)
	}

	maxItems := 10
	var items []string
	startIdx := 0
	if m.cursor >= maxItems {
		startIdx = m.cursor - maxItems + 1
	}

	for i := startIdx; i < len(m.files) && i < startIdx+maxItems; i++ {
		file := m.files[i]
		items = append(items, m.renderFileItem(file, i))
	}

	fileList := strings.Join(items, "\n")

	// Scroll indicator
	if len(m.files) > maxItems {
		scrollInfo := styles.SubtleText.Render(
			fmt.Sprintf("(%d-%d of %d)", startIdx+1, min(startIdx+maxItems, len(m.files)), len(m.files)),
		)
		fileList += "\n" + scrollInfo
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		fileList,
	)
}

func (m CommitModel) renderFileItem(file git.FileEntry, index int) string {
	// Checkbox
	checkbox := "[ ]"
	if m.selected[file.Path] {
		checkbox = styles.CheckedStyle.Render("[✓]")
	}

	// Status indicator
	var statusStyle lipgloss.Style
	switch file.Status {
	case git.StatusModified:
		statusStyle = styles.ModifiedStyle
	case git.StatusAdded:
		statusStyle = styles.StagedStyle
	case git.StatusDeleted:
		statusStyle = styles.DeletedStyle
	case git.StatusUntracked:
		statusStyle = styles.UntrackedStyle
	default:
		statusStyle = styles.MutedText
	}
	status := statusStyle.Render(file.StatusString())

	// File path
	path := file.Path
	if len(path) > 50 {
		path = "..." + path[len(path)-47:]
	}

	// Cursor
	prefix := " "
	if index == m.cursor && m.section == SectionFiles {
		prefix = styles.PastelBlue.Render("→")
	}

	return fmt.Sprintf("%s %s %s %s", prefix, checkbox, status, path)
}

func (m CommitModel) renderMessageBox() string {
	label := styles.LabelStyle.Render("Commit message")
	if m.section == SectionMessage {
		label = styles.FocusedLabelStyle.Render("Commit message")
	}

	msgStyle := styles.InputStyle
	if m.section == SectionMessage {
		msgStyle = styles.FocusedInputStyle
	}

	msgBox := msgStyle.
		Width(m.message.Width() + 2).
		Render(m.message.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		label,
		msgBox,
	)
}

func (m CommitModel) renderStatus() string {
	if m.canCommit() {
		return styles.SuccessStyle.Render("Ready to commit")
	}

	if m.selectedCount() == 0 {
		return styles.WarningStyle.Render("Select files to commit")
	}

	return styles.WarningStyle.Render("Enter a commit message")
}

func (m CommitModel) renderHelp() string {
	switch m.section {
	case SectionFiles:
		return styles.HelpStyle.Render("Space toggle • Ctrl+A all • Ctrl+R regex • Tab message • Enter commit • Esc back")
	case SectionMessage:
		return styles.HelpStyle.Render("Tab files • Ctrl+S commit • Esc back")
	default:
		return ""
	}
}

func (m CommitModel) renderRegexInput() string {
	title := styles.TitleStyle.Render("Select by Pattern")
	subtitle := styles.MutedText.Render("Enter a regex to select matching files")

	inputStyle := styles.FocusedInputStyle.Width(m.regexInput.Width + 2)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		"",
		inputStyle.Render(m.regexInput.View()),
		"",
		styles.HelpStyle.Render("Enter apply • Esc cancel"),
	)

	panel := styles.PanelStyle.
		Width(min(60, m.width-4)).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		panel,
	)
}
