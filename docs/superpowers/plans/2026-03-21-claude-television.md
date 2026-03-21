# claude-television (ctv) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Claude Code의 로컬 설정 상태를 TUI 대시보드로 한눈에 탐색하는 읽기 전용 CLI 도구를 구현한다.

**Architecture:** Go CLI(Cobra) → 데이터 수집 레이어(internal/claude, internal/scanner) → TUI 렌더링(Bubble Tea). 설정 파일은 Viper로 관리하고, 데이터 수집은 병렬로 처리한 뒤 TUI 모델에 주입한다.

**Tech Stack:** Go 1.26, Cobra, Bubble Tea, Lip Gloss, Bubbles, Viper, asdf

**Spec:** `docs/superpowers/specs/2026-03-21-claude-television-design.md`

---

## File Map

```
claude-television/
├── main.go                          # 엔트리포인트
├── cmd/
│   ├── root.go                      # root 커맨드 (기본 = dashboard)
│   ├── scan.go                      # scan 서브커맨드 (경로 등록/삭제/목록)
│   └── version.go                   # version 서브커맨드
├── internal/
│   ├── config/
│   │   └── config.go                # Viper 기반 ctv 설정 로드/저장
│   ├── claude/
│   │   ├── settings.go              # settings.json 파싱
│   │   ├── settings_test.go
│   │   ├── claudemd.go              # CLAUDE.md 파싱 (섹션 추출)
│   │   ├── claudemd_test.go
│   │   ├── plugins.go               # installed_plugins.json + enabledPlugins 조인
│   │   ├── plugins_test.go
│   │   ├── skills.go                # ~/.claude/skills/ 로컬 스킬 스캔
│   │   ├── skills_test.go
│   │   ├── hooks.go                 # hooks 설정 파싱 (nested event structure)
│   │   ├── hooks_test.go
│   │   ├── projects.go              # ~/.claude/projects/ 스캔 + 경로 디코딩
│   │   └── projects_test.go
│   ├── scanner/
│   │   ├── scanner.go               # 프로젝트 디렉토리 스캔, Claude 설정 탐지
│   │   └── scanner_test.go
│   └── tui/
│       ├── app.go                   # bubbletea 메인 모델 (탭 라우팅)
│       ├── global.go                # Global 탭
│       ├── projects.go              # Projects 탭
│       ├── skills.go                # Skills 탭
│       ├── hooks.go                 # Hooks 탭
│       ├── tabs.go                  # 탭 네비게이션 컴포넌트
│       └── styles.go                # lipgloss 스타일 정의
├── .tool-versions                   # asdf (golang 1.26.1)
├── .gitignore
├── CLAUDE.md
├── README.md
├── LICENSE                          # MIT License
├── CHANGELOG.md                     # Keep a Changelog 형식
├── go.mod
└── go.sum
```

---

### Task 0: 프로젝트 부트스트랩

**Files:**
- Create: `.tool-versions`, `.gitignore`, `go.mod`, `main.go`, `CLAUDE.md`, `README.md`, `LICENSE`, `CHANGELOG.md`

- [ ] **Step 1: asdf 버전 파일 생성**

`.tool-versions`:
```
golang 1.26.1
```

- [ ] **Step 2: Go 모듈 초기화**

Run: `go mod init github.com/nopecho/claude-television`

- [ ] **Step 3: .gitignore 생성**

`.gitignore`:
```gitignore
# Binary
ctv
claude-television

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work
go.work.sum

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Debug
__debug_bin*

# Distribution
dist/

# Internal docs (superpowers specs/plans)
docs/superpowers/
```

- [ ] **Step 4: CLAUDE.md 생성**

`CLAUDE.md`:
```markdown
# claude-television (ctv)

Claude Code의 로컬 설정 상태를 TUI 대시보드로 탐색하는 읽기 전용 CLI 도구.

## 빌드 및 실행

- `go build -o ctv .` — 빌드
- `go run .` — 실행
- `go test ./...` — 전체 테스트
- `go test ./internal/claude/... -v` — claude 패키지 테스트
- `go test ./internal/scanner/... -v` — scanner 패키지 테스트

## 프로젝트 구조

- `cmd/` — Cobra CLI 커맨드
- `internal/claude/` — Claude Code 설정 파일 파싱 (settings, CLAUDE.md, plugins, skills, hooks, projects)
- `internal/scanner/` — 프로젝트 디렉토리 스캔
- `internal/config/` — Viper 기반 ctv 자체 설정 관리
- `internal/tui/` — Bubble Tea TUI 컴포넌트

## 코딩 컨벤션

- Go 표준 프로젝트 레이아웃 (`cmd/`, `internal/`)
- 에러 처리: `fmt.Errorf("context: %w", err)` 패턴 사용
- 테스트: `testdata/` 디렉토리에 픽스처 파일 배치, 테이블 드리븐 테스트 선호
- 외부 의존성 최소화, 필요한 경우만 추가
```

- [ ] **Step 5: README.md 생성**

`README.md`:
```markdown
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

\```bash
go install github.com/nopecho/claude-television@latest
\```

### Using asdf

\```bash
asdf install golang 1.26.1
asdf local golang 1.26.1
go install github.com/nopecho/claude-television@latest
\```

## Quick Start

\```bash
# Register a directory to scan for projects
ctv scan ~/projects

# Launch the dashboard
ctv
\```

## Usage

\```
ctv                      # Launch TUI dashboard
ctv scan <path>          # Register a scan path
ctv scan --list          # List registered scan paths
ctv scan --remove <path> # Remove a scan path
ctv version              # Show version
\```

## Dashboard

\```
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
\```

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

\```yaml
scan:
  roots:
    - ~/projects
    - ~/work
  ignore:
    - node_modules
    - .git
    - vendor
\```

## Contributing

Contributions are welcome! Please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (following [Conventional Commits](https://www.conventionalcommits.org/))
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## License

MIT
\```

- [ ] **Step 6: LICENSE 생성 (MIT)**

`LICENSE`:
```
MIT License

Copyright (c) 2026 nopecho

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 7: CHANGELOG.md 생성**

`CHANGELOG.md`:
```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- Initial project setup
- TUI dashboard with 4 tabs: Global, Projects, Skills, Hooks
- `ctv scan` command for managing project scan paths
- Read-only view of Claude Code settings, CLAUDE.md, plugins, and hooks
```

- [ ] **Step 8: main.go 스켈레톤 생성**

`main.go`:
```go
package main

import "github.com/nopecho/claude-television/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 9: 의존성 설치**

Run:
```bash
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
go mod tidy
```

- [ ] **Step 10: 커밋**

```bash
git add -A
git commit -m "chore: 프로젝트 초기 셋업

