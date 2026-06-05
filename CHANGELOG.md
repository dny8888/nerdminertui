# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Interactive help overlay (`?` hotkey) to display global keyboard shortcuts.
- `ToastManager` system for temporary, non-blocking alert notifications (errors and status updates) on the UI.
- Persistent SQLite `store` integration for hashrate history tracking across application restarts.
- Centralized `i18n` dictionary (`pkg/i18n`) replacing hardcoded strings with Portuguese translations.

### Changed
- Refactored the core mining hash loop using `sync.Pool` and pre-calculated SHA256d mid-states to eliminate heap allocations, drastically increasing hashrate performance.
- Optimized Stratum parsing to reduce JSON unmarshaling overhead.
- Updated all UI components (Dashboard, Settings, Global Stats, Status Bar) to be fully responsive to terminal resizing (`tea.WindowSizeMsg`).
- Improved navigation in the Settings screen by supporting backwards cycling with `Shift+Tab`.
- Added visual validation (red borders and text) for malformed inputs in the Settings screen.

### Fixed
- Fixed UI uptime rendering bug where uptime failed to properly format days and appeared as `0m`.
## [1.1.0] - 2026-06-04
### Added
- Support for `client.reconnect` method in the Stratum pool client.
- Dynamic ExtraNonce2 entropy generator (`pkg/trivia`) using space/astronomy words to avoid nonce collisions.
- Multi-worker block header rebuilding with unique extranonce2.
- Benchmark tests for block hashing methods.

### Changed
- Refactored mining hash loop (`pkg/mining/hash.go`) to achieve zero heap allocations per hash, improving performance.

## [1.0.0] - 2026-06-02
### Added
- Core mining engine based on NerdMiner original logic
- Bubbletea TUI with dashboard and settings screens
- Timeseries line chart for dynamic hashrate display
- Stratum pool client with automatic difficulty suggestion for CPU mining

### Fixed
- Fixed endianness bug where nonces were sent to the pool in Little-Endian byte format rather than standard hex string format, causing rejected shares.
