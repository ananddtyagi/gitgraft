"""Branch selector TUI for GitGraft"""

from textual.app import App, ComposeResult
from textual.containers import Container, Vertical, Horizontal, VerticalScroll
from textual.widgets import Header, Footer, Static, Input, ListView, ListItem, Label
from textual.binding import Binding
from textual.reactive import reactive
from gitgraft.git_utils import GitRepo


class BranchListItem(ListItem):
    """Custom list item for branches"""
    
    def __init__(self, branch_name: str, is_current: bool = False):
        super().__init__()
        self.branch_name = branch_name
        self.is_current = is_current
        prefix = "→ " if is_current else "  "
        style = "bold green" if is_current else ""
        self.label = Label(f"{prefix}{branch_name}", markup=True)
    
    def compose(self) -> ComposeResult:
        yield self.label


class BranchSelectorApp(App):
    """Branch selector application with split view"""
    
    CSS = """
    Screen {
        layout: vertical;
    }
    
    #main-container {
        layout: horizontal;
        height: 100%;
    }
    
    #left-panel {
        width: 40%;
        border-right: solid $primary;
    }
    
    #right-panel {
        width: 60%;
        padding: 1 2;
    }
    
    #search-container {
        height: auto;
        padding: 1;
        background: $panel;
    }
    
    #search-input {
        width: 100%;
    }
    
    #branch-list-container {
        height: 1fr;
    }
    
    #branch-list {
        height: 100%;
    }
    
    #commits-title {
        text-style: bold;
        color: $accent;
        margin-bottom: 1;
    }
    
    #commits-container {
        height: 1fr;
    }
    
    .commit-item {
        margin-bottom: 1;
    }
    
    .commit-hash {
        color: $warning;
        text-style: bold;
    }
    
    .commit-author {
        color: $success;
    }
    
    .commit-date {
        color: $text-muted;
    }
    
    .commit-message {
        color: $text;
    }
    
    .no-commits {
        color: $text-muted;
        text-align: center;
        margin-top: 2;
    }
    
    ListItem {
        padding: 0 1;
    }
    
    ListItem:hover {
        background: $boost;
    }
    
    ListItem.-highlighted {
        background: $accent;
    }
    """
    
    BINDINGS = [
        Binding("q", "quit", "Quit", show=True),
        Binding("escape", "quit", "Quit", show=False),
        Binding("/", "focus_search", "Search", show=True),
    ]
    
    selected_branch = reactive("")
    
    def __init__(self):
        super().__init__()
        try:
            self.git_repo = GitRepo()
            self.all_branches = self.git_repo.get_branches()
            self.filtered_branches = self.all_branches.copy()
            self.current_branch = self.git_repo.get_current_branch()
        except ValueError as e:
            self.git_repo = None
            self.all_branches = []
            self.filtered_branches = []
            self.current_branch = None
            self.error_message = str(e)
    
    def compose(self) -> ComposeResult:
        """Compose the UI"""
        yield Header()
        
        if not self.git_repo:
            yield Container(
                Static(f"❌ Error: {self.error_message}", classes="no-commits"),
                id="main-container"
            )
        else:
            yield Horizontal(
                # Left panel - Branch list
                Vertical(
                    Container(
                        Input(placeholder="Search branches...", id="search-input"),
                        id="search-container"
                    ),
                    Container(
                        ListView(id="branch-list"),
                        id="branch-list-container"
                    ),
                    id="left-panel"
                ),
                # Right panel - Commit history
                Vertical(
                    Static("Commit History", id="commits-title"),
                    VerticalScroll(
                        Static("Select a branch to view commits", classes="no-commits"),
                        id="commits-container"
                    ),
                    id="right-panel"
                ),
                id="main-container"
            )
        
        yield Footer()
    
    def on_mount(self) -> None:
        """Initialize after mounting"""
        if self.git_repo:
            self.update_branch_list()
            if self.current_branch:
                self.selected_branch = self.current_branch
                self.update_commits()
    
    def update_branch_list(self):
        """Update the branch list"""
        branch_list = self.query_one("#branch-list", ListView)
        branch_list.clear()
        
        for branch in self.filtered_branches:
            is_current = branch == self.current_branch
            item = BranchListItem(branch, is_current)
            branch_list.append(item)
        
        # Highlight current branch initially
        if self.current_branch and self.current_branch in self.filtered_branches:
            try:
                index = self.filtered_branches.index(self.current_branch)
                branch_list.index = index
            except ValueError:
                pass
    
    def on_input_changed(self, event: Input.Changed) -> None:
        """Handle search input changes"""
        if event.input.id == "search-input":
            search_term = event.value.lower()
            self.filtered_branches = [
                branch for branch in self.all_branches
                if search_term in branch.lower()
            ]
            self.update_branch_list()
    
    def on_list_view_selected(self, event: ListView.Selected) -> None:
        """Handle branch selection"""
        if isinstance(event.item, BranchListItem):
            self.selected_branch = event.item.branch_name
            self.update_commits()
    
    def watch_selected_branch(self, branch: str) -> None:
        """React to selected branch changes"""
        if branch:
            self.update_commits()
    
    def update_commits(self):
        """Update the commit history display"""
        if not self.selected_branch:
            return
        
        commits_container = self.query_one("#commits-container", VerticalScroll)
        commits_container.remove_children()
        
        # Update title
        title = self.query_one("#commits-title", Static)
        title.update(f"Commit History - {self.selected_branch}")
        
        # Get commits
        commits = self.git_repo.get_commits(self.selected_branch)
        
        if not commits:
            commits_container.mount(
                Static("No commits found", classes="no-commits")
            )
            return
        
        # Display commits
        for commit_hash, author, date, message in commits:
            commit_text = (
                f"[bold yellow]{commit_hash}[/bold yellow] "
                f"[dim]{date}[/dim]\n"
                f"[green]{author}[/green]\n"
                f"{message}"
            )
            commits_container.mount(
                Static(commit_text, classes="commit-item", markup=True)
            )
    
    def action_focus_search(self) -> None:
        """Focus the search input"""
        search_input = self.query_one("#search-input", Input)
        search_input.focus()


def run_branch_selector():
    """Run the branch selector"""
    app = BranchSelectorApp()
    app.run()


if __name__ == "__main__":
    run_branch_selector()
