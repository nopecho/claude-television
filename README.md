# claude-television

> 📺 A TUI dashboard for exploring your Claude Code configuration at a glance.

Claude Code settings are scattered across multiple locations — `settings.json`, `CLAUDE.md`, plugins, hooks, and project-specific configs. **claude-television** (`ctv`) brings them all together in a single, read-only terminal dashboard.

## Features

- **Global Settings** — View `settings.json`, `settings.local.json`, and global `CLAUDE.md` in one place
- **Project Explorer** — Scan directories to see which projects have Claude Code configs
- **Skills & Plugins** — Browse installed plugins, their versions, and activation status
- **Hooks Overview** — Inspect registered hooks at a glance
- **Keyboard-driven** — Navigate with vim-style keybindings

## Installation

### From Source

Requires Go 1.26+:

```bash
go install github.com/nopecho/claude-television@latest
```

### Using asdf

```bash
asdf install golang 1.26.1
asdf local golang 1.26.1
go install github.com/nopecho/claude-television@latest
```

## Quick Start

```bash
# Register a directory to scan for projects
ctv scan ~/projects

# Launch the dashboard
ctv
```

## Usage

```
ctv                      # Launch TUI dashboard
ctv scan <path>          # Register a scan path
ctv scan --list          # List registered scan paths
ctv scan --remove <path> # Remove a scan path
ctv version              # Show version
```

## Dashboard

```
┌─ claude-television ──────────────────────────────────────┐
│                                                          │
│  [Global] [Projects] [Skills] [Hooks]                    │
│                                                          │
│  ┌─ List ──────────────┐  ┌─ Detail ──────────────────┐  │
│  │ ● Settings       ✓  │  │ model: opus              │  │
│  │ ● Local Settings  ✓  │  │ language: korean         │  │
│  │ ● CLAUDE.md      ✓  │  │ permissions:             │  │
│  │                      │  │   allow: [Bash, Read...] │  │
│  └──────────────────────┘  └──────────────────────────┘  │
│                                                          │
│  ↑↓/jk navigate  ←→/Tab switch  / filter  q quit        │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

## Keybindings

| Key | Action |
|-----|--------|
| `↑`/`k` | Move up |
| `↓`/`j` | Move down |
| `Tab`/`←`/`→` | Switch tab |
| `Enter` | Select / toggle detail |
| `/` | Filter list (post-MVP) |
| `q` / `Ctrl+C` | Quit |

## Configuration

ctv stores its config at `~/.config/ctv/config.yaml`:

```yaml
scan:
  roots:
    - ~/projects
    - ~/work
  ignore:
    - node_modules
    - .git
    - vendor
```

## Contributing

Contributions are welcome! Please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (following [Conventional Commits](https://www.conventionalcommits.org/))
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## License

MIT
