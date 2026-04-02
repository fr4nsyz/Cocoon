package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cocoon/cocoon/internal/config"
	"github.com/cocoon/cocoon/internal/logging"
	"github.com/cocoon/cocoon/pkg/isolation"
)

var ErrSandboxNotAvailable = fmt.Errorf("sandbox not available")

type Runner struct {
	cfg    *config.Config
	logger *logging.Logger
}

func NewRunner(cfg *config.Config, logger *logging.Logger) *Runner {
	return &Runner{cfg: cfg, logger: logger}
}

func (r *Runner) Run() error {
	r.logger.Info("Starting sandbox for: %s", r.cfg.ProjectDir)

	if r.cfg.CleanEnv {
		r.logger.Info("Cleaning environment variables")
		os.Clearenv()
	} else {
		secrets := isolation.ScanEnvForSecrets(os.Environ())
		if len(secrets) > 0 {
			r.logger.Warn("Detected %d sensitive environment variables", len(secrets))
			if r.cfg.Verbose {
				for key := range secrets {
					r.logger.Warn("  - %s", key)
				}
			}
		}
	}

	dockerAvailable := IsDockerAvailable()
	if dockerAvailable && !r.cfg.NoContainer {
		return r.runContainer()
	}

	if !dockerAvailable && !r.cfg.NoContainer {
		r.logger.Warn("Docker not available, falling back to wrapper mode")
	}

	return r.runWrapper()
}

func IsDockerAvailable() bool {
	cmd := exec.Command("docker", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (r *Runner) runContainer() error {
	r.logger.Info("Running in container mode")

	dockerArgs := r.buildDockerArgs()

	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = r.cfg.ProjectDir

	r.logger.Debug("Docker command: docker %s", strings.Join(dockerArgs, " "))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("container execution failed: %w", err)
	}

	return nil
}

func (r *Runner) buildDockerArgs() []string {
	args := []string{
		"run", "--rm",
		"--read-only",
		"--cap-drop=ALL",
		"--pids-limit=100",
		"--memory=512m",
		"--security-opt=no-new-privileges",
	}

	switch r.cfg.NetworkMode {
	case "none":
		args = append(args, "--network=none")
		r.logger.Info("Network mode: blocked (no external connections)")
	case "whitelist":
		args = append(args, "--network=none")
		r.logger.Warn("Whitelist mode is not yet fully implemented. Network is blocked.")
		r.logger.Info("For package installs (npm install, pip install), use --network=full temporarily")
		r.logger.Info("Allowed registries (for future use): %v", isolation.GetAllowedHosts())
	case "full":
		args = append(args, "--network=host")
		r.logger.Warn("Network mode: full (all connections allowed - not recommended for running apps)")
	default:
		args = append(args, "--network=none")
		r.logger.Info("Network mode: blocked (no external connections)")
	}

	for _, port := range r.cfg.ExposedPorts {
		args = append(args, "-p", fmt.Sprintf("%d:%d", port, port))
	}

	args = append(args, "-v", fmt.Sprintf("%s:/sandbox:rw", r.cfg.ProjectDir))
	args = append(args, "-w", "/sandbox")
	args = append(args, "-e", "HOME=/sandbox")
	args = append(args, "--user=1000")

	args = append(args, getBaseImage(string(r.cfg.ProjectType)))

	args = append(args, r.cfg.Command...)

	return args
}

func getBaseImage(projectType string) string {
	images := map[string]string{
		"python": "python:3.11-slim",
		"node":   "node:20-slim",
		"go":     "golang:1.21-slim",
		"ruby":   "ruby:3.2-slim",
	}
	if img, ok := images[string(projectType)]; ok {
		return img
	}
	return "debian:bookworm-slim"
}

func (r *Runner) runWrapper() error {
	r.logger.Info("Running in wrapper mode")

	runner := GetWrapperRunner(string(r.cfg.ProjectType))
	if runner == nil {
		return fmt.Errorf("no wrapper available for project type: %s", r.cfg.ProjectType)
	}

	return runner(r.cfg, r.logger)
}

type wrapperFunc func(cfg *config.Config, logger *logging.Logger) error

var wrappers = map[string]wrapperFunc{}

func RegisterWrapper(projectType string, fn wrapperFunc) {
	wrappers[projectType] = fn
}

func GetWrapperRunner(projectType string) wrapperFunc {
	return wrappers[projectType]
}