- Go 모듈 초기화, 의존성 설치
- .tool-versions, .gitignore, CLAUDE.md, README.md, LICENSE, CHANGELOG.md 작성
- main.go 엔트리포인트 생성"
```

---

### Task 1: Cobra CLI 프레임워크 (`cmd/`)

**Files:**
- Create: `cmd/root.go`, `cmd/version.go`, `cmd/scan.go`
- Create: `internal/config/config.go`

- [ ] **Step 1: internal/config/config.go 구현**

Viper를 사용해 `~/.config/ctv/config.yaml`을 로드/저장한다.

```go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Scan ScanConfig `mapstructure:"scan"`
}

type ScanConfig struct {
	Roots  []string `mapstructure:"roots"`
	Ignore []string `mapstructure:"ignore"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ctv")
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir())

	viper.SetDefault("scan.roots", []string{})
	viper.SetDefault("scan.ignore", []string{"node_modules", ".git", "vendor"})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("config read: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}
	return &cfg, nil
}

func Save() error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return viper.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}

func AddScanRoot(path string) error {
	roots := viper.GetStringSlice("scan.roots")
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	for _, r := range roots {
		if r == abs {
			return fmt.Errorf("path already registered: %s", abs)
		}
	}
	roots = append(roots, abs)
	viper.Set("scan.roots", roots)
	return Save()
}

func RemoveScanRoot(path string) error {
	roots := viper.GetStringSlice("scan.roots")
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	filtered := make([]string, 0, len(roots))
	found := false
	for _, r := range roots {
		if r == abs {
			found = true
			continue
		}
		filtered = append(filtered, r)
	}
	if !found {
		return fmt.Errorf("path not found: %s", abs)
	}
	viper.Set("scan.roots", filtered)
	return Save()
}
```

- [ ] **Step 2: cmd/root.go 구현**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctv",
	Short: "Claude Code TUI dashboard",
	Long:  "claude-television — A TUI dashboard for exploring your Claude Code configuration at a glance.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Task 5에서 TUI 실행 로직 추가
		fmt.Println("claude-television dashboard (coming soon)")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 3: cmd/version.go 구현**

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print ctv version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ctv %s (%s)\n", version, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

- [ ] **Step 4: cmd/scan.go 구현**

```go
package cmd

import (
	"fmt"

	"github.com/nopecho/claude-television/internal/config"
	"github.com/spf13/cobra"
)

var (
	scanList   bool
	scanRemove string
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Manage project scan paths",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if scanList {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if len(cfg.Scan.Roots) == 0 {
				fmt.Println("No scan paths registered. Use: ctv scan <path>")
				return nil
			}
			for _, r := range cfg.Scan.Roots {
				fmt.Println(r)
			}
			return nil
		}

		if scanRemove != "" {
			if _, err := config.Load(); err != nil {
				return err
			}
			if err := config.RemoveScanRoot(scanRemove); err != nil {
				return err
			}
			fmt.Printf("Removed: %s\n", scanRemove)
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("path required. Usage: ctv scan <path>")
		}

		// Load config first to initialize viper
		if _, err := config.Load(); err != nil {
			return err
		}
		if err := config.AddScanRoot(args[0]); err != nil {
			return err
		}
		fmt.Printf("Added: %s\n", args[0])
		return nil
	},
}

func init() {
	scanCmd.Flags().BoolVar(&scanList, "list", false, "List registered scan paths")
	scanCmd.Flags().StringVar(&scanRemove, "remove", "", "Remove a scan path")
	rootCmd.AddCommand(scanCmd)
}
```

- [ ] **Step 5: 빌드 및 동작 확인**

Run:
```bash
go build -o ctv .
./ctv version
./ctv scan --list
```
Expected: 버전 출력, 빈 스캔 목록 표시

- [ ] **Step 6: 커밋**

```bash
git add cmd/ internal/config/ main.go
git commit -m "feat: Cobra CLI 프레임워크 구현

- root, version, scan 커맨드 추가
- Viper 기반 설정 관리 (scan roots 등록/삭제/목록)"
```

---

### Task 2: Claude Code 설정 파싱 — settings, CLAUDE.md (`internal/claude/`)

**Files:**
- Create: `internal/claude/settings.go`, `internal/claude/settings_test.go`
- Create: `internal/claude/claudemd.go`, `internal/claude/claudemd_test.go`
- Create: `internal/claude/testdata/settings.json`, `internal/claude/testdata/CLAUDE.md`

- [ ] **Step 1: 테스트 픽스처 생성**

`internal/claude/testdata/settings.json`:
```json
{
  "model": "opus",
  "language": "korean",
  "env": {
    "CLAUDE_CODE_SHELL": "zsh"
  },
  "permissions": {
    "allow": ["Bash(go:*)", "Read"],
    "deny": ["Bash(rm -rf *)"]
  },
  "enabledPlugins": {
    "superpowers@claude-plugins-official": true,
    "obsidian@obsidian-skills": false
  }
}
```

`internal/claude/testdata/CLAUDE.md`:
```markdown
# Project Instructions

## Build
- go build -o ctv .

## Testing
- go test ./...

## Conventions
- Use standard Go layout
```

- [ ] **Step 2: settings_test.go 작성 (실패하는 테스트)**

```go
package claude_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func testdataDir() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "testdata")
}

func TestParseSettings(t *testing.T) {
	path := filepath.Join(testdataDir(), "settings.json")
	s, err := claude.ParseSettings(path)
	if err != nil {
		t.Fatalf("ParseSettings: %v", err)
	}

	if s.Model != "opus" {
		t.Errorf("Model = %q, want %q", s.Model, "opus")
	}
	if s.Language != "korean" {
		t.Errorf("Language = %q, want %q", s.Language, "korean")
	}
	if len(s.Permissions.Allow) != 2 {
		t.Errorf("Permissions.Allow len = %d, want 2", len(s.Permissions.Allow))
	}
	if len(s.EnabledPlugins) != 2 {
		t.Errorf("EnabledPlugins len = %d, want 2", len(s.EnabledPlugins))
	}
}

func TestParseSettings_NotFound(t *testing.T) {
	_, err := claude.ParseSettings("/nonexistent/settings.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./internal/claude/... -v`
Expected: FAIL (claude 패키지 없음)

- [ ] **Step 4: settings.go 구현**

```go
package claude

import (
	"encoding/json"
	"fmt"
	"os"
)

type Settings struct {
	Model          string                     `json:"model"`
	Language       string                     `json:"language"`
	Env            map[string]string          `json:"env"`
	Permissions    Permissions                `json:"permissions"`
	EnabledPlugins map[string]bool            `json:"enabledPlugins"`
	Hooks          map[string][]HookRule      `json:"hooks"`
	StatusLine     json.RawMessage            `json:"statusLine"`
	PlansDirectory string                     `json:"plansDirectory"`
	TeammateMode   string                     `json:"teammateMode"`
	Raw            map[string]any             `json:"-"`
}

type Permissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

// HookRule는 특정 이벤트에 대한 hook 규칙 (matcher + actions)
type HookRule struct {
	Matcher string       `json:"matcher"`
	Hooks   []HookAction `json:"hooks"`
}

// HookAction은 실행할 개별 hook (command 또는 prompt)
type HookAction struct {
	Type    string `json:"type"`    // "command" or "prompt"
	Command string `json:"command"` // type=command일 때
	Prompt  string `json:"prompt"`  // type=prompt일 때
	Async   bool   `json:"async"`
	Timeout int    `json:"timeout"`
}

func ParseSettings(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read settings: %w", err)
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse settings: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse raw settings: %w", err)
	}
	s.Raw = raw

	return &s, nil
}
```

- [ ] **Step 5: settings 테스트 통과 확인**

Run: `go test ./internal/claude/... -v -run TestParseSettings`
Expected: PASS

- [ ] **Step 6: claudemd_test.go 작성 (실패하는 테스트)**

```go
package claude_test

import (
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestParseClaudeMD(t *testing.T) {
	path := filepath.Join(testdataDir(), "CLAUDE.md")
	md, err := claude.ParseClaudeMD(path)
	if err != nil {
		t.Fatalf("ParseClaudeMD: %v", err)
	}

	if md.LineCount == 0 {
		t.Error("LineCount should not be 0")
	}
	if len(md.Sections) != 3 {
		t.Errorf("Sections len = %d, want 3", len(md.Sections))
	}
	if md.Sections[0] != "Build" {
		t.Errorf("Sections[0] = %q, want %q", md.Sections[0], "Build")
	}
}

func TestParseClaudeMD_NotFound(t *testing.T) {
	_, err := claude.ParseClaudeMD("/nonexistent/CLAUDE.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
```

- [ ] **Step 7: 테스트 실패 확인**

Run: `go test ./internal/claude/... -v -run TestParseClaudeMD`
Expected: FAIL

- [ ] **Step 8: claudemd.go 구현**

```go
package claude

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ClaudeMD struct {
	Path      string
	LineCount int
	Sections  []string
	Content   string
}

func ParseClaudeMD(path string) (*ClaudeMD, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open CLAUDE.md: %w", err)
	}
	defer f.Close()

	var (
		lines    int
		sections []string
		content  strings.Builder
	)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		content.WriteString(line)
		content.WriteString("\n")

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			sections = append(sections, strings.TrimPrefix(trimmed, "## "))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan CLAUDE.md: %w", err)
	}

	return &ClaudeMD{
		Path:      path,
		LineCount: lines,
		Sections:  sections,
		Content:   content.String(),
	}, nil
}
```

- [ ] **Step 9: 전체 테스트 통과 확인**

Run: `go test ./internal/claude/... -v`
Expected: 모든 테스트 PASS

- [ ] **Step 10: 커밋**

```bash
git add internal/claude/
git commit -m "feat: settings.json 및 CLAUDE.md 파싱 구현

- Settings 구조체로 settings.json 파싱 (model, permissions, enabledPlugins 등)
- ClaudeMD 구조체로 CLAUDE.md 파싱 (라인 수, 섹션 목록)
- testdata 픽스처 기반 테스트"
```

---

### Task 3: Claude Code 설정 파싱 — plugins, skills, hooks, projects

**Files:**
- Create: `internal/claude/plugins.go`, `internal/claude/plugins_test.go`
- Create: `internal/claude/skills.go`, `internal/claude/skills_test.go`
- Create: `internal/claude/hooks.go`, `internal/claude/hooks_test.go`
- Create: `internal/claude/projects.go`, `internal/claude/projects_test.go`
- Create: 추가 testdata 픽스처

- [ ] **Step 1: plugins 테스트 픽스처 생성**

`internal/claude/testdata/installed_plugins.json`:
```json
{
  "superpowers@claude-plugins-official": {
    "version": "5.0.5",
    "scope": "global",
    "installPath": "/Users/test/.claude/plugins/cache/claude-plugins-official/superpowers/5.0.5",
    "installedAt": "2026-03-01T00:00:00Z",
    "gitCommitSha": "abc123"
  }
}
```

- [ ] **Step 2: plugins_test.go 작성 (실패하는 테스트)**

```go
package claude_test

import (
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestParseInstalledPlugins(t *testing.T) {
	path := filepath.Join(testdataDir(), "installed_plugins.json")
	plugins, err := claude.ParseInstalledPlugins(path)
	if err != nil {
		t.Fatalf("ParseInstalledPlugins: %v", err)
	}
	if len(plugins) != 1 {
		t.Fatalf("len = %d, want 1", len(plugins))
	}
	p := plugins["superpowers@claude-plugins-official"]
	if p.Version != "5.0.5" {
		t.Errorf("Version = %q, want %q", p.Version, "5.0.5")
	}
}

func TestMergePluginData(t *testing.T) {
	installed := map[string]claude.InstalledPlugin{
		"superpowers@claude-plugins-official": {Version: "5.0.5"},
	}
	enabled := map[string]bool{
		"superpowers@claude-plugins-official": true,
		"obsidian@obsidian-skills":            false,
	}

	merged := claude.MergePluginData(installed, enabled)
	if len(merged) != 2 {
		t.Fatalf("len = %d, want 2", len(merged))
	}

	sp := findPlugin(merged, "superpowers@claude-plugins-official")
	if sp == nil || !sp.Enabled || sp.Version != "5.0.5" {
		t.Errorf("superpowers: got %+v", sp)
	}

	ob := findPlugin(merged, "obsidian@obsidian-skills")
	if ob == nil || ob.Enabled {
		t.Errorf("obsidian: got %+v", ob)
	}
}

func findPlugin(plugins []claude.Plugin, key string) *claude.Plugin {
	for i := range plugins {
		if plugins[i].Key == key {
			return &plugins[i]
		}
	}
	return nil
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./internal/claude/... -v -run TestParseInstalledPlugins`
Expected: FAIL

- [ ] **Step 4: plugins.go 구현**

```go
package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type InstalledPlugin struct {
	Version      string `json:"version"`
	Scope        string `json:"scope"`
	InstallPath  string `json:"installPath"`
	InstalledAt  string `json:"installedAt"`
	GitCommitSha string `json:"gitCommitSha"`
}

type Plugin struct {
	Key         string
	Name        string
	Marketplace string
	Version     string
	Enabled     bool
	Installed   bool
	InstallPath string
}

func ParseInstalledPlugins(path string) (map[string]InstalledPlugin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read installed_plugins: %w", err)
	}
	var plugins map[string]InstalledPlugin
	if err := json.Unmarshal(data, &plugins); err != nil {
		return nil, fmt.Errorf("parse installed_plugins: %w", err)
	}
	return plugins, nil
}

func MergePluginData(installed map[string]InstalledPlugin, enabled map[string]bool) []Plugin {
	seen := make(map[string]bool)
	var result []Plugin

	for key, ip := range installed {
		name, marketplace := splitPluginKey(key)
		result = append(result, Plugin{
			Key:         key,
			Name:        name,
			Marketplace: marketplace,
			Version:     ip.Version,
			Enabled:     enabled[key],
			Installed:   true,
			InstallPath: ip.InstallPath,
		})
		seen[key] = true
	}

	for key, en := range enabled {
		if seen[key] {
			continue
		}
		name, marketplace := splitPluginKey(key)
		result = append(result, Plugin{
			Key:         key,
			Name:        name,
			Marketplace: marketplace,
			Enabled:     en,
			Installed:   false,
		})
	}

	return result
}

func splitPluginKey(key string) (name, marketplace string) {
	parts := strings.SplitN(key, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return key, ""
}
```

- [ ] **Step 5: plugins 테스트 통과 확인**

Run: `go test ./internal/claude/... -v -run TestParse\|TestMerge`
Expected: PASS

- [ ] **Step 6: skills_test.go 작성 + skills.go 구현**

테스트: `internal/claude/testdata/skills/` 디렉토리에 `my-skill/` 서브디렉토리 생성 후 스캔 확인.

```go
// skills.go
package claude

import (
	"fmt"
	"os"
	"path/filepath"
)

type Skill struct {
	Name string
	Path string
}

func ScanLocalSkills(skillsDir string) ([]Skill, error) {
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan skills: %w", err)
	}

	var skills []Skill
	for _, e := range entries {
		if e.IsDir() {
			skills = append(skills, Skill{
				Name: e.Name(),
				Path: filepath.Join(skillsDir, e.Name()),
			})
		}
	}
	return skills, nil
}
```

```go
// skills_test.go
package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestScanLocalSkills(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "my-skill"), 0755)
	os.MkdirAll(filepath.Join(dir, "another-skill"), 0755)

	skills, err := claude.ScanLocalSkills(dir)
	if err != nil {
		t.Fatalf("ScanLocalSkills: %v", err)
	}
	if len(skills) != 2 {
		t.Errorf("len = %d, want 2", len(skills))
	}
}

func TestScanLocalSkills_NotExist(t *testing.T) {
	skills, err := claude.ScanLocalSkills("/nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("len = %d, want 0", len(skills))
	}
}
```

- [ ] **Step 7: skills 테스트 통과 확인**

Run: `go test ./internal/claude/... -v -run TestScanLocalSkills`
Expected: PASS

- [ ] **Step 8: hooks_test.go 작성 + hooks.go 구현**

hooks는 settings.json의 `map[string][]HookRule` 구조를 플랫한 `[]HookDetail`로 변환한다.

```go
// hooks.go
package claude

type HookDetail struct {
	Event   string
	Matcher string
	Type    string // "command" or "prompt"
	Command string
	Async   bool
	Timeout int
	Source  string // "global" or project path
}

func ExtractHooks(settings *Settings, source string) []HookDetail {
	if settings.Hooks == nil {
		return nil
	}

	var result []HookDetail
	for event, rules := range settings.Hooks {
		for _, rule := range rules {
			for _, action := range rule.Hooks {
				result = append(result, HookDetail{
					Event:   event,
					Matcher: rule.Matcher,
					Type:    action.Type,
					Command: action.Command,
					Async:   action.Async,
					Timeout: action.Timeout,
					Source:  source,
				})
			}
		}
	}
	return result
}
```

```go
// hooks_test.go
package claude_test

import (
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestExtractHooks(t *testing.T) {
	s := &claude.Settings{
		Hooks: map[string][]claude.HookRule{
			"PreToolUse": {
				{
					Matcher: "Bash",
					Hooks: []claude.HookAction{
						{Type: "command", Command: "echo test", Timeout: 10},
					},
				},
			},
			"Stop": {
				{
					Hooks: []claude.HookAction{
						{Type: "command", Command: "echo stop"},
					},
				},
			},
		},
	}
	hooks := claude.ExtractHooks(s, "global")
	if len(hooks) != 2 {
		t.Fatalf("len = %d, want 2", len(hooks))
	}

	found := false
	for _, h := range hooks {
		if h.Event == "PreToolUse" && h.Matcher == "Bash" && h.Command == "echo test" && h.Source == "global" {
			found = true
		}
	}
	if !found {
		t.Error("expected PreToolUse hook with Bash matcher")
	}
}

func TestExtractHooks_Empty(t *testing.T) {
	s := &claude.Settings{}
	hooks := claude.ExtractHooks(s, "global")
	if len(hooks) != 0 {
		t.Errorf("len = %d, want 0", len(hooks))
	}
}
```

- [ ] **Step 9: hooks 테스트 통과 확인**

Run: `go test ./internal/claude/... -v -run TestExtractHooks`
Expected: PASS

- [ ] **Step 10: projects_test.go 작성 + projects.go 구현**

```go
// projects.go
package claude

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectMeta struct {
	EncodedName string
	DecodedPath string
	HasMemory   bool
	HasSessions bool
}

func ScanProjectsMeta(projectsDir string) ([]ProjectMeta, error) {
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan projects: %w", err)
	}

	var result []ProjectMeta
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		decoded := decodeProjectPath(name)
		dirPath := filepath.Join(projectsDir, name)

		meta := ProjectMeta{
			EncodedName: name,
			DecodedPath: decoded,
		}

		if _, err := os.Stat(filepath.Join(dirPath, "memory")); err == nil {
			meta.HasMemory = true
		}

		// sessions = any UUID-like subdirectory
		subEntries, _ := os.ReadDir(dirPath)
		for _, se := range subEntries {
			if se.IsDir() && se.Name() != "memory" {
				meta.HasSessions = true
				break
			}
		}

		result = append(result, meta)
	}
	return result, nil
}

func decodeProjectPath(encoded string) string {
	if !strings.HasPrefix(encoded, "-") {
		return encoded
	}
	// Remove leading dash, replace dashes with slashes
	path := "/" + strings.TrimPrefix(encoded, "-")
	path = strings.ReplaceAll(path, "-", "/")

	// Try to find actual path by progressively joining with dashes
	return resolveAmbiguousPath(path, encoded)
}

func resolveAmbiguousPath(candidate string, encoded string) string {
	// Try the simple replacement first
	simple := "/" + strings.ReplaceAll(strings.TrimPrefix(encoded, "-"), "-", "/")
	if pathExists(simple) {
		return simple
	}

	// Fall back: try to reconstruct by checking filesystem
	parts := strings.Split(strings.TrimPrefix(encoded, "-"), "-")
	return bestEffortDecode(parts)
}

func bestEffortDecode(parts []string) string {
	if len(parts) == 0 {
		return "/"
	}

	current := "/"
	i := 0
	for i < len(parts) {
		// Try joining progressively more parts with dashes
		found := false
		for j := len(parts); j > i; j-- {
			candidate := current + strings.Join(parts[i:j], "-")
			if pathExists(candidate) {
				current = candidate + "/"
				i = j
				found = true
				break
			}
		}
		if !found {
			current += parts[i] + "/"
			i++
		}
	}
	return strings.TrimSuffix(current, "/")
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
```

```go
// projects_test.go
package claude_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestScanProjectsMeta(t *testing.T) {
	dir := t.TempDir()

	// project with memory
	p1 := filepath.Join(dir, "-Users-test-projects-app")
	os.MkdirAll(filepath.Join(p1, "memory"), 0755)

	// project with session
	p2 := filepath.Join(dir, "-Users-test-projects-api")
	os.MkdirAll(filepath.Join(p2, "abc-123-session"), 0755)

	metas, err := claude.ScanProjectsMeta(dir)
	if err != nil {
		t.Fatalf("ScanProjectsMeta: %v", err)
	}
	if len(metas) != 2 {
		t.Fatalf("len = %d, want 2", len(metas))
	}

	for _, m := range metas {
		if m.EncodedName == "-Users-test-projects-app" && !m.HasMemory {
			t.Error("expected HasMemory=true for app")
		}
		if m.EncodedName == "-Users-test-projects-api" && !m.HasSessions {
			t.Error("expected HasSessions=true for api")
		}
	}
}

func TestScanProjectsMeta_NotExist(t *testing.T) {
	metas, err := claude.ScanProjectsMeta("/nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("len = %d, want 0", len(metas))
	}
}
```

- [ ] **Step 11: projects 테스트 통과 확인**

Run: `go test ./internal/claude/... -v -run TestScanProjectsMeta`
Expected: PASS

- [ ] **Step 12: 전체 claude 패키지 테스트 통과 확인**

Run: `go test ./internal/claude/... -v`
Expected: 모든 테스트 PASS

- [ ] **Step 13: 커밋**

```bash
git add internal/claude/
git commit -m "feat: plugins, skills, hooks, projects 파싱 구현

- installed_plugins.json + enabledPlugins 조인 로직
- 로컬 스킬 디렉토리 스캔
- hooks 추출 헬퍼
- projects 메타데이터 스캔 + 경로 디코딩"
```

---

### Task 4: 프로젝트 디렉토리 스캐너 (`internal/scanner/`)

**Files:**
- Create: `internal/scanner/scanner.go`, `internal/scanner/scanner_test.go`

- [ ] **Step 1: scanner_test.go 작성 (실패하는 테스트)**

```go
package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/scanner"
)

func TestScanProjects(t *testing.T) {
	root := t.TempDir()

	// project with .claude/ dir
	p1 := filepath.Join(root, "has-claude")
	os.MkdirAll(filepath.Join(p1, ".claude"), 0755)
	os.WriteFile(filepath.Join(p1, ".claude", "settings.json"), []byte("{}"), 0644)

	// project with CLAUDE.md only
	p2 := filepath.Join(root, "has-claudemd")
	os.MkdirAll(p2, 0755)
	os.WriteFile(filepath.Join(p2, "CLAUDE.md"), []byte("# Test"), 0644)

	// project with no claude config
	p3 := filepath.Join(root, "no-claude")
	os.MkdirAll(p3, 0755)

	// ignored dir
	os.MkdirAll(filepath.Join(root, "node_modules"), 0755)

	projects, err := scanner.ScanProjects([]string{root}, []string{"node_modules"})
	if err != nil {
		t.Fatalf("ScanProjects: %v", err)
	}
	if len(projects) != 3 {
		t.Fatalf("len = %d, want 3", len(projects))
	}

	for _, p := range projects {
		switch p.Name {
		case "has-claude":
			if !p.HasClaudeDir {
				t.Error("has-claude: expected HasClaudeDir=true")
			}
		case "has-claudemd":
			if !p.HasClaudeMD {
				t.Error("has-claudemd: expected HasClaudeMD=true")
			}
		case "no-claude":
			if p.HasClaudeDir || p.HasClaudeMD {
				t.Error("no-claude: expected no claude config")
			}
		}
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./internal/scanner/... -v`
Expected: FAIL

- [ ] **Step 3: scanner.go 구현**

```go
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	Name         string
	Path         string
	HasClaudeDir bool // .claude/ directory exists
	HasClaudeMD  bool // CLAUDE.md exists
	HasSettings  bool // .claude/settings.json exists
}

func ScanProjects(roots []string, ignore []string) ([]Project, error) {
	ignoreSet := make(map[string]bool, len(ignore))
	for _, ig := range ignore {
		ignoreSet[ig] = true
	}

	var projects []Project
	for _, root := range roots {
		expanded := expandHome(root)
		entries, err := os.ReadDir(expanded)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("scan root %s: %w", root, err)
		}

		for _, e := range entries {
			if !e.IsDir() || ignoreSet[e.Name()] || e.Name()[0] == '.' {
				continue
			}

			dirPath := filepath.Join(expanded, e.Name())
			p := Project{
				Name: e.Name(),
				Path: dirPath,
			}

			if _, err := os.Stat(filepath.Join(dirPath, ".claude")); err == nil {
				p.HasClaudeDir = true
			}
			if _, err := os.Stat(filepath.Join(dirPath, "CLAUDE.md")); err == nil {
				p.HasClaudeMD = true
			}
			if _, err := os.Stat(filepath.Join(dirPath, ".claude", "settings.json")); err == nil {
				p.HasSettings = true
			}

			projects = append(projects, p)
		}
	}
	return projects, nil
}

func expandHome(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./internal/scanner/... -v`
Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/scanner/
git commit -m "feat: 프로젝트 디렉토리 스캐너 구현

- 등록된 root 경로에서 1단계 깊이로 프로젝트 탐색
- Claude 설정 유무 탐지 (.claude/, CLAUDE.md, settings.json)
- ignore 목록 및 hidden 디렉토리 필터링"
```

---

### Task 5: TUI 대시보드 — 기본 프레임 (`internal/tui/`)

**Files:**
- Create: `internal/tui/styles.go`, `internal/tui/tabs.go`, `internal/tui/app.go`
- Modify: `cmd/root.go`

- [ ] **Step 1: styles.go 구현**

```go
package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Background(lipgloss.Color("236")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244")).
				Padding(0, 2)

	listItemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Padding(0, 1)

	detailStyle = lipgloss.NewStyle().
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))

	statusIcon    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
	statusIconOff = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗")
)
```

- [ ] **Step 2: tabs.go 구현**

```go
package tui

