# Cross-Environment Sandbox Runner - Development Plan

## Project Overview

**Project Name:** SafeRun  
**Type:** CLI tool for developer sandboxing  
**Core Feature:** Zero-config sandbox that automatically isolates student projects from the host system, preventing malicious file writes, unauthorized network calls, and accidental API key exposure.  
**Target Users:** Student developers, coding bootcamps, hackathon participants

---

## Architecture Decision

**Recommended Approach:** Hybrid container-based + language-level sandboxing

| Layer | Technology | Purpose |
|-------|------------|---------|
| Primary isolation | Docker/Podman containers | Process, filesystem, network isolation |
| Fallback/lightweight | Language wrappers (Python/Node) | Soft sandboxing when Docker unavailable |
| CLI interface | Go binary | Cross-platform, single executable |

### Why Hybrid?

- **When Docker is available:** Full container isolation with security flags
- **When Docker is unavailable:** Language-level wrapping with filesystem redirection + network proxy
- **Cross-platform:** Works on Linux, macOS, Windows (with Docker as optional dependency)

### Core Design Principles

1. **Zero-friction:** Students run `saferun python main.py` and it just works
2. **Security-first defaults:** Block filesystem outside project, block network by default
3. **Dev UX preserved:** Frontend hot-reload, localhost UIs work normally
4. **Fail gracefully:** If container fails, fallback to soft sandboxing

---

## File Structure

```
saferun/
├── cmd/
│   ├── saferun/
│   │   └── main.go           # CLI entry point
│   └── saferun-container/    # Container mode binary
├── pkg/
│   ├── sandbox/
│   │   ├── runner.go         # Core sandbox execution logic
│   │   ├── container.go      # Docker/podman integration
│   │   ├── wrapper.go        # Language-level sandboxing
│   │   └── config.go        # Configuration structs
│   ├── isolation/
│   │   ├── filesystem.go     # Filesystem guard logic
│   │   ├── network.go       # Network control logic
│   │   └── secrets.go       # Secrets detection & blocking
│   ├── detection/
│   │   ├── project.go       # Project type detection (Python/Node/etc)
│   │   ├── ports.go         # Port detection for UI exposure
│   │   └── env.go           # Environment variable scanner
│   └── ui/
│       ├── status.go        # Status display (blocked actions, warnings)
│       └── formatter.go     # Output formatting
├── internal/
│   ├── config/
│   │   └── config.go        # Config loading from CLI/env/flags
│   └── logging/
│       └── logging.go       # Structured logging
├── scripts/
│   ├── install.sh           # Installation script
│   ├── quickstart.sh        # Quick demo setup
│   └── test-projects/      # Test projects for validation
├── testdata/
│   ├── malicious/          # Test malicious packages
│   └── benign/             # Test benign projects
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile              # Build container for CLI
├── .dockerignore
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

---

## MVP Scope (2-3 Weeks)

### Week 1: Core Sandbox Execution

- [ ] **Day 1-2:** CLI skeleton with argument parsing  
  - `saferun <command>` wrapper that captures args  
  - Detect project type from directory contents

- [ ] **Day 3-4:** Container-based execution  
  - Build Docker command with safe flags  
  - Mount project folder read-write  
  - Set working directory

- [ ] **Day 5:** Filesystem isolation enforcement  
  - `--read-only` container flag  
  - Block writes outside project mount

### Week 2: Network & Security

- [ ] **Day 6-7:** Network control  
  - `--network=none` default  
  - Whitelist detection (npm/yarn registries)  
  - Port exposure for localhost UIs

- [ ] **Day 8-9:** Secrets protection  
  - Scan env vars for API key patterns  
  - Block or warn when secrets exposed  
  - Option to inject clean env into container

- [ ] **Day 10:** Logging & status output  
  - Show blocked actions  
  - Log all sandbox activity  
  - Verbose mode for debugging

### Week 3: Polish & Testing

- [ ] **Day 11-12:** Fallback wrapper mode  
  - If Docker unavailable, use language-level sandboxing  
  - Python: intercept `os` and `requests` calls via monkey-patching  
  - Node: wrapper script that redirects fs/network calls

- [ ] **Day 13:** Project detection & auto-config  
  - Auto-detect Python (requirements.txt, setup.py)  
  - Auto-detect Node (package.json)  
  - Auto-detect ports from config files

- [ ] **Day 14:** Final polish  
  - Error messages  
  - Installation script  
  - README with examples  
  - Basic test suite

---

## Implementation Details

### 1. CLI Entry Point (main.go)

```go
// Pseudocode structure
func main() {
    // Parse flags: --project-dir, --network-mode, --expose-ports, --verbose
    // Detect project type
    // Run sandbox (container or wrapper based on availability)
    // Forward stdin/stdout/stderr
    // Exit with proper code
}
```

**Change:** Add new file `cmd/saferun/main.go`

---

### 2. Container Runner (container.go)

```go
// BuildSafeDockerCmd builds Docker command with security flags
func BuildSafeDockerCmd(projectDir string, cmd []string) []string {
    return []string{
        "docker", "run", "--rm",
        "--read-only",                    // Filesystem guard
        "--cap-drop=ALL",                  // Drop all capabilities
        "--network=none",                  // Block network by default
        "--pids-limit=100",                // Limit processes
        "--memory=512m",                   // Memory limit
        "--security-opt=no-new-privileges",
        "-v", projectDir + ":/sandbox:rw",
        "-w", "/sandbox",
        "-e", "HOME=/sandbox",
        "saferun-base-image",
    }
}
```

**Change:** Add `pkg/sandbox/container.go`

---

### 3. Filesystem Guard (filesystem.go)

```go
// CheckPathAllowed verifies path is within project directory
func CheckPathAllowed(path string, projectDir string) bool {
    absProject, _ := filepath.Abs(projectDir)
    absPath, _ := filepath.Abs(path)
    return strings.HasPrefix(absPath, absProject)
}
```

**Change:** Add `pkg/isolation/filesystem.go`

---

### 4. Network Control (network.go)

```go
// AllowedHosts contains whitelisted registries
var AllowedHosts = []string{
    "registry.npmjs.org",
    "pypi.org",
    "pip.confederation.tech",
    "npmjs.org",
}

