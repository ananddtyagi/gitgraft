# GitGraft Quick Start Guide

Get started with GitGraft in 5 minutes!

## Installation

```bash
# Clone the repository
git clone https://github.com/ananddtyagi/gitgraft.git
cd gitgraft

# Install the package
pip install -e .

# Verify installation
graph --help
```

## First Run

1. **Run Onboarding**
   ```bash
   graph onboard
   ```
   
   This will:
   - Welcome you to GitGraft
   - Offer to create a shell alias (`gg` → `graph`)
   - Save your preferences

2. **Explore the Main Menu**
   ```bash
   graph
   # or if you created the alias:
   gg
   ```
   
   Use arrow keys to navigate, Enter to select.

3. **Try the Branch Selector**
   ```bash
   graph branch
   # or:
   gg branch
   ```

## Daily Usage

### Common Commands

```bash
# Show main menu
gg

# View branches
gg branch

# Re-run onboarding
gg onboard
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate up/down |
| `/` | Search/filter |
| `Enter` | Select item |
| `q` | Quit |
| `Esc` | Exit/Back |

## Tips & Tricks

### Branch Selector Tips

1. **Quick Search**: Press `/` and type to filter branches
2. **Navigate**: Use arrow keys to browse branches
3. **View History**: Select a branch to see its commit history
4. **Current Branch**: Look for the `→` marker

### Main Menu Tips

1. **Search Commands**: Press `/` to filter available commands
2. **Quick Navigation**: Use arrow keys and Enter
3. **Help Text**: Look at the bottom of the screen for hints

## What's Next?

- Explore the branch selector
- Try searching with `/`
- Check out the commit history
- Customize your workflow

## Need Help?

- Run `graph --help` for command-line options
- Check the main [README.md](README.md) for detailed documentation
- See [CONTRIBUTING.md](CONTRIBUTING.md) if you want to contribute

Happy Git-ing! 🌳
