"""Git utilities for GitGraft"""

import git
from pathlib import Path
from typing import List, Optional, Tuple


class GitRepo:
    """Git repository interface"""
    
    def __init__(self, path: Optional[str] = None):
        """Initialize Git repository
        
        Args:
            path: Path to repository. If None, uses current directory
        """
        try:
            if path:
                self.repo = git.Repo(path)
            else:
                self.repo = git.Repo(".", search_parent_directories=True)
        except git.InvalidGitRepository:
            raise ValueError("Not a git repository")
    
    def get_branches(self) -> List[str]:
        """Get all local branches
        
        Returns:
            List of branch names
        """
        return [str(branch.name) for branch in self.repo.branches]
    
    def get_current_branch(self) -> Optional[str]:
        """Get current branch name
        
        Returns:
            Current branch name or None if detached HEAD
        """
        try:
            return str(self.repo.active_branch.name)
        except TypeError:
            return None
    
    def get_commits(self, branch: str, max_count: int = 50) -> List[Tuple[str, str, str, str]]:
        """Get commits for a branch
        
        Args:
            branch: Branch name
            max_count: Maximum number of commits to return
            
        Returns:
            List of tuples (hash, author, date, message)
        """
        commits = []
        try:
            for commit in self.repo.iter_commits(branch, max_count=max_count):
                short_hash = commit.hexsha[:7]
                author = commit.author.name
                date = commit.committed_datetime.strftime("%Y-%m-%d %H:%M")
                message = commit.message.split('\n')[0]  # First line only
                commits.append((short_hash, author, date, message))
        except git.GitCommandError:
            pass
        return commits
    
    def checkout_branch(self, branch: str) -> bool:
        """Checkout a branch
        
        Args:
            branch: Branch name
            
        Returns:
            True if successful, False otherwise
        """
        try:
            self.repo.git.checkout(branch)
            return True
        except git.GitCommandError:
            return False
