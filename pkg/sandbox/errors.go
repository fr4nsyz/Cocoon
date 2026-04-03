package sandbox

import "fmt"

var (
	ErrDockerNotAvailable     = fmt.Errorf("docker not available")
	ErrProjectNotFound        = fmt.Errorf("project directory not found")
	ErrProjectTypeUnsupported = fmt.Errorf("unsupported project type")
	ErrRuntimeNotFound        = fmt.Errorf("runtime not found")
	ErrContainerFailed        = fmt.Errorf("container execution failed")
	ErrWrapperFailed          = fmt.Errorf("wrapper execution failed")
)

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration error: %s - %s", e.Field, e.Message)
}

type RuntimeError struct {
	Runtime string
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s error: %s", e.Runtime, e.Message)
}