// IsNetworkAllowed checks if host is whitelisted
func IsNetworkAllowed(host string) bool {
    for _, allowed := range AllowedHosts {
        if strings.Contains(host, allowed) {
            return true
        }
    }
    return false
}
```

**Change:** Add `pkg/isolation/network.go`

---

### 5. Secrets Detection (secrets.go)

```go
// SensitivePatterns contains regex for API keys
var SensitivePatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)(api[_-]?key|apikey)['"]?\s*[:=]\s*['"]?[\w-]{20,}`),
    regexp.MustCompile(`(?i)secret['"]?\s*[:=]\s*['"]?[\w-]{20,}`),
    regexp.MustCompile(`(?i)password['"]?\s*[:=]\s*['"]?[\w-]{8,}`),
}

// ScanEnvForSecrets scans environment variables
func ScanEnvForSecrets(env []string) map[string]string {
    // Returns map of detected secrets
}
```

**Change:** Add `pkg/isolation/secrets.go`

---

### 6. Project Detection (project.go)

```go
// DetectProjectType determines project type from files
func DetectProjectType(dir string) ProjectType {
    if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
        return NodeProject
    }
    if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
        return PythonProject
    }
    // ... more project types
    return UnknownProject
}
```

**Change:** Add `pkg/detection/project.go`

---

### 7. Port Detection (ports.py)

```go
// DetectExposedPorts reads config files for common dev ports
func DetectExposedPorts(dir string) []int {
    ports := []int{}
    // Check package.json for "start": "react-scripts start" (3000)
    // Check package.json for "port" in config
    // Check .env for PORT=3000
    // Check webpack/vite config files
    return ports
}
```

**Change:** Add `pkg/detection/ports.go`

---

### 8. Language Wrapper Fallback (wrapper.go)

```go
// RunWithWrapper runs project with language-level sandboxing
func RunWithWrapper(projectType ProjectType, cmd []string) error {
    switch projectType {
    case PythonProject:
        return runPythonWrapper(cmd)
    case NodeProject:
        return runNodeWrapper(cmd)
    }
    return errors.New("unsupported project type for wrapper mode")
}
```

**Change:** Add `pkg/sandbox/wrapper.go`

---

### 9. Status UI (status.go)

```go
// LogBlockedAction logs and optionally displays blocked action
func LogBlockedAction(action string, details string) {
    log.Printf("[BLOCKED] %s: %s", action, details)
    if verboseMode {
        fmt.Println("⚠️  Blocked:", action)
    }
}
```

**Change:** Add `pkg/ui/status.go`

---

## Configuration Options

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--project-dir` | `.` | Project directory to sandbox |
| `--network` | `none` | Network mode: none, whitelist, full |
| `--expose-ports` | `auto` | Ports to expose (auto-detect or comma-separated) |
| `--verbose` | `false` | Show blocked actions in real-time |
| `--no-container` | `false` | Force wrapper mode (no Docker) |
| `--clean-env` | `false` | Strip all env vars before running |

