# Cocoon

Zero-config sandbox for running student projects safely.

Your projects run in their own space. Nothing can touch your machine or network.

## What is Cocoon?

Cocoon isolates your development projects to prevent:
- Malicious file writes outside your project folder
- Unauthorized network calls
- Accidental exposure of API keys and secrets

Students can install once and run projects safely without configuration.

## Installation

```bash
git clone git@github.com:fr4nsyz/Cocoon.git
cd Cocoon
go build -o cocoon ./cmd/cocoon
```

Then run `./cocoon --help` to verify.

## Requirements

- **Docker** (optional) - For container sandboxing mode
- **Python, Node.js, or Go** - For wrapper fallback mode

If Docker is not available, Cocoon automatically falls back to lightweight wrapper mode.

## Quick Start

### Basic Usage

```bash
# Python project
cocoon python main.py

# Node.js project
cocoon npm start

# Go project
cocoon go run main.go

# With custom command
cocoon python -m flask run
```

### Network Modes

**Default (blocked):** All network connections blocked
```bash
cocoon go run main.go
```

**Full (for installing dependencies):**
```bash
cocoon --network=full npm install
cocoon --network=full pip install -r requirements.txt
cocoon --network=full go mod download
```

**Local (localhost only):** Outbound blocked, but localhost connections allowed. Requires `--expose-ports`:
```bash
cocoon --network=local --expose-ports=8080 go run main.go
```

**Whitelist:** Currently same as blocked.

### Exposing Ports

Ports are auto-detected from your project files. Manual override:

```bash
cocoon --expose-ports=8080 go run main.go
```

### Wrapper Mode (No Docker)

Force lightweight mode without Docker:

```bash
cocoon --no-container python main.py
```

### Clean Environment

Strip all environment variables before running:

```bash
cocoon --clean-env python main.py
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-d, --project-dir` | `.` | Project directory to sandbox |
| `--network` | `none` | Network mode: none, local, whitelist, full |
| `--expose-ports` | `auto` | Ports to expose (auto or comma-separated) |
| `-v, --verbose` | false | Show sandbox activity in real-time |
| `--no-container` | false | Force wrapper mode (no Docker) |
| `--clean-env` | false | Strip all env vars before running |

## Common Workflows

### Run a Python Flask App

```bash
cd my-flask-app
cocoon python app.py
```

### Run a Go HTTP Server

```bash
cd my-go-app
cocoon --expose-ports=8080 go run main.go
```

The server will be accessible at localhost:8080. To allow localhost connections (e.g., for database proxies), use `--network=local`:

```bash
cocoon --network=local --expose-ports=8080 go run main.go
```

### Install Dependencies

```bash
# Go
cocoon --network=full go mod download

# Node.js
cocoon --network=full npm install

# Python
cocoon --network=full pip install -r requirements.txt
```

### Run Tests

```bash
cocoon go test ./...
# or
cocoon pytest
# or
cocoon npm test
```

## How It Works

### Container Mode (with Docker)

When Docker is available, Cocoon runs your project inside an isolated container with:

- **Filesystem:** Read-only except for project folder
- **Network:** Blocked by default (use `--network=full` for installs)
- **Capabilities:** All dropped (`--cap-drop=ALL`)
- **Memory:** Limited to 512MB
- **Processes:** Limited to 100

> **Note about dependency installation:** When using `--network=full` for installing dependencies (e.g., `npm install`, `pip install`, `go mod download`), the container mounts your project directory as read-write. This means that while the installation process runs inside the container, any created files (like `node_modules/`, installed packages, or downloaded modules) will appear directly in your host project directory. Postinstall scripts execute inside the container but can modify files in your host project folder. For maximum isolation during installation, consider using a temporary directory or reviewing packages before installation.

### Wrapper Mode (fallback)

When Docker is unavailable, Cocoon runs your project directly with basic process isolation.

## Troubleshooting

### "Docker not available"

Install Docker or use `--no-container` for lightweight mode.

### "Unsupported project type"

Ensure your project has the right files:
- Python: `requirements.txt`, `setup.py`, or `pyproject.toml`
- Node.js: `package.json`
- Go: `go.mod`

### "Runtime not found"

Make sure Python, Node.js, or Go is installed on your system.

### Port already in use

Specify a different port:
```bash
cocoon --expose-ports=8081 go run main.go
```

## Security

Cocoon is designed for developers working with:
- Untrusted third-party packages and dependencies
- Runtime code that might contain vulnerabilities
- Learning projects that might accidentally expose secrets

### What Cocoon Protects Against:
- Malicious code writing files outside your project directory
- Unauthorized network connections when using default (blocked) network mode
- Accidental leakage of environment variables and secrets
- Resource exhaustion attacks (memory/process limits)
- Privilege escalation attempts (dropped capabilities, non-root user)
- Runtime exploitation of vulnerable dependencies

### Limitations to Understand:
- **During dependency installation with `--network=full`**: Container gets direct host network access and project directory is mounted read-write, allowing postinstall scripts to modify host project files
- **Not a container escape guarantee**: While Cocoon applies strong isolation practices, determined attackers might still find container escape vulnerabilities (as with any container-based solution)
- **Not for production workloads**: Designed for development/learning scenarios

It is **not** a replacement for production-grade sandboxing. For production use, consider gVisor, Firejail, or proper container orchestration.

## Future Features

- [ ] Network whitelist (allow npm/pypi registries)
- [ ] GUI dashboard showing blocked actions
- [ ] Windows support
- [ ] Auto-detect hardcoded secrets in source files

## License

MIT
