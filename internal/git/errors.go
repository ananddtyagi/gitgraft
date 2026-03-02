package git

import (
	"errors"
	"fmt"
	"strings"
)

// Error types for git operations
var (
	ErrNotARepository     = errors.New("not a git repository")
	ErrBranchExists       = errors.New("branch already exists")
	ErrBranchNotFound     = errors.New("branch not found")
	ErrUncommittedChanges = errors.New("uncommitted changes exist")
	ErrNothingToCommit    = errors.New("nothing to commit")
	ErrPreCommitFailed    = errors.New("pre-commit hook failed")
	ErrPushRejected       = errors.New("push rejected by remote")
	ErrMergeConflict      = errors.New("merge conflict")
	ErrDetachedHead       = errors.New("detached HEAD state")
)

// GitError wraps a git error with additional context
type GitError struct {
	Op       string   // Operation that failed (e.g., "checkout", "commit")
	Err      error    // Underlying error
	Message  string   // User-friendly message
	Output   string   // Raw git output
	Actions  []Action // Suggested recovery actions
}

// Action represents a recovery action for an error
type Action struct {
	Label       string // Display label (e.g., "Stash changes")
	Description string // Longer description
	Command     string // Git command to run (if applicable)
	ActionType  ActionType
}

// ActionType represents the type of recovery action
type ActionType int

const (
	ActionStash ActionType = iota
	ActionCommit
	ActionDiscard
	ActionRename
	ActionDelete
	ActionRetry
	ActionSkip
	ActionPull
	ActionForce
	ActionCancel
)

func (e *GitError) Error() string {
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func (e *GitError) Unwrap() error {
	return e.Err
}

// NewUncommittedChangesError creates an error for uncommitted changes
func NewUncommittedChangesError(op string, output string) *GitError {
	return &GitError{
		Op:      op,
		Err:     ErrUncommittedChanges,
		Message: "You have uncommitted changes that would be overwritten",
		Output:  output,
		Actions: []Action{
			{
				Label:       "Stash changes",
				Description: "Save changes temporarily and restore later",
				Command:     "git stash",
				ActionType:  ActionStash,
			},
			{
				Label:       "Commit first",
				Description: "Commit current changes before proceeding",
				ActionType:  ActionCommit,
			},
			{
				Label:       "Discard changes",
				Description: "Permanently discard all uncommitted changes",
				Command:     "git checkout -- .",
				ActionType:  ActionDiscard,
			},
			{
				Label:       "Cancel",
				Description: "Return without making changes",
				ActionType:  ActionCancel,
			},
		},
	}
}

// NewBranchExistsError creates an error for when a branch already exists
func NewBranchExistsError(branchName string) *GitError {
	return &GitError{
		Op:      "create branch",
		Err:     ErrBranchExists,
		Message: fmt.Sprintf("Branch '%s' already exists", branchName),
		Actions: []Action{
			{
				Label:       "Choose different name",
				Description: "Enter a new branch name",
				ActionType:  ActionRename,
			},
			{
				Label:       "Delete existing branch",
				Description: fmt.Sprintf("Delete '%s' and create new", branchName),
				Command:     fmt.Sprintf("git branch -D %s", branchName),
				ActionType:  ActionDelete,
			},
			{
				Label:       "Switch to existing",
				Description: fmt.Sprintf("Switch to existing branch '%s'", branchName),
				Command:     fmt.Sprintf("git checkout %s", branchName),
				ActionType:  ActionCancel,
			},
		},
	}
}

// NewPreCommitFailedError creates an error for pre-commit hook failures
func NewPreCommitFailedError(output string) *GitError {
	return &GitError{
		Op:      "commit",
		Err:     ErrPreCommitFailed,
		Message: "Pre-commit hook failed",
		Output:  output,
		Actions: []Action{
			{
				Label:       "View output",
				Description: "See the full hook output",
				ActionType:  ActionRetry,
			},
			{
				Label:       "Retry commit",
				Description: "Fix issues and try again",
				ActionType:  ActionRetry,
			},
			{
				Label:       "Skip hooks",
				Description: "Commit without running hooks (not recommended)",
				Command:     "--no-verify",
				ActionType:  ActionSkip,
			},
			{
				Label:       "Cancel",
				Description: "Return without committing",
				ActionType:  ActionCancel,
			},
		},
	}
}

// NewPushRejectedError creates an error for rejected pushes
func NewPushRejectedError(output string) *GitError {
	return &GitError{
		Op:      "push",
		Err:     ErrPushRejected,
		Message: "Push was rejected by the remote",
		Output:  output,
		Actions: []Action{
			{
				Label:       "Pull and rebase",
				Description: "Fetch remote changes and rebase your commits",
				Command:     "git pull --rebase",
				ActionType:  ActionPull,
			},
			{
				Label:       "Pull and merge",
				Description: "Fetch remote changes and merge",
				Command:     "git pull",
				ActionType:  ActionPull,
			},
			{
				Label:       "Force push",
				Description: "Overwrite remote (dangerous!)",
				Command:     "git push --force-with-lease",
				ActionType:  ActionForce,
			},
			{
				Label:       "Cancel",
				Description: "Return without pushing",
				ActionType:  ActionCancel,
			},
		},
	}
}

// NewNotARepositoryError creates an error for non-repository directories
func NewNotARepositoryError(path string) *GitError {
	return &GitError{
		Op:      "init",
		Err:     ErrNotARepository,
		Message: fmt.Sprintf("'%s' is not a git repository", path),
		Actions: []Action{
			{
				Label:       "Initialize repository",
				Description: "Create a new git repository here",
				Command:     "git init",
				ActionType:  ActionRetry,
			},
			{
				Label:       "Exit",
				Description: "Close Git-Graft",
				ActionType:  ActionCancel,
			},
		},
	}
}

// ParseGitError parses git command output and returns an appropriate error
func ParseGitError(op string, output string, exitCode int) *GitError {
	output = strings.TrimSpace(output)
	lowerOutput := strings.ToLower(output)

	// Check for common error patterns
	switch {
	case strings.Contains(lowerOutput, "not a git repository"):
		return NewNotARepositoryError(".")

	case strings.Contains(lowerOutput, "already exists"):
		// Extract branch name if possible
		return &GitError{
			Op:      op,
			Err:     ErrBranchExists,
			Message: "Branch already exists",
			Output:  output,
		}

	case strings.Contains(lowerOutput, "uncommitted changes"),
		strings.Contains(lowerOutput, "local changes"),
		strings.Contains(lowerOutput, "would be overwritten"):
		return NewUncommittedChangesError(op, output)

	case strings.Contains(lowerOutput, "pre-commit hook"),
		strings.Contains(lowerOutput, "hook") && exitCode == 1:
		return NewPreCommitFailedError(output)

	case strings.Contains(lowerOutput, "rejected"),
		strings.Contains(lowerOutput, "failed to push"):
		return NewPushRejectedError(output)

	case strings.Contains(lowerOutput, "conflict"):
		return &GitError{
			Op:      op,
			Err:     ErrMergeConflict,
			Message: "Merge conflict detected",
			Output:  output,
		}

	case strings.Contains(lowerOutput, "nothing to commit"):
		return &GitError{
			Op:      op,
			Err:     ErrNothingToCommit,
			Message: "Nothing to commit, working tree clean",
			Output:  output,
		}

	default:
		return &GitError{
			Op:      op,
			Err:     errors.New(output),
			Message: output,
			Output:  output,
		}
	}
}
