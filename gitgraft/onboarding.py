"""Onboarding flow for GitGraft"""

import os
import subprocess
from textual.app import App, ComposeResult
from textual.containers import Container, Vertical, Horizontal
from textual.widgets import Header, Footer, Static, Button, Label
from textual.binding import Binding
from gitgraft.config import Config


class OnboardingApp(App):
    """Onboarding application"""
    
    CSS = """
    Screen {
        align: center middle;
    }
    
    #onboarding-container {
        width: 80;
        height: auto;
        border: solid $primary;
        padding: 2;
    }
    
    #title {
        text-align: center;
        text-style: bold;
        color: $accent;
        margin-bottom: 1;
    }
    
    #content {
        margin-bottom: 2;
        text-align: center;
    }
    
    #buttons {
        align: center middle;
        width: 100%;
        height: auto;
    }
    
    Button {
        margin: 0 1;
    }
    
    .success {
        color: $success;
        text-align: center;
        margin-top: 1;
    }
    """
    
    BINDINGS = [
        Binding("q", "quit", "Quit", show=True),
    ]
    
    def __init__(self):
        super().__init__()
        self.config = Config()
        self.stage = "welcome"  # welcome, alias, complete
    
    def compose(self) -> ComposeResult:
        """Compose the UI"""
        yield Header()
        yield Container(
            Static("🌳 Welcome to GitGraft!", id="title"),
            Static(
                "GitGraft is a TUI (Terminal User Interface) for Git commands.\n\n"
                "Let's get you set up!",
                id="content"
            ),
            Horizontal(
                Button("Continue", variant="primary", id="continue-btn"),
                Button("Skip", variant="default", id="skip-btn"),
                id="buttons"
            ),
            id="onboarding-container"
        )
        yield Footer()
    
    def on_button_pressed(self, event: Button.Pressed) -> None:
        """Handle button presses"""
        if event.button.id == "continue-btn":
            if self.stage == "welcome":
                self.show_alias_stage()
            elif self.stage == "alias":
                self.create_alias()
                self.show_complete_stage()
        elif event.button.id == "skip-btn":
            if self.stage == "welcome":
                self.show_alias_stage()
            elif self.stage == "alias":
                self.config.set_alias_preference(False)
                self.show_complete_stage()
        elif event.button.id == "finish-btn":
            self.config.set_onboarding_complete()
            self.exit()
    
    def show_alias_stage(self):
        """Show alias configuration stage"""
        self.stage = "alias"
        container = self.query_one("#onboarding-container")
        
        # Update content
        container.query_one("#title").update("⚡ Shell Alias")
        container.query_one("#content").update(
            "Would you like to create a shell alias?\n\n"
            "This will allow you to use 'gg' instead of 'git-graft'\n\n"
            "Example: 'gg' will show the main menu\n"
            "         'gg branch' will show the branch selector"
        )
        
        # Update buttons
        buttons = container.query_one("#buttons")
        buttons.remove_children()
        buttons.mount(Button("Yes, create alias", variant="primary", id="continue-btn"))
        buttons.mount(Button("No, thanks", variant="default", id="skip-btn"))
    
    def create_alias(self):
        """Create shell alias"""
        self.config.set_alias_preference(True)
        
        # Determine shell configuration file
        shell = os.environ.get("SHELL", "")
        config_file = None
        
        if "zsh" in shell:
            config_file = os.path.expanduser("~/.zshrc")
        elif "bash" in shell:
            bashrc = os.path.expanduser("~/.bashrc")
            bash_profile = os.path.expanduser("~/.bash_profile")
            config_file = bashrc if os.path.exists(bashrc) else bash_profile
        
        if config_file:
            alias_line = "\n# GitGraft alias\nalias gg='git-graft'\n"
            
            # Check if alias already exists
            try:
                with open(config_file, 'r') as f:
                    content = f.read()
                    if "alias gg='git-graft'" not in content and 'alias gg="git-graft"' not in content:
                        with open(config_file, 'a') as f:
                            f.write(alias_line)
            except FileNotFoundError:
                with open(config_file, 'w') as f:
                    f.write(alias_line)
    
    def show_complete_stage(self):
        """Show completion stage"""
        self.stage = "complete"
        container = self.query_one("#onboarding-container")
        
        # Update content
        container.query_one("#title").update("✅ All Set!")
        
        if self.config.get_alias_preference():
            content = (
                "GitGraft is now configured!\n\n"
                "Alias 'gg' has been added to your shell config.\n"
                "Restart your terminal or run 'source ~/.bashrc' (or ~/.zshrc)\n\n"
                "Try these commands:\n"
                "  • git-graft (or gg) - Main menu\n"
                "  • git-graft branch (or gg branch) - Branch selector"
            )
        else:
            content = (
                "GitGraft is now configured!\n\n"
                "Try these commands:\n"
                "  • git-graft - Main menu\n"
                "  • git-graft branch - Branch selector"
            )
        
        container.query_one("#content").update(content)
        
        # Update buttons
        buttons = container.query_one("#buttons")
        buttons.remove_children()
        buttons.mount(Button("Finish", variant="success", id="finish-btn"))


def run_onboarding():
    """Run the onboarding flow"""
    app = OnboardingApp()
    app.run()


if __name__ == "__main__":
    run_onboarding()
