package container

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}
