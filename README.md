# Git-Graft

A chill TUI for git commands.

```
  ╔═╗╦╔╦╗  ╔═╗╦═╗╔═╗╔═╗╔╦╗
  ║ ╦║ ║───║ ╦╠╦╝╠═╣╠╣  ║
  ╚═╝╩ ╩   ╚═╝╩╚═╩ ╩╚   ╩
```

## Install

### One-liner (requires Go 1.21+)

```bash
go install github.com/anandtyagi/gitgraft/cmd/graft@latest
```

Then add the alias (optional):

```bash
echo "alias gg='graft'" >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc
```

### From source

```bash
git clone https://github.com/anandtyagi/gitgraft.git
cd gitgraft
make install
```

## Usage

```bash
graft    # or 'gg' with alias
```

## Features

- **Command Center** - Quick access with fuzzy search
- **New Branch** - Create from any base branch
- **Switch Branch** - Split view with commit history
- **Commit** - Visual file staging with regex select
- **Smart Errors** - Recovery actions for common issues

## Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate |
| `Enter` | Select |
| `Tab` | Next field |
| `Esc` | Back |
| `Space` | Toggle |
| `Ctrl+A` | Select all |
| `Ctrl+R` | Regex select |
| `/` | Search |
| `q` | Quit |

## Config

`~/.config/graft/config.yaml`

```yaml
default_branch: main
push_after_commit: true
```

## License

MIT
