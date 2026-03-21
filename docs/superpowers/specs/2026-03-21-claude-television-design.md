# claude-television (ctv) — 설계 문서

## 개요

Claude Code의 로컬 설정 상태를 TUI 대시보드로 한눈에 탐색하는 읽기 전용 CLI 도구.

- **이름:** claude-television (ctv)
- **언어:** Go + Cobra (CLI) + Bubble Tea (TUI)
- **버전 관리:** asdf
- **성격:** 읽기 전용 대시보드, 설정 변경 기능 없음

## 해결하는 문제

Claude Code 설정은 여러 위치에 분산되어 있다:
- `~/.claude/settings.json`, `settings.local.json` (글로벌)
- `~/.claude/CLAUDE.md` (글로벌 지시)
- `~/.claude/plugins/` (플러그인/스킬)
- `~/.claude/projects/` (프로젝트별 세션, 메모리)
- 각 프로젝트 디렉토리의 `.claude/settings.json`, `CLAUDE.md` (프로젝트 로컬)

현재는 이 파일들을 직접 열어 확인해야 하며, 어떤 프로젝트에 어떤 설정이 적용되어 있는지 한눈에 볼 방법이 없다.

## 커맨드

```
ctv                     # TUI 대시보드 바로 실행
ctv scan <path>         # 프로젝트 스캔 경로 등록
ctv scan --list         # 등록된 스캔 경로 목록
ctv scan --remove <path># 스캔 경로 제거
ctv version             # 버전 출력
```

## TUI 대시보드 구조

### 레이아웃

```
┌─ claude-television ──────────────────────────────────────┐
│                                                          │
│  [Global] [Projects] [Skills] [Hooks]          ← 탭     │
│                                                          │
│  ┌─ 목록 패널 ─────────┐  ┌─ 상세 패널 ──────────────┐  │
│  │                      │  │                          │  │
│  │  (탭별 목록)         │  │  (선택 항목 상세)        │  │
│  │                      │  │                          │  │
│  └──────────────────────┘  └──────────────────────────┘  │
│                                                          │
│  ↑↓ navigate  ←→/Tab switch  Enter select  q quit       │
└──────────────────────────────────────────────────────────┘
```

### 탭별 내용

#### Global 탭
글로벌 Claude Code 설정의 전체 현황.

- **Settings** — `~/.claude/settings.json` 내용 (model, language, env, permissions 등)
- **Local Settings** — `~/.claude/settings.local.json` 내용 (오버라이드)
- **CLAUDE.md** — 글로벌 CLAUDE.md 내용 미리보기 (라인 수, 섹션 목록)

목록 패널: 세 항목(Settings, Local Settings, CLAUDE.md)을 나열.
상세 패널: 선택한 항목의 내용을 포맷팅하여 표시.

#### Projects 탭
스캔된 디렉토리에서 발견한 프로젝트 목록과 각 프로젝트의 Claude Code 설정 현황.

- 프로젝트별 Claude Code 설정 존재 여부 표시 (아이콘)
- 선택 시 해당 프로젝트의 상세 표시:
  - `.claude/settings.json` 존재 여부 및 내용
  - `CLAUDE.md` 존재 여부, 라인 수, 섹션 헤더 미리보기
  - `~/.claude/projects/` 하위 프로젝트 메모리/세션 존재 여부

목록 패널: 프로젝트 이름 + 설정 존재 아이콘 (✓/✗).
상세 패널: 선택한 프로젝트의 설정 상세.

#### Skills 탭
설치된 플러그인/스킬의 목록과 상태.

- `~/.claude/settings.json`의 `enabledPlugins` 기반
- `~/.claude/plugins/cache/` 에서 실제 설치 정보 확인
- 플러그인별: 이름, 마켓플레이스, 버전, 활성화 상태

목록 패널: 플러그인 목록 + 활성화 상태.
상세 패널: 선택한 플러그인의 버전, 마켓플레이스, 포함된 스킬 목록.

#### Hooks 탭
등록된 hooks의 전체 목록.

- 글로벌 hooks (`~/.claude/settings.json` 내 hooks 설정)
- 프로젝트별 hooks (각 프로젝트 `.claude/settings.json` 내 hooks)
- hook별: 이벤트 타입, matcher, command

