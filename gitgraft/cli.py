"""Main CLI entry point for GitGraft"""

import sys
import argparse
from pathlib import Path
from gitgraft.config import Config
from gitgraft.onboarding import OnboardingApp
from gitgraft.menu import MainMenuApp
from gitgraft.branch_selector import BranchSelectorApp


def main():
    """Main entry point for the graph command"""
    config = Config()
    
    parser = argparse.ArgumentParser(description="GitGraft - TUI for Git commands")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")
    
    # Add subcommands
    subparsers.add_parser("branch", help="Browse and select git branches")
    subparsers.add_parser("onboard", help="Run onboarding flow")
    
    args = parser.parse_args()
    
    # Check if this is first run
    if not config.is_configured() and args.command != "onboard":
        print("Welcome to GitGraft! Running onboarding...")
        app = OnboardingApp()
        app.run()
        return
    
    # Handle commands
    if args.command == "branch":
        app = BranchSelectorApp()
        app.run()
    elif args.command == "onboard":
        app = OnboardingApp()
        app.run()
    else:
        # No command specified, show main menu
        app = MainMenuApp()
        app.run()


if __name__ == "__main__":
    main()
