# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Core mining engine based on NerdMiner original logic
- Bubbletea TUI with dashboard and settings screens
- Timeseries line chart for dynamic hashrate display
- Stratum pool client with automatic difficulty suggestion for CPU mining

### Fixed
- Fixed endianness bug where nonces were sent to the pool in Little-Endian byte format rather than standard hex string format, causing rejected shares.