type TabID int

const (
	TabGlobal TabID = iota
	TabProjects
	TabSkills
	TabHooks
)

var tabNames = []string{"Global", "Projects", "Skills", "Hooks"}

func renderTabs(active TabID) string {
	var tabs string
	for i, name := range tabNames {
		if TabID(i) == active {
			tabs += activeTabStyle.Render("[" + name + "]")
		} else {
			tabs += inactiveTabStyle.Render(" " + name + " ")
		}
	}
	return tabs
}
```

- [ ] **Step 3: app.go 구현 — 메인 TUI 모델**

데이터 수집 → 모델 초기화 → 렌더링의 핵심 구조를 구현한다. 각 탭의 상세 렌더링은 Task 6-9에서 구현하므로, 여기서는 탭 전환과 기본 레이아웃만 구현한다.

```go
package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nopecho/claude-television/internal/claude"
	"github.com/nopecho/claude-television/internal/scanner"
)

type DashboardData struct {
	Settings      *claude.Settings
	LocalSettings *claude.Settings
	ClaudeMD      *claude.ClaudeMD
	Plugins       []claude.Plugin
	LocalSkills   []claude.Skill
	Hooks         []claude.HookDetail
	ProjectsMeta  []claude.ProjectMeta
	Projects      []scanner.Project
}

