package secret

import "github.com/JulienBreux/run-cli/internal/run/model/common/keytopath"

// SecretKeySelector selects a key of a Secret.
type SecretKeySelector struct {
	Secret string `json:"secret"`
	Key    string `json:"key"`
}

// SecretSource represents a secret that is mounted as a volume.
type SecretSource struct {
	Secret      string                 `json:"secret"`
	Items       []*keytopath.KeyToPath `json:"items"`
	DefaultMode int32                  `json:"defaultMode"`
}
