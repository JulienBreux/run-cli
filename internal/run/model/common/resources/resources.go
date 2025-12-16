package resources

// Resources represents the computational resources of a container.
type Resources struct {
	Limits          map[string]string `json:"limits,omitempty"`
	CPUIdle         bool              `json:"cpuIdle,omitempty"`
	StartupCPUBoost bool              `json:"startupCpuBoost,omitempty"`
}