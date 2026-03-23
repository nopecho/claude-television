# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.3.3-pre](https://github.com/nopecho/claude-television/compare/v0.3.2-pre...v0.3.3-pre) (2026-03-23)


### Bug Fixes

* TUI bug fixes and UI improvements ([#15](https://github.com/nopecho/claude-television/issues/15)) ([4c6c6d1](https://github.com/nopecho/claude-television/commit/4c6c6d1f504844e24888404e1deab3abaf066d15))

## [0.3.2-pre](https://github.com/nopecho/claude-television/compare/v0.3.1-pre...v0.3.2-pre) (2026-03-21)


### Bug Fixes

* trigger release-please to include tui-redesign ([e67d1ba](https://github.com/nopecho/claude-television/commit/e67d1ba4813bd3d2e2f7a9309ea4ac8bc3c024f4))

## [0.3.1-pre](https://github.com/nopecho/claude-television/compare/v0.3.0-pre...v0.3.1-pre) (2026-03-21)


### Bug Fixes

* dot-prefixed path decoding and global data deduplication ([#11](https://github.com/nopecho/claude-television/issues/11)) ([3b3060a](https://github.com/nopecho/claude-television/commit/3b3060ab66cb07dfdb11dc9677c502d6dd90a320))

## [0.3.0-pre](https://github.com/nopecho/claude-television/compare/v0.2.0-pre...v0.3.0-pre) (2026-03-21)


### Features

* comprehensive improvements and new features ([#9](https://github.com/nopecho/claude-television/issues/9)) ([019c049](https://github.com/nopecho/claude-television/commit/019c049b61e0091e3cfe5f2e95c7503659136f31))

## [0.1.0-pre](https://github.com/nopecho/claude-television/compare/v0.0.1-pre...v0.1.0-pre) (2026-03-21)


### Features

* introduce channel system with hybrid TUI dashboard ([#3](https://github.com/nopecho/claude-television/issues/3)) ([6c1a8a9](https://github.com/nopecho/claude-television/commit/6c1a8a952c7772ed42987d58dfe58a22773f8ec6))

## [0.0.1-pre](https://github.com/nopecho/claude-television/compare/v0.0.0-pre...v0.0.1-pre) (2026-03-21)


### Bug Fixes

* downgrade goreleaser config version to 1 ([aaba9e3](https://github.com/nopecho/claude-television/commit/aaba9e34fa474252e7acf696fadaba4b7be737a7))

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
