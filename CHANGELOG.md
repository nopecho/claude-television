# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## 0.0.0-pre (2026-03-21)


### Features

* Cobra CLI 프레임워크 구현 ([f639c75](https://github.com/nopecho/claude-television/commit/f639c7529e9fd3ce2726a498727b5a8e50f5db6d))
* plugins, skills, hooks, projects 파싱 구현 ([9f727f2](https://github.com/nopecho/claude-television/commit/9f727f2d8592bb579459e57b4737f93c1632f262))
* settings.json 및 CLAUDE.md 파싱 구현 ([44e21c7](https://github.com/nopecho/claude-television/commit/44e21c7752f5b262ca79240673520a42d9f80186))
* setup automated release pipeline and Korean README ([8c09a20](https://github.com/nopecho/claude-television/commit/8c09a20e236ed9405930fa4e2f9db285310299b3))
* TUI 대시보드 구현 (4개 탭 렌더링) ([b7db451](https://github.com/nopecho/claude-television/commit/b7db45128f4d117c9de09a5565af41f9e68d8d1f))
* 병렬 데이터 수집 및 TUI 실행 연결 ([8423876](https://github.com/nopecho/claude-television/commit/8423876153baa3a534480a8fdfa67826a27e49b2))
* 프로젝트 디렉토리 스캐너 구현 ([7bb813c](https://github.com/nopecho/claude-television/commit/7bb813cad85cf5e67e68c6e2a40053b53876249d))


### Miscellaneous Chores

* force release as v0.0.0-pre ([04fc783](https://github.com/nopecho/claude-television/commit/04fc7833f649ea47c19c220d75a69df5d80ef375))

## [Unreleased]

### Added

- Initial project setup
- TUI dashboard with 4 tabs: Global, Projects, Skills, Hooks
- `ctv scan` command for managing project scan paths
- Read-only view of Claude Code settings, CLAUDE.md, plugins, and hooks