### Environment Variables

- `SAFERUN_NETWORK_MODE` - Default network mode
- `SAFERUN_IMAGE` - Custom container image
- `SAFERUN_WHITELIST` - Comma-separated allowed hosts

---

## Security Defaults

### Container Flags (Always Applied)

```bash
--read-only
--cap-drop=ALL
--network=none
--pids-limit=100
--memory=512m
--security-opt=no-new-privileges
--user=nonroot
```

### What This Prevents

| Threat | Protection |
|--------|-------------|
| Malicious file writes | `--read-only` + project-only mount |
| Network exfiltration | `--network=none` (block all) |
| Privilege escalation | `--cap-drop=ALL`, `--user=nonroot` |
| Resource exhaustion | `--memory=512m`, `--pids-limit=100` |
| Secrets exposure | Env var scanning + clean env option |

---

## Frontend UI Handling

### Automatic Port Detection & Exposure

```go
// Example: React app typically uses port 3000
// saferun npm start
// → Auto-detects port 3000 from package.json
// → Runs: docker run -p 3000:3000 ...
// → User opens localhost:3000 in browser
```

### Hot Reload Support

- Project directory mounted read-write (`-v project:/sandbox:rw`)
- Frontend dev server (webpack, vite, etc.) works normally
- File changes inside container reflected to host
- No special configuration needed

---

## Testing Strategy

### Test Projects

```
testdata/
├── benign/
│   ├── python-flask-app/    # Flask with routes
│   ├── node-express/        # Express API
│   ├── react-app/           # Create-react-app
│   └── python-django/       # Django with runserver
└── malicious/
    ├── write-outside/       # Attempts to write to /tmp
    ├── network-exfil/      # Tries to send data out
    └── steal-env/          # Reads API keys from env
```

### Validation Tests

1. **Benign projects:** Run normally, verify UI accessible
2. **Malicious projects:** Verify writes blocked, network blocked, secrets protected
3. **Performance:** Measure overhead vs. native execution
4. **UX:** User testing with non-technical students

---

## Distribution

### Installation

```bash
# Quick install (macOS/Linux)
curl -sSL https://get.saferun.dev | sh

# Or via Go
go install github.com/saferun/saferun@latest
```

### Usage

```bash
# Python project
saferun python main.py

# Node project
saferun npm start

# With custom ports
saferun --expose-ports=3000,8080 npm start

# Allow specific network
saferun --network=whitelist npm install
```

---

## Future Enhancements (Post-MVP)

1. **Auto-detect secrets in code:** Scan source files for hardcoded keys
2. **GUI dashboard:** Show blocked actions, sandbox status
3. **Windows support:** Use Windows Sandbox or Job objects
4. **Template projects:** Pre-sandboxed boilerplates
5. **CI/CD integration:** Run tests in sandbox
6. **Multi-project workspaces:** Handle monorepos

---

## Summary

This plan delivers:

1. **Zero-friction sandboxing:** `saferun <command>` just works
2. **Security-first defaults:** Filesystem, network, secrets all protected
3. **Dev UX preserved:** Frontend UIs, hot-reload function normally
4. **Cross-platform:** Works on Linux/macOS (Windows post-MVP)
5. **Fail gracefully:** Container mode with wrapper fallback

The MVP can be built in 2-3 weeks by a solo developer with Go experience.

---

## Next Steps

1. **Confirm architecture:** Do we want the hybrid container+wrapper approach?
2. **Tech stack confirmation:** Go for CLI, Docker for primary isolation?
3. **MVP priority:** Which features are must-have vs. nice-to-have?
