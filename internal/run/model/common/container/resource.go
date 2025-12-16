package container

// Resources represents the computational resources of a container.
type Resources struct {
	Limits          map[string]string `json:"limits"`
	CPUIdle         bool              `json:"cpuIdle"`
	StartupCPUBoost bool              `json:"startupCpuBoost"`
}