type model struct {
	data       DashboardData
	activeTab  TabID
	cursor     int
	width      int
	height     int
	skillCache []skillItem // Skills 탭 캐시 (초기화 시 계산)
}

func NewModel(data DashboardData) model {
	m := model{
		data:      data,
		activeTab: TabGlobal,
	}
	m.skillCache = m.buildSkillItems() // 초기화 시 캐시
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "right", "l":
			m.activeTab = (m.activeTab + 1) % TabID(len(tabNames))
			m.cursor = 0
		case "shift+tab", "left", "h":
			m.activeTab = (m.activeTab - 1 + TabID(len(tabNames))) % TabID(len(tabNames))
			m.cursor = 0
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			max := m.listLen() - 1
			if m.cursor < max {
				m.cursor++
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	header := titleStyle.Render("📺 claude-television") + "\n\n"
	tabs := renderTabs(m.activeTab) + "\n\n"

	listWidth := m.width*35/100 - 4
	detailWidth := m.width*65/100 - 4
	contentHeight := m.height - 8

	list := borderStyle.Width(listWidth).Height(contentHeight).Render(m.renderList())
	detail := borderStyle.Width(detailWidth).Height(contentHeight).Render(m.renderDetail())

	content := lipgloss.JoinHorizontal(lipgloss.Top, list, detail)
	help := helpStyle.Render("\n  ↑↓/jk navigate  ←→/Tab switch  q quit")

	return header + tabs + content + help
}

func (m model) listLen() int {
	switch m.activeTab {
	case TabGlobal:
		return 3 // Settings, Local Settings, CLAUDE.md
	case TabProjects:
		return len(m.data.Projects)
	case TabSkills:
		return len(m.skillCache)
	case TabHooks:
		return len(m.data.Hooks)
	}
	return 0
}

func (m model) renderList() string {
	switch m.activeTab {
	case TabGlobal:
		return m.renderGlobalList()
	case TabProjects:
		return m.renderProjectsList()
	case TabSkills:
		return m.renderSkillsList()
	case TabHooks:
		return m.renderHooksList()
	}
	return ""
}

func (m model) renderDetail() string {
	switch m.activeTab {
	case TabGlobal:
		return m.renderGlobalDetail()
	case TabProjects:
		return m.renderProjectsDetail()
	case TabSkills:
		return m.renderSkillsDetail()
	case TabHooks:
		return m.renderHooksDetail()
	}
	return ""
}

// 각 탭 렌더링은 별도 파일에서 정의 (global.go, projects.go, skills.go, hooks.go)
// Task 5에서는 stub 파일을 생성하고, Task 6-9에서 내용을 교체한다.

func Run(data DashboardData) error {
	p := tea.NewProgram(NewModel(data), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
```

- [ ] **Step 4: 탭 stub 파일 생성**

각 탭별 파일에 placeholder 렌더링 메서드를 정의한다. 같은 package이므로 app.go의 model에서 호출 가능.

`internal/tui/global.go`:
```go
package tui

func (m model) renderGlobalList() string   { return "Settings\nLocal Settings\nCLAUDE.md" }
func (m model) renderGlobalDetail() string { return "Select an item" }
```

`internal/tui/projects.go`:
```go
package tui

import "fmt"

func (m model) renderProjectsList() string   { return fmt.Sprintf("%d projects", len(m.data.Projects)) }
func (m model) renderProjectsDetail() string { return "Select a project" }
```

`internal/tui/skills.go`:
```go
package tui

import "fmt"

func (m model) renderSkillsList() string   { return fmt.Sprintf("%d plugins", len(m.data.Plugins)) }
func (m model) renderSkillsDetail() string { return "Select a plugin" }
```

`internal/tui/hooks.go`:
```go
package tui

import "fmt"

func (m model) renderHooksList() string   { return fmt.Sprintf("%d hooks", len(m.data.Hooks)) }
func (m model) renderHooksDetail() string { return "Select a hook" }
```

- [ ] **Step 5: cmd/root.go에 데이터 수집 + TUI 실행 연결**

`cmd/root.go`의 `RunE`를 업데이트. `sync.WaitGroup`으로 독립적인 I/O를 병렬 처리한다:

```go
RunE: func(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	home, _ := os.UserHomeDir()
	claudeDir := filepath.Join(home, ".claude")

	data := tui.DashboardData{}
	var wg sync.WaitGroup

	// 병렬: Global settings
	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Settings, _ = claude.ParseSettings(filepath.Join(claudeDir, "settings.json"))
		data.LocalSettings, _ = claude.ParseSettings(filepath.Join(claudeDir, "settings.local.json"))
		data.ClaudeMD, _ = claude.ParseClaudeMD(filepath.Join(claudeDir, "CLAUDE.md"))
	}()

	// 병렬: Plugins & Skills
	var installed map[string]claude.InstalledPlugin
	wg.Add(1)
	go func() {
		defer wg.Done()
		installed, _ = claude.ParseInstalledPlugins(filepath.Join(claudeDir, "plugins", "installed_plugins.json"))
		data.LocalSkills, _ = claude.ScanLocalSkills(filepath.Join(claudeDir, "skills"))
	}()

	// 병렬: Projects meta
	wg.Add(1)
	go func() {
		defer wg.Done()
		data.ProjectsMeta, _ = claude.ScanProjectsMeta(filepath.Join(claudeDir, "projects"))
	}()

	// 병렬: Project directory scan
	wg.Add(1)
	go func() {
		defer wg.Done()
		data.Projects, _ = scanner.ScanProjects(cfg.Scan.Roots, cfg.Scan.Ignore)
	}()

	wg.Wait()

	// 병렬 완료 후 의존성 있는 작업
	var enabled map[string]bool
	if data.Settings != nil {
		enabled = data.Settings.EnabledPlugins
		data.Hooks = claude.ExtractHooks(data.Settings, "global")
	}
	data.Plugins = claude.MergePluginData(installed, enabled)

	return tui.Run(data)
},
```

필요한 import 추가: `"os"`, `"path/filepath"`, 그리고 내부 패키지들.

- [ ] **Step 5: 빌드 및 수동 동작 확인**

Run:
```bash
go build -o ctv .
./ctv
```
Expected: TUI 화면이 뜨고 탭 전환(Tab/←→)과 종료(q)가 동작

- [ ] **Step 6: 커밋**

```bash
git add internal/tui/ cmd/root.go
git commit -m "feat: TUI 대시보드 기본 프레임 구현

- Bubble Tea 기반 메인 모델, 탭 전환, 키바인딩
- lipgloss 스타일 정의
- cmd/root.go에서 데이터 수집 후 TUI 실행 연결"
```

---

### Task 6: TUI — Global 탭 렌더링

**Files:**
- Create: `internal/tui/global.go`
- Modify: `internal/tui/app.go` (placeholder 교체)

- [ ] **Step 1: global.go 구현 (stub 파일 내용을 교체)**

```go
package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/claude"
)

var globalItems = []string{"Settings", "Local Settings", "CLAUDE.md"}

func (m model) renderGlobalList() string {
	var b strings.Builder
	for i, item := range globalItems {
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + item))
		} else {
			b.WriteString(listItemStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) renderGlobalDetail() string {
	switch m.cursor {
	case 0: // Settings
		return m.renderSettingsDetail(m.data.Settings, "settings.json")
	case 1: // Local Settings
		return m.renderSettingsDetail(m.data.LocalSettings, "settings.local.json")
	case 2: // CLAUDE.md
		return m.renderClaudeMDDetail()
	}
	return ""
}

func (m model) renderSettingsDetail(s *claude.Settings, name string) string {
	if s == nil {
		return detailStyle.Render(name + " not found")
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(name) + "\n\n")

	if s.Model != "" {
		b.WriteString(fmt.Sprintf("  model:     %s\n", s.Model))
	}
	if s.Language != "" {
		b.WriteString(fmt.Sprintf("  language:  %s\n", s.Language))
	}
	if s.TeammateMode != "" {
		b.WriteString(fmt.Sprintf("  teammate:  %s\n", s.TeammateMode))
	}

	if len(s.Env) > 0 {
		b.WriteString("\n  env:\n")
		for k, v := range s.Env {
			b.WriteString(fmt.Sprintf("    %s: %s\n", k, v))
		}
	}

	if len(s.Permissions.Allow) > 0 {
		b.WriteString(fmt.Sprintf("\n  permissions.allow: (%d rules)\n", len(s.Permissions.Allow)))
		for _, p := range s.Permissions.Allow {
			b.WriteString(fmt.Sprintf("    %s %s\n", statusIcon, p))
		}
	}
	if len(s.Permissions.Deny) > 0 {
		b.WriteString(fmt.Sprintf("\n  permissions.deny: (%d rules)\n", len(s.Permissions.Deny)))
		for _, p := range s.Permissions.Deny {
			b.WriteString(fmt.Sprintf("    %s %s\n", statusIconOff, p))
		}
	}

	return b.String()
}

func (m model) renderClaudeMDDetail() string {
	md := m.data.ClaudeMD
	if md == nil {
		return detailStyle.Render("CLAUDE.md not found")
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("CLAUDE.md") + "\n\n")
	b.WriteString(fmt.Sprintf("  Lines: %d\n\n", md.LineCount))

	if len(md.Sections) > 0 {
		b.WriteString("  Sections:\n")
		for _, s := range md.Sections {
			b.WriteString(fmt.Sprintf("    • %s\n", s))
		}
	}

	return b.String()
}
```

- [ ] **Step 2: 빌드 후 수동 확인**

Run: `go build -o ctv . && ./ctv`
Expected: Global 탭에서 Settings/Local Settings/CLAUDE.md 선택 시 상세 표시

- [ ] **Step 4: 커밋**

```bash
git add internal/tui/
git commit -m "feat: Global 탭 렌더링 구현

- Settings, Local Settings, CLAUDE.md 목록 및 상세 표시
- model, language, permissions, env 등 포맷팅 출력"
```

---

### Task 7: TUI — Projects 탭 렌더링

**Files:**
- Create: `internal/tui/projects.go`
- Modify: `internal/tui/app.go` (placeholder 교체)

- [ ] **Step 1: projects.go 구현**

```go
package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/claude"
	"github.com/nopecho/claude-television/internal/scanner"
)

