package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/saferun/saferun/pkg/detection"
)

type Config struct {
	ProjectDir   string
	ProjectType  detection.ProjectType
	Command      []string
	NetworkMode  string
	ExposePorts  string
	ExposedPorts []int
	CleanEnv     bool
	Verbose      bool
	NoContainer  bool
}

func ResolveProjectDir(dir string) (string, error) {
	if dir == "" {
		dir = "."
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project directory: %w", err)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		return "", fmt.Errorf("project directory does not exist: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("project path is not a directory: %s", absDir)
	}

	return absDir, nil
}
