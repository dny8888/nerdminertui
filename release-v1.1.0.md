# Release v1.1.0 Plan

## Goal
Release version `v1.1.0` of NerdMiner TUI by committing pending improvements, untracking the compiled binary `bin/nerdtui`, updating `CHANGELOG.md` with semantic groups, and tag creation.

## Tasks
- [x] Task 1: Untrack `bin/nerdtui` from Git → Verify: `git status` shows `bin/nerdtui` deleted/untracked
- [x] Task 2: Commit pending Go files and `.gitignore` with message `feat: add stratum reconnect and hash optimizations` → Verify: `git status` shows only `CHANGELOG.md` modified (and this plan)
- [x] Task 3: Edit `CHANGELOG.md` to format `[Unreleased]` changes under `## [1.1.0] - 2026-06-04`, adding details for the new reconnect & hash optimizations → Verify: `CHANGELOG.md` updated
- [x] Task 4: Commit `CHANGELOG.md` change → Verify: `git status` shows working tree clean (excluding this plan)
- [x] Task 5: Run full test suite → Verify: `go test ./...` returns green
- [x] Task 6: Run `bin/release v1.1.0` script to create and push tag → Verify: Output shows tag `v1.1.0` created and pushed to origin

## Done When
- [x] Working tree is clean (excluding this plan)
- [x] Git tag `v1.1.0` is pushed to remote origin
