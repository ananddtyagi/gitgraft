package app

import (
	"github.com/anandtyagi/gitgraft/internal/config"
	"github.com/anandtyagi/gitgraft/internal/git"
	"github.com/anandtyagi/gitgraft/internal/tui/styles"
	"github.com/anandtyagi/gitgraft/internal/tui/views"
	"github.com/anandtyagi/gitgraft/internal/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Re-export state constants for convenience
const (
	StateOnboarding    = types.StateOnboarding
	StateCommandCenter = types.StateCommandCenter
	StateNewBranch     = types.StateNewBranch
	StateSwitchBranch  = types.StateSwitchBranch
	StateCommit        = types.StateCommit
	StateError         = types.StateError
	StatePushPrompt    = types.StatePushPrompt
)

// AppState is an alias for types.AppState
type AppState = types.AppState

// Model is the main application model
type Model struct {
	// Core
	state       AppState
	prevState   AppState
	width       int
	height      int
	keys        KeyMap
	config      *config.Config
	gitClient   *git.Client
	gitError    *git.GitError

	// Views
	onboarding    views.OnboardingModel
	commandCenter views.CommandCenterModel
	newBranch     views.NewBranchModel
	switchBranch  views.SwitchBranchModel
	commit        views.CommitModel
	errorView     views.ErrorModel

	// Flags
	ready         bool
	pushAfterCommit bool
}

// New creates a new application model
func New() Model {
	return Model{
		state:  StateCommandCenter, // Default, will be changed if first run
		keys:   DefaultKeyMap(),
		ready:  false,
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadConfig(),
		m.initGitClient(),
	)
}

func (m Model) loadConfig() tea.Cmd {
	return func() tea.Msg {
		cfg, _ := config.Load()
		return ConfigLoadedMsg{
			FirstRun: cfg.FirstRun,
			Alias:    cfg.Alias,
		}
	}
}

func (m Model) initGitClient() tea.Cmd {
	return func() tea.Msg {
		client, err := git.NewClient()
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "init",
				Err:     err,
				Message: err.Error(),
			}}
		}
		return gitClientInitMsg{client: client}
	}
}

