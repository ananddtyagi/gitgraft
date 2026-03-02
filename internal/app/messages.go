package app

import (
	"github.com/anandtyagi/gitgraft/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

// Navigation messages
type (
	// NavigateMsg requests navigation to a specific view
	NavigateMsg struct {
		State AppState
	}

	// NavigateBackMsg requests navigation to the previous view
	NavigateBackMsg struct{}
)

// Git operation result messages
type (
	// BranchesLoadedMsg contains the loaded branches
	BranchesLoadedMsg struct {
		Branches []git.Branch
		Current  string
	}

	// StatusLoadedMsg contains the loaded file status
	StatusLoadedMsg struct {
		Files []git.FileEntry
	}

	// CommitsLoadedMsg contains loaded commit history
	CommitsLoadedMsg struct {
		Commits []git.CommitInfo
		Branch  string
	}

	// BranchCreatedMsg indicates a branch was created
	BranchCreatedMsg struct {
		Name string
	}

	// BranchSwitchedMsg indicates the branch was switched
	BranchSwitchedMsg struct {
		Name string
	}

	// CommitSuccessMsg indicates a successful commit
	CommitSuccessMsg struct {
		Hash    string
		Message string
	}

	// PushSuccessMsg indicates a successful push
	PushSuccessMsg struct{}

	// PushPromptMsg asks if user wants to push
	PushPromptMsg struct{}
)

// Error messages
type (
	// ErrorMsg wraps an error for display
	ErrorMsg struct {
		Err *git.GitError
	}

	// ErrorActionMsg indicates user selected an error action
	ErrorActionMsg struct {
		Action git.Action
	}
)

// Config messages
type (
	// ConfigLoadedMsg contains the loaded config
	ConfigLoadedMsg struct {
		FirstRun bool
		Alias    string
	}

	// AliasSetMsg indicates the alias was configured
	AliasSetMsg struct {
		Alias string
	}

	// OnboardingCompleteMsg indicates onboarding is done
	OnboardingCompleteMsg struct{}
)

// Navigation commands
func NavigateTo(state AppState) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{State: state}
	}
}

func NavigateBack() tea.Cmd {
	return func() tea.Msg {
		return NavigateBackMsg{}
	}
}

// Git commands
func LoadBranches(client *git.Client) tea.Cmd {
	return func() tea.Msg {
		branches, err := client.ListBranches()
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "list branches",
				Err:     err,
				Message: err.Error(),
			}}
		}

		current, _ := client.GetCurrentBranch()

		return BranchesLoadedMsg{
			Branches: branches,
			Current:  current,
		}
	}
}

func LoadStatus(client *git.Client) tea.Cmd {
	return func() tea.Msg {
		files, err := client.GetStatus()
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "status",
				Err:     err,
				Message: err.Error(),
			}}
		}

		return StatusLoadedMsg{Files: files}
	}
}

func LoadCommits(client *git.Client, branch string, limit int) tea.Cmd {
	return func() tea.Msg {
		commits, err := client.GetBranchCommits(branch, limit)
		if err != nil {
			return CommitsLoadedMsg{Commits: []git.CommitInfo{}, Branch: branch}
		}

		return CommitsLoadedMsg{
			Commits: commits,
			Branch:  branch,
		}
	}
}

func CreateBranch(client *git.Client, name, baseBranch string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if baseBranch != "" {
			err = client.CreateBranch(name, baseBranch)
		} else {
			err = client.CreateBranchFromCurrent(name)
		}

		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "create branch",
				Err:     err,
				Message: err.Error(),
			}}
		}

		return BranchCreatedMsg{Name: name}
	}
}

func SwitchBranch(client *git.Client, name string) tea.Cmd {
	return func() tea.Msg {
		err := client.SwitchBranch(name)
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "switch branch",
				Err:     err,
				Message: err.Error(),
			}}
		}

		return BranchSwitchedMsg{Name: name}
	}
}

func StageFiles(client *git.Client, paths []string) tea.Cmd {
	return func() tea.Msg {
		err := client.StageFiles(paths)
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "stage",
				Err:     err,
				Message: err.Error(),
			}}
		}

		// Reload status after staging
		files, _ := client.GetStatus()
		return StatusLoadedMsg{Files: files}
	}
}

func CreateCommit(client *git.Client, message string) tea.Cmd {
	return func() tea.Msg {
		err := client.Commit(message)
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "commit",
				Err:     err,
				Message: err.Error(),
			}}
		}

		return CommitSuccessMsg{Message: message}
	}
}

func Push(client *git.Client) tea.Cmd {
	return func() tea.Msg {
		// Check if upstream is set
		hasUpstream, _ := client.HasUpstream()
		var err error
		if hasUpstream {
			err = client.Push()
		} else {
			err = client.PushWithUpstream()
		}

		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "push",
				Err:     err,
				Message: err.Error(),
			}}
		}

		return PushSuccessMsg{}
	}
}

func StashChanges(client *git.Client, message string) tea.Cmd {
	return func() tea.Msg {
		err := client.Stash(message)
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "stash",
				Err:     err,
				Message: err.Error(),
			}}
		}
		return nil
	}
}

func DiscardChanges(client *git.Client) tea.Cmd {
	return func() tea.Msg {
		err := client.DiscardChanges()
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok {
				return ErrorMsg{Err: gitErr}
			}
			return ErrorMsg{Err: &git.GitError{
				Op:      "discard",
				Err:     err,
				Message: err.Error(),
			}}
		}
		return nil
	}
}
