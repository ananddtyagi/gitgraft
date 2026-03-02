package git

import (
	"fmt"
	"regexp"
	"strings"
)

// StageFile stages a single file
func (c *Client) StageFile(path string) error {
	_, err := c.RunGitCommand("add", path)
	return err
}

// StageFiles stages multiple files
func (c *Client) StageFiles(paths []string) error {
	args := append([]string{"add"}, paths...)
	_, err := c.RunGitCommand(args...)
	return err
}

// UnstageFile unstages a single file
func (c *Client) UnstageFile(path string) error {
	_, err := c.RunGitCommand("reset", "HEAD", path)
	return err
}

// UnstageFiles unstages multiple files
func (c *Client) UnstageFiles(paths []string) error {
	args := append([]string{"reset", "HEAD"}, paths...)
	_, err := c.RunGitCommand(args...)
	return err
}

// StageAll stages all changes
func (c *Client) StageAll() error {
	_, err := c.RunGitCommand("add", "-A")
	return err
}

// UnstageAll unstages all changes
func (c *Client) UnstageAll() error {
	_, err := c.RunGitCommand("reset", "HEAD")
	return err
}

// StageByRegex stages files matching a regex pattern
func (c *Client) StageByRegex(pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	files, err := c.GetUnstagedFiles()
	if err != nil {
		return nil, err
	}

	matched := []string{}
	for _, f := range files {
		if re.MatchString(f.Path) {
			matched = append(matched, f.Path)
		}
	}

	if len(matched) > 0 {
		if err := c.StageFiles(matched); err != nil {
			return nil, err
		}
	}

	return matched, nil
}

// Commit creates a new commit with the staged changes
// Uses git CLI to support hooks
func (c *Client) Commit(message string) error {
	hasStagedChanges, err := c.HasStagedChanges()
	if err != nil {
		return err
	}

	if !hasStagedChanges {
		return &GitError{
			Op:      "commit",
			Err:     ErrNothingToCommit,
			Message: "No staged changes to commit",
		}
	}

	_, err = c.RunGitCommand("commit", "-m", message)
	return err
}

// CommitWithBody creates a commit with a subject and body
func (c *Client) CommitWithBody(subject, body string) error {
	hasStagedChanges, err := c.HasStagedChanges()
	if err != nil {
		return err
	}

	if !hasStagedChanges {
		return &GitError{
			Op:      "commit",
			Err:     ErrNothingToCommit,
			Message: "No staged changes to commit",
		}
	}

	message := subject
	if body != "" {
		message = subject + "\n\n" + body
	}

	_, err = c.RunGitCommand("commit", "-m", message)
	return err
}

// CommitNoVerify creates a commit skipping hooks
func (c *Client) CommitNoVerify(message string) error {
	hasStagedChanges, err := c.HasStagedChanges()
	if err != nil {
		return err
	}

	if !hasStagedChanges {
		return &GitError{
			Op:      "commit",
			Err:     ErrNothingToCommit,
			Message: "No staged changes to commit",
		}
	}

	_, err = c.RunGitCommand("commit", "-m", message, "--no-verify")
	return err
}

// AmendCommit amends the last commit with staged changes
func (c *Client) AmendCommit(message string) error {
	if message != "" {
		_, err := c.RunGitCommand("commit", "--amend", "-m", message)
		return err
	}
	_, err := c.RunGitCommand("commit", "--amend", "--no-edit")
	return err
}

// Push pushes to the remote
func (c *Client) Push() error {
	_, err := c.RunGitCommand("push")
	return err
}

// PushWithUpstream pushes and sets upstream tracking
func (c *Client) PushWithUpstream() error {
	branch, err := c.GetCurrentBranch()
	if err != nil {
		return err
	}
	_, err = c.RunGitCommand("push", "-u", "origin", branch)
	return err
}

// ForcePush force pushes with lease (safer)
func (c *Client) ForcePush() error {
	_, err := c.RunGitCommand("push", "--force-with-lease")
	return err
}

// Pull pulls changes from remote
func (c *Client) Pull() error {
	_, err := c.RunGitCommand("pull")
	return err
}

// PullRebase pulls with rebase
func (c *Client) PullRebase() error {
	_, err := c.RunGitCommand("pull", "--rebase")
	return err
}

// GetDiff returns the diff of staged changes
func (c *Client) GetDiff() (string, error) {
	return c.RunGitCommand("diff", "--cached")
}

// GetFileDiff returns the diff for a specific file
func (c *Client) GetFileDiff(path string, staged bool) (string, error) {
	if staged {
		return c.RunGitCommand("diff", "--cached", path)
	}
	return c.RunGitCommand("diff", path)
}

// GetLastCommitMessage returns the message of the last commit
func (c *Client) GetLastCommitMessage() (string, error) {
	return c.RunGitCommand("log", "-1", "--pretty=%B")
}

// GetRecentCommitMessages returns recent commit messages for style reference
func (c *Client) GetRecentCommitMessages(count int) ([]string, error) {
	output, err := c.RunGitCommand("log", fmt.Sprintf("-%d", count), "--pretty=%s")
	if err != nil {
		return nil, err
	}

	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

// HasUpstream returns true if the current branch has an upstream set
func (c *Client) HasUpstream() (bool, error) {
	branch, err := c.GetCurrentBranch()
	if err != nil {
		return false, err
	}

	_, err = c.RunGitCommand("config", "--get", fmt.Sprintf("branch.%s.remote", branch))
	return err == nil, nil
}
