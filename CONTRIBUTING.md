# Contributing to GitGraft

Thank you for your interest in contributing to GitGraft! This document provides guidelines and instructions for contributing.

## Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/yourusername/gitgraft.git
   cd gitgraft
   ```

2. **Install Dependencies**
   ```bash
   pip install -e .
   ```

3. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Code Style

- Follow PEP 8 style guidelines
- Use meaningful variable and function names
- Add docstrings to functions and classes
- Keep functions focused and concise

## Testing

Before submitting a PR:

1. **Test your changes manually**
   ```bash
   graph branch     # Test branch selector
   graph onboard    # Test onboarding
   graph            # Test main menu
   ```

2. **Verify module imports**
   ```bash
   python3 -c "from gitgraft.branch_selector import BranchSelectorApp"
   python3 -c "from gitgraft.menu import MainMenuApp"
   python3 -c "from gitgraft.onboarding import OnboardingApp"
   ```

3. **Test in different git repositories**
   - Repository with multiple branches
   - Repository with single branch
   - Repository with many commits

## Pull Request Process

1. **Update Documentation**
   - Update README.md if you add new features
   - Add comments for complex logic
   - Update keyboard shortcuts table if needed

2. **Commit Messages**
   - Use clear, descriptive commit messages
   - Start with a verb (Add, Fix, Update, Remove)
   - Example: "Add search functionality to commit viewer"

3. **Submit PR**
   - Provide a clear description of changes
   - Reference any related issues
   - Include screenshots for UI changes

## Adding New Features

### New TUI Screen

1. Create a new file in `gitgraft/` (e.g., `commit_viewer.py`)
2. Inherit from `textual.app.App`
3. Define CSS styling
4. Implement `compose()` method
5. Add event handlers
6. Add to menu in `menu.py`
7. Add command line argument in `cli.py`

Example structure:
```python
from textual.app import App, ComposeResult
from textual.widgets import Header, Footer

class MyNewApp(App):
    CSS = """
    /* Your CSS here */
    """
    
    def compose(self) -> ComposeResult:
        yield Header()
        # Your widgets here
        yield Footer()
```

### New Git Operation

1. Add method to `GitRepo` class in `git_utils.py`
2. Handle errors appropriately
3. Return structured data
4. Add docstring with parameters and return type

## Project Structure

```
gitgraft/
├── gitgraft/
│   ├── __init__.py          # Package init
│   ├── cli.py               # CLI entry point
│   ├── config.py            # Configuration
│   ├── git_utils.py         # Git operations
│   ├── onboarding.py        # Onboarding TUI
│   ├── menu.py              # Main menu TUI
│   └── branch_selector.py   # Branch selector TUI
├── docs/
│   └── screenshots/         # UI screenshots
├── setup.py                 # Package setup
├── requirements.txt         # Dependencies
├── README.md               # Documentation
└── CONTRIBUTING.md         # This file
```

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the code
- Suggestions for improvements

Thank you for contributing! 🌳
