# internal/config, Technical Design

> Design specification for the `internal/config` module. Focuses on HOW the config loading is structured.

## Interface

### Classes / Functions

| Symbol | Signature | Return | Observation |
|---------|-----------|---------|------------|
| `Load` | `func Load() (*Config, error)` | `(*Config, error)` | Load properties from files and environment overrides. Performs tilde expansion on StorePath. 🟢 |
| `Config.Validate` | `func (c *Config) Validate() error` | `error` | Validate structural boundaries of loaded configs. 🟢 |
| `ExpandPath` | `func ExpandPath(p string) (string, error)` | `(string, error)` | Replaces leading `~/` with `os.UserHomeDir()`. Located in `internal/config/paths.go`. 🟢 |

## Main Flow
1. **Bind Environment Variables**: Use Viper to register bindings for env variables starting with `NM_` (e.g. `NM_POOL_URL` mappings). 🟢
2. **Apply Defaults**: Register default fallbacks:
   - `PoolURL`: `"public-pool.io"`
   - `PoolPort`: `21496`
   - `PollInterval`: `5s`
   - `CPUTarget`: `0.5`
   - `StorePath`: `"~/.nerdtui/metrics.db"`
   - `Theme`: `"dark"`
   - `MockMining`: `false` 🟢
3. **Load files**: Try loading settings from standard config files (e.g. `config.yaml` or `.env.local` if available). 🟢
4. **Expand Tilde Path (REQ-CONFIG-PATH-01)**: Call `ExpandPath` on `StorePath`. If an error occurs (such as an inaccessible home directory), abort and propagate the error. 🟢
5. **Validation Check**: Execute `c.Validate()` to check boundaries. 🟢

## Alternative Flows
- **Error on Validation Fail**: Returns a non-nil error if `BTCAddress == ""` (with `MockMining == false`) or if `CPUTarget` violates limits. 🟢
- **Home Directory Error**: If `os.UserHomeDir()` fails inside `ExpandPath`, propagate the error up to abort initialization. 🟢

## Dependencies
- `github.com/spf13/viper`: Used for bindings and parsing. 🟢
- `os`: Accesses user directories (`os.UserHomeDir`). 🟢

## Design Decisions Identified

| Decision | Evidence | Confidence |
|---------|---------------------|-----------|
| Env Precedence | `config.go:viper.AutomaticEnv()` | 🟢 |
| Strict Constraints | `config.go:Validate()` | 🟢 |
| Path Pre-processing | `paths.go:ExpandPath` | 🟢 |

## Internal State
- This module is stateless. It generates a config struct at startup and remains immutable throughout the app lifecycle. 🟢