func (m model) renderProjectsList() string {
	if len(m.data.Projects) == 0 {
		return listItemStyle.Render("No projects found.\nUse: ctv scan <path>")
	}

	var b strings.Builder
	for i, p := range m.data.Projects {
		icon := statusIconOff
		if p.HasClaudeDir || p.HasClaudeMD {
			icon = statusIcon
		}

		line := fmt.Sprintf("%s %s", icon, p.Name)
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + line))
		} else {
			b.WriteString(listItemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) renderProjectsDetail() string {
	if len(m.data.Projects) == 0 || m.cursor >= len(m.data.Projects) {
		return ""
	}

	p := m.data.Projects[m.cursor]
	var b strings.Builder
	b.WriteString(titleStyle.Render(p.Name) + "\n")
	b.WriteString(fmt.Sprintf("  %s\n\n", p.Path))

	b.WriteString(fmt.Sprintf("  .claude/          %s\n", boolIcon(p.HasClaudeDir)))
	b.WriteString(fmt.Sprintf("  CLAUDE.md         %s\n", boolIcon(p.HasClaudeMD)))
	b.WriteString(fmt.Sprintf("  settings.json     %s\n", boolIcon(p.HasSettings)))

	// Check project meta from ~/.claude/projects/
	meta := m.findProjectMeta(p.Path)
	if meta != nil {
		b.WriteString(fmt.Sprintf("\n  memory            %s\n", boolIcon(meta.HasMemory)))
		b.WriteString(fmt.Sprintf("  sessions          %s\n", boolIcon(meta.HasSessions)))
	}

	// Parse and show project-level CLAUDE.md if exists
	if p.HasClaudeMD {
		md, err := claude.ParseClaudeMD(p.Path + "/CLAUDE.md")
		if err == nil {
			b.WriteString(fmt.Sprintf("\n  CLAUDE.md (%d lines):\n", md.LineCount))
			for _, s := range md.Sections {
				b.WriteString(fmt.Sprintf("    • %s\n", s))
			}
		}
	}

	return b.String()
}

func (m model) findProjectMeta(path string) *claude.ProjectMeta {
	for i := range m.data.ProjectsMeta {
		if m.data.ProjectsMeta[i].DecodedPath == path {
			return &m.data.ProjectsMeta[i]
		}
	}
	return nil
}

func boolIcon(b bool) string {
	if b {
		return statusIcon
	}
	return statusIconOff
}
```

- [ ] **Step 2: 빌드 후 수동 확인**

Run: `go build -o ctv . && ./ctv`
Expected: Projects 탭에서 프로젝트 목록 + 상세 표시

- [ ] **Step 3: 커밋**

```bash
git add internal/tui/
git commit -m "feat: Projects 탭 렌더링 구현

- 프로젝트 목록에 Claude 설정 유무 아이콘 표시
- 상세에서 .claude/, CLAUDE.md, settings.json, memory, sessions 표시"
```

---

### Task 8: TUI — Skills 탭 렌더링

**Files:**
- Create: `internal/tui/skills.go`
- Modify: `internal/tui/app.go` (placeholder 교체)

- [ ] **Step 1: skills.go 구현**

```go
package tui

import (
	"fmt"
	"strings"
)

func (m model) renderSkillsList() string {
	var b strings.Builder

	items := m.skillCache
	if len(items) == 0 {
		return listItemStyle.Render("No plugins or skills found.")
	}

	for i, item := range items {
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + item.display))
		} else {
			b.WriteString(listItemStyle.Render("  " + item.display))
		}
		b.WriteString("\n")
	}
	return b.String()
}

