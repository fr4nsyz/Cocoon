package main

import (
	"os"
	"os/exec"

	"github.com/cocoon/cocoon/internal/config"
	"github.com/cocoon/cocoon/internal/logging"
	"github.com/cocoon/cocoon/pkg/sandbox"
)

func init() {
	sandbox.RegisterWrapper("python", runPythonWrapper)
	sandbox.RegisterWrapper("node", runNodeWrapper)
}

func runPythonWrapper(cfg *config.Config, logger *logging.Logger) error {
	logger.Info("Running Python wrapper")

	pythonCmd := findPython()
	if pythonCmd == "" {
		return sandbox.ErrRuntimeNotFound
	}

	args := cfg.Command
	if len(args) > 0 && (args[0] == "python" || args[0] == "python3" || args[0] == "py") {
		args = args[1:]
	}

	cmd := exec.Command(pythonCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = cfg.ProjectDir

	if cfg.CleanEnv {
		cmd.Env = []string{"PATH=" + os.Getenv("PATH")}
	}

	return cmd.Run()
}

func runNodeWrapper(cfg *config.Config, logger *logging.Logger) error {
	logger.Info("Running Node wrapper")

	nodeCmd := findNode()
	if nodeCmd == "" {
		return sandbox.ErrRuntimeNotFound
	}

	args := cfg.Command
	if len(args) > 0 && (args[0] == "node" || args[0] == "nodejs") {
		args = args[1:]
	}

	cmd := exec.Command(nodeCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = cfg.ProjectDir

	if cfg.CleanEnv {
		cmd.Env = []string{"PATH=" + os.Getenv("PATH")}
	}

	return cmd.Run()
}

func findPython() string {
	commands := []string{"python3", "python", "py"}
	for _, cmd := range commands {
		if path, err := exec.LookPath(cmd); err == nil {
			return path
		}
	}
	return ""
}

func findNode() string {
	commands := []string{"node", "nodejs"}
	for _, cmd := range commands {
		if path, err := exec.LookPath(cmd); err == nil {
			return path
		}
	}
	return ""
}
