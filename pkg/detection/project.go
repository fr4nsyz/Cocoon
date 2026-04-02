package detection

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProjectType string

const (
	UnknownProject ProjectType = "unknown"
	PythonProject  ProjectType = "python"
	NodeProject    ProjectType = "node"
	GoProject      ProjectType = "go"
	RubyProject    ProjectType = "ruby"
)

func DetectProjectType(dir string) ProjectType {
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		return NodeProject
	}
	if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
		return PythonProject
	}
	if _, err := os.Stat(filepath.Join(dir, "setup.py")); err == nil {
		return PythonProject
	}
	if _, err := os.Stat(filepath.Join(dir, "pyproject.toml")); err == nil {
		return PythonProject
	}
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return GoProject
	}
	if _, err := os.Stat(filepath.Join(dir, "Gemfile")); err == nil {
		return RubyProject
	}
	return UnknownProject
}

func DetectExposedPorts(projectDir string, exposePorts string) []int {
	if exposePorts == "" || exposePorts == "auto" {
		return detectPortsAuto(projectDir)
	}
	return parsePortString(exposePorts)
}

func detectPortsAuto(dir string) []int {
	ports := []int{}

	if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		var pkg struct {
			Scripts map[string]string `json:"scripts"`
			Config  map[string]any    `json:"config"`
		}
		if json.Unmarshal(data, &pkg) == nil {
			for script, cmd := range pkg.Scripts {
				if strings.Contains(cmd, "react-scripts start") ||
					strings.Contains(cmd, "vite") ||
					strings.Contains(cmd, "next start") ||
					strings.Contains(cmd, "vue-cli-service serve") {
					ports = appendUnique(ports, 3000)
				}
				if strings.Contains(cmd, "angular-cli") || strings.Contains(cmd, "ng serve") {
					ports = appendUnique(ports, 4200)
				}
				if script == "start" && !strings.Contains(cmd, "react") {
					ports = appendUnique(ports, 3000)
				}
			}
			if port, ok := pkg.Config["port"]; ok {
				switch v := port.(type) {
				case float64:
					ports = appendUnique(ports, int(v))
				case string:
					if p, err := strconv.Atoi(v); err == nil {
						ports = appendUnique(ports, p)
					}
				}
			}
		}
	}

	if data, err := os.ReadFile(filepath.Join(dir, ".env")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "PORT=") {
				portStr := strings.TrimPrefix(line, "PORT=")
				portStr = strings.TrimSpace(portStr)
				if p, err := strconv.Atoi(portStr); err == nil {
					ports = appendUnique(ports, p)
				}
			}
		}
	}

	if data, err := os.ReadFile(filepath.Join(dir, "vite.config.js")); err == nil {
		content := string(data)
		if strings.Contains(content, "port:") {
			ports = appendUnique(ports, 5173)
		}
	}

	if data, err := os.ReadFile(filepath.Join(dir, "webpack.config.js")); err == nil {
		content := string(data)
		if strings.Contains(content, "port:") {
			ports = appendUnique(ports, 8080)
		}
	}

	return ports
}

func parsePortString(s string) []int {
	ports := []int{}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if p, err := strconv.Atoi(part); err == nil && p > 0 && p < 65536 {
			ports = append(ports, p)
		}
	}
	return ports
}

func appendUnique(slice []int, item int) []int {
	for _, v := range slice {
		if v == item {
			return slice
		}
	}
	return append(slice, item)
}
