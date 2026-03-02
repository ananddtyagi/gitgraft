package git

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
)

// Branch represents a git branch
type Branch struct {
	Name      string
	IsCurrent bool
	IsRemote  bool
	LastCommit *CommitInfo
}

// CommitInfo holds basic commit information
type CommitInfo struct {
	Hash      string
	ShortHash string
	Message   string
	Author    string
	Date      time.Time
}

// ListBranches returns all local branches sorted by most recent commit
func (c *Client) ListBranches() ([]Branch, error) {
	branches := []Branch{}

	// Get current branch
	head, err := c.repo.Head()
	if err != nil {
		return nil, err
	}
	currentBranch := ""
	if head.Name().IsBranch() {
		currentBranch = head.Name().Short()
	}

	// List all branches
	refs, err := c.repo.Branches()
	if err != nil {
		return nil, err
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()

		// Get last commit for this branch
		commit, err := c.repo.CommitObject(ref.Hash())
		if err != nil {
			return nil
		}

		branches = append(branches, Branch{
			Name:      name,
			IsCurrent: name == currentBranch,
			IsRemote:  false,
			LastCommit: &CommitInfo{
				Hash:      commit.Hash.String(),
				ShortHash: commit.Hash.String()[:7],
				Message:   strings.Split(commit.Message, "\n")[0],
				Author:    commit.Author.Name,
				Date:      commit.Author.When,
			},
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort by most recent commit
	sort.Slice(branches, func(i, j int) bool {
		if branches[i].IsCurrent {
			return true
		}
		if branches[j].IsCurrent {
			return false
		}
		if branches[i].LastCommit == nil {
			return false
		}
		if branches[j].LastCommit == nil {
			return true
		}
		return branches[i].LastCommit.Date.After(branches[j].LastCommit.Date)
	})

	return branches, nil
}

// GetCurrentBranch returns the current branch name
func (c *Client) GetCurrentBranch() (string, error) {
	head, err := c.repo.Head()
	if err != nil {
		return "", err
	}

	if head.Name().IsBranch() {
		return head.Name().Short(), nil
	}

	// Detached HEAD
	return head.Hash().String()[:7], nil
}

// CreateBranch creates a new branch from the specified base
func (c *Client) CreateBranch(name, baseBranch string) error {
	// Check if branch already exists
	branches, err := c.ListBranches()
	if err != nil {
		return err
	}

	for _, b := range branches {
		if b.Name == name {
			return NewBranchExistsError(name)
		}
	}

	// Use git CLI for branch creation
	_, err = c.RunGitCommand("checkout", "-b", name, baseBranch)
	return err
}

// CreateBranchFromCurrent creates a new branch from the current HEAD
func (c *Client) CreateBranchFromCurrent(name string) error {
	// Check if branch already exists
	branches, err := c.ListBranches()
	if err != nil {
		return err
	}

	for _, b := range branches {
		if b.Name == name {
			return NewBranchExistsError(name)
		}
	}

	// Use git CLI for branch creation
	_, err = c.RunGitCommand("checkout", "-b", name)
	return err
}

// SwitchBranch switches to the specified branch
func (c *Client) SwitchBranch(name string) error {
	// Check for uncommitted changes
	clean, err := c.IsClean()
	if err != nil {
		return err
	}

	if !clean {
		return NewUncommittedChangesError("checkout", "")
	}

	_, err = c.RunGitCommand("checkout", name)
	return err
}

// DeleteBranch deletes a branch
func (c *Client) DeleteBranch(name string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}
	_, err := c.RunGitCommand("branch", flag, name)
	return err
}

// GetBranchCommits returns the commit history for a branch
func (c *Client) GetBranchCommits(branchName string, limit int) ([]CommitInfo, error) {
	ref, err := c.repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return nil, err
	}

	commit, err := c.repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	commits := []CommitInfo{}
	iter := commit
	for i := 0; i < limit && iter != nil; i++ {
		commits = append(commits, CommitInfo{
			Hash:      iter.Hash.String(),
			ShortHash: iter.Hash.String()[:7],
			Message:   strings.Split(iter.Message, "\n")[0],
			Author:    iter.Author.Name,
			Date:      iter.Author.When,
		})

		if iter.NumParents() == 0 {
			break
		}
		iter, err = iter.Parent(0)
		if err != nil {
			break
		}
	}

	return commits, nil
}

// RenameBranch renames a branch
func (c *Client) RenameBranch(oldName, newName string) error {
	_, err := c.RunGitCommand("branch", "-m", oldName, newName)
	return err
}

// GetDefaultBranch tries to determine the default branch (main or master)
func (c *Client) GetDefaultBranch() string {
	branches, err := c.ListBranches()
	if err != nil {
		return "main"
	}

	for _, b := range branches {
		if b.Name == "main" {
			return "main"
		}
	}

	for _, b := range branches {
		if b.Name == "master" {
			return "master"
		}
	}

	// Check remote default
	output, err := c.RunGitCommand("symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil {
		parts := strings.Split(output, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}

	return "main"
}

// FilterBranches returns branches matching the filter string
func (c *Client) FilterBranches(branches []Branch, filter string) []Branch {
	if filter == "" {
		return branches
	}

	filter = strings.ToLower(filter)
	filtered := []Branch{}

	for _, b := range branches {
		if strings.Contains(strings.ToLower(b.Name), filter) {
			filtered = append(filtered, b)
		}
	}

	return filtered
}

// Stash stashes the current changes
func (c *Client) Stash(message string) error {
	if message != "" {
		_, err := c.RunGitCommand("stash", "push", "-m", message)
		return err
	}
	_, err := c.RunGitCommand("stash")
	return err
}

// StashPop pops the most recent stash
func (c *Client) StashPop() error {
	_, err := c.RunGitCommand("stash", "pop")
	return err
}

// DiscardChanges discards all uncommitted changes
func (c *Client) DiscardChanges() error {
	_, err := c.RunGitCommand("checkout", "--", ".")
	if err != nil {
		return err
	}
	// Also clean untracked files
	_, err = c.RunGitCommand("clean", "-fd")
	return err
}

// FormatBranchDisplay formats a branch for display
func (b Branch) FormatBranchDisplay() string {
	prefix := "  "
	if b.IsCurrent {
		prefix = "* "
	}

	commitInfo := ""
	if b.LastCommit != nil {
		commitInfo = fmt.Sprintf(" (%s) %s", b.LastCommit.ShortHash, truncateString(b.LastCommit.Message, 40))
	}

	return fmt.Sprintf("%s%s%s", prefix, b.Name, commitInfo)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
