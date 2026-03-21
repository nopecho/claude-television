# claude-television

[English 🇺🇸](./README.md) | [한국어 🇰🇷](./README_ko.md)

> 📺 한눈에 파악하는 Claude Code 설정 TUI 대시보드.

Claude Code 설정은 `settings.json`, `CLAUDE.md`, 플러그인, 훅, MCP 서버, 메모리 등 여러 위치에 분산되어 있습니다. **claude-television** (`ctv`)은 새롭게 도입된 채널(Channel) 시스템을 통해 이 모든 정보를 하나의 터미널 대시보드에서 직관적으로 보여줍니다.

## 주요 기능

- **하이브리드 TUI 대시보드** — 좌측의 채널 목록과 우측의 상세 탭으로 구성된 분할 화면을 제공합니다.
- **채널 자동 탐색** — `~/.claude/projects/`에 등록된 프로젝트들을 자동으로 찾아 채널로 동기화합니다.
- **통합 정보 제공** — 전역/로컬 설정뿐만 아니라 `CLAUDE.md`, MCP 서버, 플러그인, 로컬 스킬, 프로젝트 상태(Health), 메모리 파일, Git 상태 등 종합적인 컨텍스트를 확인할 수 있습니다.
- **빠른 캐싱 시스템** — mtime 기반 캐싱을 통해 대시보드 로딩 속도를 최적화합니다.
- **키보드 기반 동작 (Keyboard-driven)** — Vim 스타일 단축키를 지원하며, 선택한 프로젝트 디렉토리로 즉시 이동(`cd`)할 수도 있습니다.

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
# ~/.claude/projects/에서 프로젝트 채널 자동 탐색 및 등록
ctv init

# 대시보드 실행
ctv
```

### 디렉토리 즉시 이동 기능 활성화
TUI 화면에서 `Alt+Enter`를 눌러 해당 프로젝트 디렉토리로 즉시 이동(`cd`)하려면, 사용 중인 쉘 설정 파일(`~/.zshrc` 또는 `~/.bashrc`)에 아래 함수를 추가하세요:

```bash
ctv() { local dir; dir="$(command ctv "$@")"; [ -d "$dir" ] && cd "$dir" || command ctv "$@"; }
```

## 사용법

```
ctv          # TUI 대시보드 실행
ctv init     # 채널 자동 탐색 및 등록
ctv version  # 버전 확인
```

## 대시보드 화면

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

## 단축키

| 키                  | 동작                           |
| ------------------- | ------------------------------ |
| `↑` / `k`           | 채널 목록 위로 이동            |
| `↓` / `j`           | 채널 목록 아래로 이동          |
| `Tab` / `→` / `l`   | 다음 상세 탭으로 전환          |
| `Shift+Tab`/`←`/`h` | 이전 상세 탭으로 전환          |
| `Ctrl+d` / `u`      | 상세 탭 내용 스크롤 (위/아래)  |
| `Alt+Enter`         | 프로젝트 디렉토리로 이동 (`cd`)|
| `/`                 | 채널 이름 퍼지 검색            |
| `?`                 | 상세 내용 내에서 검색          |
| `p`                 | 채널 고정 / 해제 (Pin)         |
| `e`                 | 프로젝트 설정 편집 ($EDITOR)   |
| `g`                 | 채널 그룹 관리 (Group)         |
| `q` / `Ctrl+C`      | 종료                           |

## 설정

ctv는 이제 `~/.config/ctv/config.json` 파일을 통해 설정을 관리합니다:

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
  },
  "editor": "vim"
}
```

- **`auto_sync`**: 실행 시 새로운 프로젝트를 자동으로 탐색합니다.
- **`cache_ttl`**: 채널 데이터 캐시 유지 시간입니다.
- **`pins`**: 목록 상단에 고정할 채널 이름이나 ID 배열입니다.
- **`groups`**: 채널 이름이나 ID를 지정하여 그룹 단위로 분류합니다.
- **`editor`**: 설정 파일을 열 때 사용할 에디터 명령어 (예: `vim`, `code`, `nano`).

## 기여하기

기여는 언제나 환영합니다! 변경하고 싶은 내용에 대해 먼저 이슈를 열어 논의해 주세요.

1. 저장소를 포크(Fork)합니다.
2. 기능 브랜치를 생성합니다 (`git checkout -b feat/amazing-feature`)
3. 변경 사항을 커밋합니다 ([Conventional Commits](https://www.conventionalcommits.org/) 준수)
4. 브랜치에 푸시합니다 (`git push origin feat/amazing-feature`)
5. Pull Request를 엽니다.

## 라이선스

MIT
