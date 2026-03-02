package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Client provides git operations using a hybrid approach:
// - go-git for read operations (status, log, branches)
// - git CLI for write operations (commit, push) to support hooks
type Client struct {
	repo     *git.Repository
	repoPath string
}

// NewClient creates a new git client for the current directory
func NewClient() (*Client, error) {
	return NewClientAt(".")
}

// NewClientAt creates a new git client for the specified directory
func NewClientAt(path string) (*Client, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpenWithOptions(absPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil, NewNotARepositoryError(absPath)
		}
		return nil, err
	}

	// Get the actual repo root
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	repoPath := wt.Filesystem.Root()

	return &Client{
		repo:     repo,
		repoPath: repoPath,
	}, nil
}

// RepoPath returns the repository root path
func (c *Client) RepoPath() string {
	return c.repoPath
}

// Repository returns the underlying go-git repository
func (c *Client) Repository() *git.Repository {
	return c.repo
}

// RunGitCommand runs a git command and returns the output
func (c *Client) RunGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.repoPath
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return outputStr, ParseGitError(args[0], outputStr, exitErr.ExitCode())
		}
		return outputStr, err
	}

	return outputStr, nil
}

// RunGitCommandWithStdin runs a git command with stdin input
func (c *Client) RunGitCommandWithStdin(input string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.repoPath
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return outputStr, ParseGitError(args[0], outputStr, exitErr.ExitCode())
		}
		return outputStr, err
	}

	return outputStr, nil
}

// IsClean returns true if the working tree has no uncommitted changes
func (c *Client) IsClean() (bool, error) {
	wt, err := c.repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := wt.Status()
	if err != nil {
		return false, err
	}

	return status.IsClean(), nil
}

// HasRemote returns true if the repository has a remote named "origin"
func (c *Client) HasRemote() bool {
	_, err := c.repo.Remote("origin")
	return err == nil
}

// GetRemoteURL returns the URL for the origin remote
func (c *Client) GetRemoteURL() (string, error) {
	remote, err := c.repo.Remote("origin")
	if err != nil {
		return "", err
	}

	config := remote.Config()
	if len(config.URLs) > 0 {
		return config.URLs[0], nil
	}

	return "", nil
}