type skillItem struct {
	display    string
	isPlugin   bool
	pluginIdx  int
	skillIdx   int
}

func (m model) buildSkillItems() []skillItem {
	var items []skillItem
	for i, p := range m.data.Plugins {
		icon := statusIconOff
		if p.Enabled {
			icon = statusIcon
		}
		tag := "plugin"
		if !p.Installed {
			tag = "not installed"
		}
		items = append(items, skillItem{
			display:   fmt.Sprintf("%s %s [%s]", icon, p.Name, tag),
			isPlugin:  true,
			pluginIdx: i,
		})
	}
	for i, s := range m.data.LocalSkills {
		items = append(items, skillItem{
			display:  fmt.Sprintf("%s %s [local]", statusIcon, s.Name),
			isPlugin: false,
			skillIdx: i,
		})
	}
	return items
}

func (m model) renderSkillsDetail() string {
	items := m.skillCache
	if len(items) == 0 || m.cursor >= len(items) {
		return ""
	}

	item := items[m.cursor]
	var b strings.Builder

	if item.isPlugin {
		p := m.data.Plugins[item.pluginIdx]
		b.WriteString(titleStyle.Render(p.Name) + "\n\n")
		b.WriteString(fmt.Sprintf("  key:         %s\n", p.Key))
		b.WriteString(fmt.Sprintf("  marketplace: %s\n", p.Marketplace))
		b.WriteString(fmt.Sprintf("  version:     %s\n", p.Version))
		b.WriteString(fmt.Sprintf("  enabled:     %s\n", boolIcon(p.Enabled)))
		b.WriteString(fmt.Sprintf("  installed:   %s\n", boolIcon(p.Installed)))
		if p.InstallPath != "" {
			b.WriteString(fmt.Sprintf("  path:        %s\n", p.InstallPath))
		}
	} else {
		s := m.data.LocalSkills[item.skillIdx]
		b.WriteString(titleStyle.Render(s.Name) + "\n\n")
		b.WriteString(fmt.Sprintf("  type: local skill\n"))
		b.WriteString(fmt.Sprintf("  path: %s\n", s.Path))
	}

	return b.String()
}
```

- [ ] **Step 2: 빌드 후 수동 확인**

Run: `go build -o ctv . && ./ctv`
Expected: Skills 탭에서 플러그인/스킬 목록 + 상세 표시

- [ ] **Step 3: 커밋**

```bash
git add internal/tui/
git commit -m "feat: Skills 탭 렌더링 구현

