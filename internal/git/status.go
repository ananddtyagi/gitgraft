package git

import (
	"github.com/go-git/go-git/v5"
)

// FileStatus represents the status of a file in the working tree
type FileStatus int

const (
	StatusUntracked FileStatus = iota
	StatusModified
	StatusAdded
	StatusDeleted
	StatusRenamed
	StatusCopied
	StatusUnmodified
	StatusIgnored
	StatusConflicted
)

// FileEntry represents a file with its status
type FileEntry struct {
	Path         string
	Status       FileStatus
	StagedStatus FileStatus
	IsStaged     bool
	OldPath      string // For renames
}

// GetStatus returns the current status of the working tree
func (c *Client) GetStatus() ([]FileEntry, error) {
	wt, err := c.repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := wt.Status()
	if err != nil {
		return nil, err
	}

	entries := []FileEntry{}

	for path, fileStatus := range status {
		entry := FileEntry{
			Path:         path,
			Status:       convertStatus(fileStatus.Worktree),
			StagedStatus: convertStatus(fileStatus.Staging),
			IsStaged:     fileStatus.Staging != git.Unmodified && fileStatus.Staging != git.Untracked,
		}

		// Handle untracked files
		if fileStatus.Worktree == git.Untracked {
			entry.Status = StatusUntracked
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// GetStagedFiles returns only staged files
func (c *Client) GetStagedFiles() ([]FileEntry, error) {
	allFiles, err := c.GetStatus()
	if err != nil {
		return nil, err
	}

	staged := []FileEntry{}
	for _, f := range allFiles {
		if f.IsStaged {
			staged = append(staged, f)
		}
	}

	return staged, nil
}

// GetUnstagedFiles returns only unstaged files (modified or untracked)
func (c *Client) GetUnstagedFiles() ([]FileEntry, error) {
	allFiles, err := c.GetStatus()
	if err != nil {
		return nil, err
	}

	unstaged := []FileEntry{}
	for _, f := range allFiles {
		if !f.IsStaged && f.Status != StatusUnmodified {
			unstaged = append(unstaged, f)
		}
	}

	return unstaged, nil
}

// GetChangedFiles returns all changed files (staged and unstaged)
func (c *Client) GetChangedFiles() ([]FileEntry, error) {
	return c.GetStatus()
}

// convertStatus converts go-git status to our FileStatus
func convertStatus(s git.StatusCode) FileStatus {
	switch s {
	case git.Untracked:
		return StatusUntracked
	case git.Modified:
		return StatusModified
	case git.Added:
		return StatusAdded
	case git.Deleted:
		return StatusDeleted
	case git.Renamed:
		return StatusRenamed
	case git.Copied:
		return StatusCopied
	case git.Unmodified:
		return StatusUnmodified
	case git.UpdatedButUnmerged:
		return StatusConflicted
	default:
		return StatusUnmodified
	}
}

// StatusString returns a short status string (like git status --short)
func (f FileEntry) StatusString() string {
	staged := statusChar(f.StagedStatus)
	unstaged := statusChar(f.Status)

	if f.Status == StatusUntracked {
		return "??"
	}

	return string(staged) + string(unstaged)
}

func statusChar(s FileStatus) rune {
	switch s {
	case StatusModified:
		return 'M'
	case StatusAdded:
		return 'A'
	case StatusDeleted:
		return 'D'
	case StatusRenamed:
		return 'R'
	case StatusCopied:
		return 'C'
	case StatusUntracked:
		return '?'
	case StatusIgnored:
		return '!'
	case StatusConflicted:
		return 'U'
	default:
		return ' '
	}
}

// StatusDescription returns a human-readable status description
func (f FileEntry) StatusDescription() string {
	switch f.Status {
	case StatusUntracked:
		return "untracked"
	case StatusModified:
		return "modified"
	case StatusAdded:
		return "new file"
	case StatusDeleted:
		return "deleted"
	case StatusRenamed:
		return "renamed"
	case StatusCopied:
		return "copied"
	case StatusConflicted:
		return "conflicted"
	default:
		return ""
	}
}

// HasChanges returns true if there are any changes
func (c *Client) HasChanges() (bool, error) {
	files, err := c.GetStatus()
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}

// HasStagedChanges returns true if there are staged changes
func (c *Client) HasStagedChanges() (bool, error) {
	files, err := c.GetStagedFiles()
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}
