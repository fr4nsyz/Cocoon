# Cocoon

Zero-config sandbox for running student projects safely.

Your projects run in their own space—nothing can touch your machine or network.

## What is Cocoon?

Cocoon isolates your development projects to prevent:
- Malicious file writes outside your project folder
- Unauthorized network calls
- Accidental exposure of API keys and secrets

Students can install once and run projects safely without configuration.

## Installation

### Quick Install (macOS / Linux)

```bash
curl -sSL https://get.cocoon.dev | sh
```

### Via Go

```bash
go install github.com/cocoon/cocoon@latest
```

### Manual

1. Download the latest release for your platform from [GitHub Releases](https://github.com/cocoon/cocoon/releases)
2. Extract and add to your PATH
3. Run `cocoon --help` to verify

## Requirements

- **Docker** (optional) - For container sandboxing mode
- **Python or Node.js** - For wrapper fallback mode

If Docker is not available, Cocoon automatically falls back to lightweight wrapper mode.

## Quick Start

### Basic Usage

```bash
# Python project
cocoon python main.py

# Node.js project
cocoon npm start

# With custom command
cocoon python -m flask run
```

### Network Modes

**Default (blocked):** All network connections blocked
```bash
cocoon npm start
```

**Full (for installing dependencies):**
```bash
cocoon --network=full npm install
cocoon --network=full pip install -r requirements.txt
```

**Whitelist (coming soon):** Currently same as blocked.

### Exposing Ports

Ports are auto-detected from your project files. Manual override:

```bash
cocoon --expose-ports=3000,8080 npm start
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
| `--network` | `none` | Network mode: none, whitelist, full |
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

### Run a Node.js React App

```bash
cd my-react-app
cocoon npm start
```

The React dev server will be accessible at localhost:3000.

### Install Dependencies

```bash
# Node.js
cocoon --network=full npm install

# Python
cocoon --network=full pip install -r requirements.txt
```

### Run Tests

```bash
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

### Wrapper Mode (fallback)

When Docker is unavailable, Cocoon runs your project directly with basic process isolation.

## Troubleshooting

### "Docker not available"

Install Docker or use `--no-container` for lightweight mode.

### "Unsupported project type"

Ensure your project has the right files:
- Python: `requirements.txt`, `setup.py`, or `pyproject.toml`
- Node.js: `package.json`

### "Runtime not found"

Make sure Python or Node.js is installed on your system.

### Port already in use

Specify a different port:
```bash
cocoon --expose-ports=3001 npm start
```

## Security

Cocoon is designed for student developers working with:
- Tutorial code from the internet
- Untrusted npm/pip packages
- Learning projects that might accidentally expose secrets

It is **not** a replacement for production-grade sandboxing. For production use, consider gVisor, Firejail, or proper container orchestration.

## Future Features

- [ ] Network whitelist (allow npm/pypi registries)
- [ ] GUI dashboard showing blocked actions
- [ ] Windows support
- [ ] Auto-detect hardcoded secrets in source files

## License

MIT
