package container

// Port represents a network port in a single container.
type Port struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int32  `json:"containerPort,omitempty"`
}