- 플러그인 목록 (enabled/disabled, installed 상태)
- 로컬 스킬 표시
- 상세에서 key, marketplace, version, path 표시"
```

---

### Task 9: TUI — Hooks 탭 렌더링

**Files:**
- Create: `internal/tui/hooks.go`
- Modify: `internal/tui/app.go` (placeholder 교체)

- [ ] **Step 1: hooks.go 구현**

```go
package tui

import (
	"fmt"
	"strings"
)

func (m model) renderHooksList() string {
	if len(m.data.Hooks) == 0 {
		return listItemStyle.Render("No hooks registered.")
	}

	var b strings.Builder
	for i, h := range m.data.Hooks {
		line := fmt.Sprintf("%s [%s]", h.Event, h.Source)
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + line))
		} else {
			b.WriteString(listItemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) renderHooksDetail() string {
	if len(m.data.Hooks) == 0 || m.cursor >= len(m.data.Hooks) {
		return ""
	}

	h := m.data.Hooks[m.cursor]
	var b strings.Builder
	b.WriteString(titleStyle.Render(h.Event) + "\n\n")
	b.WriteString(fmt.Sprintf("  event:   %s\n", h.Event))
	b.WriteString(fmt.Sprintf("  type:    %s\n", h.Type))
	if h.Matcher != "" {
		b.WriteString(fmt.Sprintf("  matcher: %s\n", h.Matcher))
	}
	b.WriteString(fmt.Sprintf("  command: %s\n", h.Command))
	if h.Async {
		b.WriteString(fmt.Sprintf("  async:   %v\n", h.Async))
	}
	if h.Timeout > 0 {
		b.WriteString(fmt.Sprintf("  timeout: %ds\n", h.Timeout))
	}
	b.WriteString(fmt.Sprintf("  source:  %s\n", h.Source))

	return b.String()
}
```

- [ ] **Step 2: 빌드 후 전체 수동 확인**

Run: `go build -o ctv . && ./ctv`
Expected: 4개 탭 모두 정상 동작 (목록 + 상세 렌더링, 탭 전환, 커서 이동)

- [ ] **Step 4: 커밋**

```bash
git add internal/tui/
git commit -m "feat: Hooks 탭 렌더링 및 TUI placeholder 제거

- Hooks 탭 목록/상세 렌더링 구현
- 4개 탭 모두 실제 렌더링으로 전환"
```

---

### Task 10: 최종 검증 및 정리

**Files:**
- Modify: `README.md` (필요 시 업데이트)
- Modify: `CLAUDE.md` (필요 시 업데이트)

- [ ] **Step 1: 전체 테스트 실행**

Run: `go test ./... -v`
Expected: 모든 테스트 PASS

- [ ] **Step 2: 빌드 확인**

Run: `go build -o ctv .`
Expected: 바이너리 정상 생성

- [ ] **Step 3: E2E 수동 검증**

```bash
./ctv version
./ctv scan ~/projects
./ctv scan --list
./ctv
# TUI에서:
# - Global 탭: Settings, Local Settings, CLAUDE.md 상세 확인
# - Projects 탭: 프로젝트 목록 + 설정 현황 확인
# - Skills 탭: 플러그인/스킬 목록 확인
# - Hooks 탭: hooks 목록 확인
# - 탭 전환, 커서 이동, 종료 동작 확인
./ctv scan --remove ~/projects
```

- [ ] **Step 4: go vet, go mod tidy 정리**

Run:
```bash
go vet ./...
go mod tidy
```

- [ ] **Step 5: 최종 커밋**

```bash
git add -A
git commit -m "chore: 최종 정리 및 검증

- go vet 통과
- go mod tidy 정리
- 전체 테스트 통과 확인"
```
