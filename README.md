# claude-television

[English 🇺🇸](./README.md) | [한국어 🇰🇷](./README_ko.md)

> 📺 A TUI dashboard for exploring your Claude Code configuration at a glance.

Claude Code settings are scattered across multiple locations — `settings.json`, `CLAUDE.md`, plugins, hooks, MCP servers, and memories. **claude-television** (`ctv`) brings them all together in a single, read-only terminal dashboard using a powerful channel-based system.

## Features

- **Hybrid TUI Dashboard** — Navigate through a split-pane layout with a channel list on the left and detailed tabs on the right.
- **Channel System** — Automatically discovers and syncs your Claude Code projects from `~/.claude/projects/`.
- **Comprehensive Insights** — View local/global settings, `CLAUDE.md` sections, registered hooks, MCP servers, Git status, and Memory files.
- **Fast & Cached** — Uses an mtime-based caching system for instant load times.
- **Keyboard-driven** — Navigate effortlessly with Vim-style keybindings and even `cd` directly into projects.

## Installation

### Quick Install (macOS / Linux)

The easiest way to install the latest pre-compiled binary is using our install script:

```bash
curl -fsSL https://raw.githubusercontent.com/nopecho/claude-television/main/scripts/install.sh | bash
```

Alternatively, you can download the binary directly from the [GitHub Releases](https://github.com/nopecho/claude-television/releases) page.

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
# Discover and register projects from ~/.claude/projects/
ctv init

# Launch the dashboard
ctv
```

### Enable Directory Navigation
To navigate directly to a project directory when pressing `Alt+Enter` in the TUI, add this function to your shell config (`~/.zshrc` or `~/.bashrc`):

```bash
ctv() { local dir; dir="$(command ctv "$@")"; [ -d "$dir" ] && cd "$dir" || command ctv "$@"; }
```

## Usage

```
ctv          # Launch TUI dashboard
ctv init     # Discover and register channels
ctv version  # Show version
```

## Dashboard

```
┌─ claude-television ──────────────────────────────────────┐
│                                                          │
│  ┌─ Channels ──────────┐  ┌─ Settings [CLAUDE.md] ────┐  │
│  │ ● my-awesome-app    │  │ # Project Instructions    │  │
│  │ ● ctv-backend       │  │                           │  │
│  │ ○ legacy-api        │  │ ## Build                  │  │
│  │ ✕ broken-project    │  │ - go build -o app .       │  │
│  └─────────────────────┘  └───────────────────────────┘  │
│                                                          │
│  ↑↓/jk navigate  ←→/Tab switch  Ctrl+d/u scroll  q quit  │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

## Keybindings

| Key            | Action                            |
| -------------- | --------------------------------- |
| `↑` / `k`      | Move up in channel list           |
| `↓` / `j`      | Move down in channel list         |
| `Tab`/`←`/`→`  | Switch detail tab                 |
| `Ctrl+d` / `u` | Scroll detail view down / up      |
| `Alt+Enter`    | Navigate to project dir (cd)      |
| `/`            | Fuzzy search channels             |
| `q` / `Ctrl+C` | Quit                              |

## Configuration

ctv stores its config at `~/.config/ctv/config.json`:

```json
{
  "channels": {
    "auto_sync": true,
    "cache_ttl": "24h",
    "pins": [
      "my-awesome-app"
    ],
    "groups": {
      "Work": ["ctv-backend"]
    }
  }
}
```

- **`auto_sync`**: Automatically discover new projects on launch.
- **`cache_ttl`**: Duration to keep channel data cached.
- **`pins`**: Array of channel names or IDs to pin to the top.
- **`groups`**: Group channels by assigning names or IDs to a group label.

## Contributing

Contributions are welcome! Please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (following [Conventional Commits](https://www.conventionalcommits.org/))
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## License

MIT
