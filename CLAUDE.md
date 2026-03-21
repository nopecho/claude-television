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
