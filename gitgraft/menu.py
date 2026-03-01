"""Main menu TUI for GitGraft"""

from textual.app import App, ComposeResult
from textual.containers import Container, Vertical
from textual.widgets import Header, Footer, Static, Input, ListView, ListItem, Label
from textual.binding import Binding
from gitgraft.branch_selector import BranchSelectorApp


class MenuItem(ListItem):
    """Custom list item for menu options"""
    
    def __init__(self, title: str, description: str, command: str):
        super().__init__()
        self.title = title
        self.description = description
        self.command = command
    
    def compose(self) -> ComposeResult:
        yield Label(f"[bold]{self.title}[/bold]\n[dim]{self.description}[/dim]", markup=True)


class MainMenuApp(App):
    """Main menu application"""
    
    CSS = """
    Screen {
        align: center middle;
    }
    
    #menu-container {
        width: 70;
        height: auto;
        border: solid $primary;
        padding: 2;
    }
    
    #title {
        text-align: center;
        text-style: bold;
        color: $accent;
        margin-bottom: 2;
    }
    
    #search-container {
        height: auto;
        margin-bottom: 1;
    }
    
    #search-input {
        width: 100%;
    }
    
    #menu-list {
        height: 20;
        border: solid $primary-background;
    }
    
    MenuItem {
        padding: 1;
    }
    
    MenuItem:hover {
        background: $boost;
    }
    
    MenuItem.-highlighted {
        background: $accent;
    }
    
    #help-text {
        text-align: center;
        color: $text-muted;
        margin-top: 1;
    }
    """
    
    BINDINGS = [
        Binding("q", "quit", "Quit", show=True),
        Binding("escape", "quit", "Quit", show=False),
        Binding("enter", "select", "Select", show=True),
        Binding("/", "focus_search", "Search", show=True),
    ]
    
    def __init__(self):
        super().__init__()
        self.menu_items = [
            {
                "title": "📋 Branch Selector",
                "description": "Browse and view commit history for branches",
                "command": "branch",
            },
            {
                "title": "⚙️  Run Onboarding",
                "description": "Configure GitGraft settings",
                "command": "onboard",
            },
        ]
        self.filtered_items = self.menu_items.copy()
    
    def compose(self) -> ComposeResult:
        """Compose the UI"""
        yield Header()
        yield Container(
            Static("🌳 GitGraft - Main Menu", id="title"),
            Container(
                Input(placeholder="Search commands...", id="search-input"),
                id="search-container"
            ),
            ListView(id="menu-list"),
            Static("Use ↑↓ to navigate, Enter to select, / to search, Q to quit", id="help-text"),
            id="menu-container"
        )
        yield Footer()
    
    def on_mount(self) -> None:
        """Initialize after mounting"""
        self.update_menu_list()
    
    def update_menu_list(self):
        """Update the menu list"""
        menu_list = self.query_one("#menu-list", ListView)
        menu_list.clear()
        
        for item in self.filtered_items:
            menu_item = MenuItem(
                item["title"],
                item["description"],
                item["command"]
            )
            menu_list.append(menu_item)
    
    def on_input_changed(self, event: Input.Changed) -> None:
        """Handle search input changes"""
        if event.input.id == "search-input":
            search_term = event.value.lower()
            self.filtered_items = [
                item for item in self.menu_items
                if search_term in item["title"].lower() or search_term in item["description"].lower()
            ]
            self.update_menu_list()
    
    def on_list_view_selected(self, event: ListView.Selected) -> None:
        """Handle menu selection"""
        if isinstance(event.item, MenuItem):
            self.execute_command(event.item.command)
    
    def action_select(self) -> None:
        """Handle Enter key to select current item"""
        menu_list = self.query_one("#menu-list", ListView)
        if menu_list.highlighted_child:
            if isinstance(menu_list.highlighted_child, MenuItem):
                self.execute_command(menu_list.highlighted_child.command)
    
    def execute_command(self, command: str):
        """Execute a menu command"""
        self.exit()
        
        if command == "branch":
            app = BranchSelectorApp()
            app.run()
        elif command == "onboard":
            from gitgraft.onboarding import OnboardingApp
            app = OnboardingApp()
            app.run()
    
    def action_focus_search(self) -> None:
        """Focus the search input"""
        search_input = self.query_one("#search-input", Input)
        search_input.focus()


def run_main_menu():
    """Run the main menu"""
    app = MainMenuApp()
    app.run()


if __name__ == "__main__":
    run_main_menu()
