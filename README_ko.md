# claude-television

[English 🇺🇸](./README.md) | [한국어 🇰🇷](./README_ko.md)

> 📺 한눈에 파악하는 Claude Code 설정 TUI 대시보드.

Claude Code 설정은 `settings.json`, `CLAUDE.md`, 플러그인, 훅, 프로젝트별 설정 등 여러 위치에 분산되어 있습니다. **claude-television** (`ctv`)은 이 모든 설정을 읽기 전용 터미널 대시보드 하나로 모아 보여줍니다.

## 주요 기능

- **전역 설정 (Global Settings)** — `~/.claude/settings.json`, `settings.local.json`, 전역 `CLAUDE.md`를 한 곳에서 확인합니다.
- **프로젝트 탐색기 (Project Explorer)** — 디렉토리를 스캔하여 로컬 `.claude/` 디렉토리나 `CLAUDE.md` 파일이 존재하는 프로젝트들을 한눈에 보여줍니다.
- **스킬 및 플러그인 (Skills & Plugins)** — 설치된 스킬과 플러그인 목록, 버전, 경로 및 활성화 상태 등의 세부 정보를 제공합니다.
- **훅 개요 (Hooks Overview)** — 실행 트리거 조건 및 연결된 쉘 스크립트 등 등록된 훅(Hook) 정보들을 직관적으로 검사할 수 있습니다.
- **키보드 기반 동작 (Keyboard-driven)** — Vim 스타일의 키 바인딩으로 터미널에서 빠르고 편리하게 탐색할 수 있습니다.

## 설치 방법

### 빠른 설치 (macOS / Linux)

가장 간단한 설치 방법은 터미널에서 아래 명령어를 실행하여 설치 스크립트를 사용하는 것입니다:

```bash
curl -fsSL https://raw.githubusercontent.com/nopecho/claude-television/main/scripts/install.sh | bash
```

또는 [GitHub Releases](https://github.com/nopecho/claude-television/releases) 페이지에서 직접 압축 파일을 다운로드할 수도 있습니다.

### 소스 코드에서 설치

Go 1.26 이상 필요:

```bash
go install github.com/nopecho/claude-television@latest
```

### asdf 사용

```bash
asdf install golang 1.26.1
asdf local golang 1.26.1
go install github.com/nopecho/claude-television@latest
```

## 빠른 시작

```bash
# 프로젝트를 스캔할 디렉토리 등록
ctv scan ~/projects

# 대시보드 실행
ctv
```

## 사용법

```
ctv                      # TUI 대시보드 실행
ctv scan <path>          # 스캔할 경로 등록
ctv scan --list          # 등록된 스캔 경로 목록 보기
ctv scan --remove <path> # 스캔 경로 제거
ctv version              # 버전 확인
```

## 대시보드 화면

```
┌─ claude-television ──────────────────────────────────────┐
│                                                          │
│  [Global] [Projects] [Skills] [Hooks]                    │
│                                                          │
│  ┌─ List ──────────────┐  ┌─ Detail ──────────────────┐  │
│  │ ● Settings       ✓  │  │ model: opus               │  │
│  │ ● Local Settings  ✓ │  │ language: korean          │  │
│  │ ● CLAUDE.md      ✓  │  │ permissions:              │  │
│  │                     │  │   allow: [Bash, Read...]  │  │
│  └─────────────────────┘  └───────────────────────────┘  │
│                                                          │
│  ↑↓/jk navigate  ←→/Tab switch  / filter  q quit         │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

## 단축키

| 키                | 동작                   |
| ----------------- | ---------------------- |
| `↑` / `k`         | 위로 이동              |
| `↓` / `j`         | 아래로 이동            |
| `Tab` / `←` / `→` | 탭 전환                |
| `Enter`           | 선택 / 상세 정보 토글  |
| `/`               | 목록 필터링 (MVP 이후) |
| `q` / `Ctrl+C`    | 종료                   |

## 설정

ctv는 설정 파일을 `~/.config/ctv/config.yaml`에 저장합니다:

```yaml
scan:
  roots:               # 프로젝트를 탐색할 최상위 디렉토리 목록
    - ~/projects
    - ~/work
  ignore:              # 스캔 시 제외할 디렉토리명 (성능 향상 목적)
    - node_modules
    - .git
    - vendor
```

- **`scan.roots`**: `ctv scan <path>` 명령어로 등록된 디렉토리들입니다. `ctv`는 이 디렉토리들을 재귀적으로 탐색하여 프로젝트들을 찾습니다.
- **`scan.ignore`**: 성능 향상을 위해 스캔 시 무시할 디렉토리 목록입니다. 기본값으로 `node_modules`나 `.git` 같이 크기가 큰 디렉토리들이 포함되어 있습니다. 설정 파일이 없는 경우 합리적인 기본값들이 사용됩니다.

## 기여하기

기여는 언제나 환영합니다! 변경하고 싶은 내용에 대해 먼저 이슈를 열어 논의해 주세요.

1. 저장소를 포크(Fork)합니다.
2. 기능 브랜치를 생성합니다 (`git checkout -b feat/amazing-feature`)
3. 변경 사항을 커밋합니다 ([Conventional Commits](https://www.conventionalcommits.org/) 준수)
4. 브랜치에 푸시합니다 (`git push origin feat/amazing-feature`)
5. Pull Request를 엽니다.

## 라이선스

MIT
