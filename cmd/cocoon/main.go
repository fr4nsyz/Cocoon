package main

import (
	"os"

	"github.com/cocoon/cocoon/internal/config"
	"github.com/cocoon/cocoon/internal/logging"
	"github.com/cocoon/cocoon/pkg/detection"
	"github.com/cocoon/cocoon/pkg/sandbox"
	"github.com/spf13/cobra"
)

var (
	cfg         *config.Config
	logger      *logging.Logger
	projectDir  string
	networkMode string
	exposePorts string
	verbose     bool
	noContainer bool
	cleanEnv    bool
)

var rootCmd = &cobra.Command{
	Use:   "cocoon [command]",
	Short: "Zero-config sandbox for running student projects safely",
	Long: `Cocoon - Run your projects in a secure sandbox automatically

Your projects run safely in their own space—nothing can touch your 
machine or network.

Examples:
  cocoon python main.py
  cocoon npm start
  cocoon --expose-ports=3000,8080 npm start
  cocoon --network=whitelist pip install -r requirements.txt`,
	RunE: runSandbox,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.Flags().StringVarP(&projectDir, "project-dir", "d", ".", "Project directory to sandbox")
	rootCmd.Flags().StringVar(&networkMode, "network", "none", "Network mode: none, whitelist, full")
	rootCmd.Flags().StringVar(&exposePorts, "expose-ports", "auto", "Ports to expose (auto-detect or comma-separated)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show blocked actions in real-time")
	rootCmd.Flags().BoolVar(&noContainer, "no-container", false, "Force wrapper mode (no Docker)")
	rootCmd.Flags().BoolVar(&cleanEnv, "clean-env", false, "Strip all env vars before running")

	rootCmd.Flags().BoolP("help", "h", false, "Show help")
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		c.Printf(c.Long + "\n\n")
		c.Printf("Usage:\n  %s\n\n", c.Use)
		c.Printf("Flags:\n%s\n", c.Flags().FlagUsages())
	})
}

func runSandbox(cmd *cobra.Command, args []string) error {
	logger = logging.New(verbose)

	absProjectDir, err := config.ResolveProjectDir(projectDir)
	if err != nil {
		return err
	}

	projectType := detection.DetectProjectType(absProjectDir)
	logger.Info("Detected project type: %s", projectType)

	cfg = &config.Config{
		ProjectDir:  absProjectDir,
		ProjectType: projectType,
		Command:     args,
		NetworkMode: networkMode,
		ExposePorts: exposePorts,
		CleanEnv:    cleanEnv,
		Verbose:     verbose,
		NoContainer: noContainer,
	}

	exposedPorts := detection.DetectExposedPorts(absProjectDir, exposePorts)
	cfg.ExposedPorts = exposedPorts

	if len(exposedPorts) > 0 {
		logger.Info("Exposing ports: %v", exposedPorts)
	}

	runner := sandbox.NewRunner(cfg, logger)
	return runner.Run()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		if err == sandbox.ErrSandboxNotAvailable {
			logger.Error("Docker not available. Install Docker or use --no-container for lightweight sandboxing.")
		}
		os.Exit(1)
	}
}