type gitClientInitMsg struct {
	client *git.Client
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.updateViewSizes()

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case gitClientInitMsg:
		m.gitClient = msg.client
		m.ready = true
		// Initialize views that need git client
		if m.gitClient != nil {
			m.commandCenter = views.NewCommandCenter(m.gitClient)
			m.newBranch = views.NewNewBranch(m.gitClient)
			m.switchBranch = views.NewSwitchBranch(m.gitClient)
			m.commit = views.NewCommit(m.gitClient)
			m = m.updateViewSizes()
		}

	case ConfigLoadedMsg:
		cfg, _ := config.Load()
		m.config = cfg
		if msg.FirstRun {
			m.state = StateOnboarding
			m.onboarding = views.NewOnboarding()
			m.onboarding = m.onboarding.SetSize(m.width, m.height)
		}

	case OnboardingCompleteMsg:
		m.state = StateCommandCenter
		if m.config != nil {
			m.config.MarkOnboardingComplete()
			config.Save(m.config)
		}

	case NavigateMsg:
		m.prevState = m.state
		m.state = msg.State
		// Initialize view data when navigating
		switch msg.State {
		case StateSwitchBranch:
			cmds = append(cmds, LoadBranches(m.gitClient))
		case StateCommit:
			cmds = append(cmds, LoadStatus(m.gitClient))
		}

	case NavigateBackMsg:
		m.state = StateCommandCenter

	case ErrorMsg:
		m.gitError = msg.Err
		m.prevState = m.state
		m.state = StateError
		m.errorView = views.NewError(msg.Err)
		m.errorView = m.errorView.SetSize(m.width, m.height)

	case BranchCreatedMsg:
		// After creating branch, go back to command center
		m.state = StateCommandCenter
		m.commandCenter = m.commandCenter.SetMessage("Branch '" + msg.Name + "' created and checked out")

	case BranchSwitchedMsg:
		m.state = StateCommandCenter
		m.commandCenter = m.commandCenter.SetMessage("Switched to branch '" + msg.Name + "'")

	case CommitSuccessMsg:
		if m.config != nil && m.config.PushAfterCommit && m.gitClient.HasRemote() {
			m.state = StatePushPrompt
			m.pushAfterCommit = true
		} else {
			m.state = StateCommandCenter
			m.commandCenter = m.commandCenter.SetMessage("Commit created successfully")
		}

	case PushSuccessMsg:
		m.state = StateCommandCenter
		m.commandCenter = m.commandCenter.SetMessage("Changes pushed successfully")

	case BranchesLoadedMsg:
		m.switchBranch = m.switchBranch.SetBranches(msg.Branches, msg.Current)
		// Load commits for the first branch
		if len(msg.Branches) > 0 {
			cmds = append(cmds, LoadCommits(m.gitClient, msg.Branches[0].Name, 20))
		}

	case CommitsLoadedMsg:
		m.switchBranch = m.switchBranch.SetCommits(msg.Branch, msg.Commits)

	case StatusLoadedMsg:
		m.commit = m.commit.SetFiles(msg.Files)
	}

	// Route to current view's update
	var cmd tea.Cmd
	switch m.state {
	case StateOnboarding:
		m.onboarding, cmd = m.onboarding.Update(msg)
		cmds = append(cmds, cmd)
		// Check if onboarding is complete
		if m.onboarding.Done() {
			cmds = append(cmds, func() tea.Msg { return OnboardingCompleteMsg{} })
		}

	case StateCommandCenter:
		m.commandCenter, cmd = m.commandCenter.Update(msg)
		cmds = append(cmds, cmd)
		// Check for navigation requests
		if nav := m.commandCenter.NavigateTo(); nav != nil {
			cmds = append(cmds, NavigateTo(*nav))
		}

	case StateNewBranch:
		m.newBranch, cmd = m.newBranch.Update(msg)
		cmds = append(cmds, cmd)
		// Check if branch creation requested
		if m.newBranch.ShouldCreate() {
			name, base := m.newBranch.GetBranchInfo()
			cmds = append(cmds, CreateBranch(m.gitClient, name, base))
			m.newBranch = m.newBranch.Reset()
		}
		// Check for back navigation
		if m.newBranch.ShouldGoBack() {
			cmds = append(cmds, NavigateBack())
		}

	case StateSwitchBranch:
		m.switchBranch, cmd = m.switchBranch.Update(msg)
		cmds = append(cmds, cmd)
		// Check if branch switch requested
		if branch := m.switchBranch.SelectedBranch(); branch != "" {
			cmds = append(cmds, SwitchBranch(m.gitClient, branch))
			m.switchBranch = m.switchBranch.Reset()
		}
		// Check for back navigation
		if m.switchBranch.ShouldGoBack() {
			cmds = append(cmds, NavigateBack())
		}
		// Load commits when selection changes
		if newBranch := m.switchBranch.BranchForCommits(); newBranch != "" {
			cmds = append(cmds, LoadCommits(m.gitClient, newBranch, 20))
		}

	case StateCommit:
		m.commit, cmd = m.commit.Update(msg)
		cmds = append(cmds, cmd)
		// Check if commit requested
		if m.commit.ShouldCommit() {
			files, message := m.commit.GetCommitInfo()
			if len(files) > 0 && message != "" {
				cmds = append(cmds, StageAndCommit(m.gitClient, files, message))
				m.commit = m.commit.Reset()
			}
		}
		// Check for back navigation
		if m.commit.ShouldGoBack() {
			cmds = append(cmds, NavigateBack())
		}

	case StateError:
		m.errorView, cmd = m.errorView.Update(msg)
		cmds = append(cmds, cmd)
		// Check for action selection
		if action := m.errorView.SelectedAction(); action != nil {
			switch action.ActionType {
			case git.ActionCancel:
				m.state = StateCommandCenter
			case git.ActionStash:
				cmds = append(cmds, StashChanges(m.gitClient, "graft: auto-stash"))
				m.state = m.prevState
			case git.ActionDiscard:
				cmds = append(cmds, DiscardChanges(m.gitClient))
				m.state = m.prevState
			case git.ActionCommit:
				m.state = StateCommit
				cmds = append(cmds, LoadStatus(m.gitClient))
			default:
				m.state = StateCommandCenter
			}
		}

	case StatePushPrompt:
		// Handle push prompt
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "y", "Y", "enter":
				cmds = append(cmds, Push(m.gitClient))
			case "n", "N", "esc":
				m.state = StateCommandCenter
				m.commandCenter = m.commandCenter.SetMessage("Commit created (not pushed)")
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateViewSizes() Model {
	if m.width == 0 || m.height == 0 {
		return m
	}

	m.onboarding = m.onboarding.SetSize(m.width, m.height)
	m.commandCenter = m.commandCenter.SetSize(m.width, m.height)
	m.newBranch = m.newBranch.SetSize(m.width, m.height)
	m.switchBranch = m.switchBranch.SetSize(m.width, m.height)
	m.commit = m.commit.SetSize(m.width, m.height)
	m.errorView = m.errorView.SetSize(m.width, m.height)

	return m
}

// View renders the application
func (m Model) View() string {
	if !m.ready && m.state != StateOnboarding && m.state != StateError {
		return m.renderLoading()
	}

	var content string

	switch m.state {
	case StateOnboarding:
		content = m.onboarding.View()
	case StateCommandCenter:
		content = m.commandCenter.View()
	case StateNewBranch:
		content = m.newBranch.View()
	case StateSwitchBranch:
		content = m.switchBranch.View()
	case StateCommit:
		content = m.commit.View()
	case StateError:
		content = m.errorView.View()
	case StatePushPrompt:
		content = m.renderPushPrompt()
	default:
		content = "Unknown state"
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(content)
}

func (m Model) renderLoading() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			styles.Logo(),
			"",
			styles.MutedText.Render("Loading..."),
		),
	)
}

func (m Model) renderPushPrompt() string {
	prompt := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.TitleStyle.Render("Commit created successfully!"),
		"",
		styles.PrimaryText.Render("Push changes to remote?"),
		"",
		styles.HelpStyle.Render("[Y]es  [N]o"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		styles.PanelStyle.Render(prompt),
	)
}
