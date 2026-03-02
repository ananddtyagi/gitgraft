package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShellType represents the type of shell
type ShellType int

const (
	ShellUnknown ShellType = iota
	ShellBash
	ShellZsh
	ShellFish
)

// DetectShell detects the current shell
func DetectShell() ShellType {
	shell := os.Getenv("SHELL")

	switch {
	case strings.Contains(shell, "zsh"):
		return ShellZsh
	case strings.Contains(shell, "bash"):
		return ShellBash
	case strings.Contains(shell, "fish"):
		return ShellFish
	default:
		return ShellUnknown
	}
}

// ShellConfigPath returns the config file path for the given shell
func ShellConfigPath(shellType ShellType) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch shellType {
	case ShellZsh:
		return filepath.Join(home, ".zshrc"), nil
	case ShellBash:
		// Check for .bashrc first, then .bash_profile
		bashrc := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(bashrc); err == nil {
			return bashrc, nil
		}
		return filepath.Join(home, ".bash_profile"), nil
	case ShellFish:
		return filepath.Join(home, ".config", "fish", "config.fish"), nil
	default:
		return "", fmt.Errorf("unknown shell type")
	}
}

// AliasCommand returns the alias command for the given shell
func AliasCommand(shellType ShellType, alias, command string) string {
	switch shellType {
	case ShellFish:
		return fmt.Sprintf("alias %s='%s'", alias, command)
	default:
		return fmt.Sprintf("alias %s='%s'", alias, command)
	}
}

// AddAlias adds an alias to the shell config file
func AddAlias(alias string) error {
	shellType := DetectShell()
	if shellType == ShellUnknown {
		return fmt.Errorf("could not detect shell type")
	}

	configPath, err := ShellConfigPath(shellType)
	if err != nil {
		return err
	}

	// Get the graft executable path
	execPath, err := os.Executable()
	if err != nil {
		execPath = "graft" // Fallback to assuming it's in PATH
	}

	aliasLine := AliasCommand(shellType, alias, execPath)
	marker := "# Added by Git-Graft"

	// Check if alias already exists
	if exists, err := aliasExists(configPath, alias); err == nil && exists {
		return fmt.Errorf("alias '%s' already exists in %s", alias, configPath)
	}

	// Append alias to config file
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	content := fmt.Sprintf("\n%s\n%s\n", marker, aliasLine)
	if _, err := f.WriteString(content); err != nil {
		return err
	}

	return nil
}

// RemoveAlias removes the graft alias from the shell config
func RemoveAlias(alias string) error {
	shellType := DetectShell()
	if shellType == ShellUnknown {
		return fmt.Errorf("could not detect shell type")
	}

	configPath, err := ShellConfigPath(shellType)
	if err != nil {
		return err
	}

	// Read the file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Find and remove the alias lines
	lines := strings.Split(string(content), "\n")
	var newLines []string
	skipNext := false

	for _, line := range lines {
		if strings.Contains(line, "# Added by Git-Graft") {
			skipNext = true
			continue
		}
		if skipNext && strings.Contains(line, fmt.Sprintf("alias %s=", alias)) {
			skipNext = false
			continue
		}
		skipNext = false
		newLines = append(newLines, line)
	}

	// Write back
	return os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// aliasExists checks if an alias already exists in the config file
func aliasExists(configPath, alias string) (bool, error) {
	f, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	aliasPattern := fmt.Sprintf("alias %s=", alias)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, aliasPattern) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

// GetShellName returns a human-readable shell name
func GetShellName(shellType ShellType) string {
	switch shellType {
	case ShellBash:
		return "Bash"
	case ShellZsh:
		return "Zsh"
	case ShellFish:
		return "Fish"
	default:
		return "Unknown"
	}
}

// GetReloadCommand returns the command to reload shell config
func GetReloadCommand(shellType ShellType) string {
	configPath, _ := ShellConfigPath(shellType)

	switch shellType {
	case ShellFish:
		return "source " + configPath
	default:
		return "source " + configPath
	}
}