목록 패널: hook 목록 (글로벌/프로젝트 구분).
상세 패널: 선택한 hook의 이벤트, matcher, command 상세.

## 설정 파일

`~/.config/ctv/config.yaml`:

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

`ctv scan <path>` 커맨드로 이 파일을 관리한다.

## 프로젝트 구조

```
claude-television/
├── cmd/
│   ├── root.go          # cobra root 커맨드 (기본 = dashboard)
│   ├── scan.go          # scan 서브커맨드
│   └── version.go       # version 서브커맨드
├── internal/
│   ├── config/
│   │   └── config.go    # ctv 자체 설정 파일 (scan roots 등) 로드/저장
│   ├── claude/
│   │   ├── settings.go  # settings.json 파싱
│   │   ├── claudemd.go  # CLAUDE.md 파싱 (섹션 추출)
│   │   ├── plugins.go   # plugins/cache 스캔, enabledPlugins 매핑
│   │   ├── hooks.go     # hooks 설정 파싱
│   │   └── projects.go  # ~/.claude/projects/ 스캔
│   ├── scanner/
│   │   └── scanner.go   # 프로젝트 디렉토리 스캔, Claude 설정 탐지
│   └── tui/
│       ├── app.go       # bubbletea 메인 모델
│       ├── global.go    # Global 탭 컴포넌트
│       ├── projects.go  # Projects 탭 컴포넌트
│       ├── skills.go    # Skills 탭 컴포넌트
│       ├── hooks.go     # Hooks 탭 컴포넌트
│       ├── tabs.go      # 탭 네비게이션
│       └── styles.go    # lipgloss 스타일 정의
├── main.go
├── go.mod
├── go.sum
├── .tool-versions       # asdf 버전 관리 (golang)
└── README.md
```

## 의존성

| 라이브러리 | 용도 |
|-----------|------|
| `github.com/spf13/cobra` | CLI 프레임워크 |
| `github.com/charmbracelet/bubbletea` | TUI 프레임워크 |
| `github.com/charmbracelet/lipgloss` | TUI 스타일링 |
| `github.com/charmbracelet/bubbles` | TUI 컴포넌트 (list, viewport 등) |
| `gopkg.in/yaml.v3` | ctv 설정 파일 파싱 |

## 데이터 흐름

```
1. ctv 실행
2. ~/.config/ctv/config.yaml 로드 → 스캔 경로 확인
3. 병렬로 데이터 수집:
   a. ~/.claude/settings.json, settings.local.json 파싱
   b. ~/.claude/CLAUDE.md 파싱
   c. ~/.claude/plugins/ 스캔
   d. ~/.claude/projects/ 스캔
   e. 등록된 경로에서 프로젝트 디렉토리 스캔
4. 수집된 데이터로 TUI 모델 초기화
5. 대시보드 렌더링
```

## 키바인딩

| 키 | 동작 |
|----|------|
| `↑`/`k` | 목록에서 위로 이동 |
| `↓`/`j` | 목록에서 아래로 이동 |
| `Tab`/`←`/`→` | 탭 전환 |
| `Enter` | 항목 선택/상세 토글 |
| `/` | 목록 필터링 (프로젝트 검색) |
| `q`/`Ctrl+C` | 종료 |

## 성공 기준

1. `ctv` 실행 시 1초 이내에 대시보드가 표시된다
2. 글로벌 settings, CLAUDE.md 내용을 정확히 파싱하여 보여준다
3. 스캔 경로의 프로젝트를 탐지하고 Claude 설정 유무를 표시한다
4. 설치된 플러그인/스킬 목록과 활성화 상태를 보여준다
5. hooks 설정을 글로벌/프로젝트 구분하여 보여준다
6. 탭 전환과 목록 탐색이 부드럽게 동작한다

## MVP 범위

**포함:**
- Global 탭 (settings, local settings, CLAUDE.md)
- Projects 탭 (프로젝트 스캔, 설정 존재 여부, 상세)
- Skills 탭 (플러그인 목록, 활성화 상태)
- Hooks 탭 (글로벌 hooks)
- `ctv scan` 커맨드

**제외 (향후):**
- 설정 편집 기능
- 프로젝트별 hooks 개별 표시
- 실시간 파일 감시 (auto-refresh)
- 원격 프로젝트 지원
