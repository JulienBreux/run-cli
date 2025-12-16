package container

// ExecAction describes a "run in container" action.
type ExecAction struct {
	Command []string `json:"command,omitempty"`
}
