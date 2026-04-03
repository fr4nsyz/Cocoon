package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cocoon/cocoon/pkg/detection"
)

var ErrProjectNotFound = errors.New("project directory not found")

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
		if os.IsNotExist(err) {
			return "", ErrProjectNotFound
		}
		return "", fmt.Errorf("project directory error: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("project path is not a directory: %s", absDir)
	}

	return absDir, nil
}
