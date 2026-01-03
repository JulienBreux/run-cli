package security

// Security represents the security settings of a service.
type Security struct {
	InvokerIAMDisabled      bool   `json:"invokerIamDisabled"`
	ServiceAccount          string `json:"serviceAccount,omitempty"`
	EncryptionKey           string `json:"encryptionKey,omitempty"`
	BinaryAuthorization     string `json:"binaryAuthorization,omitempty"`
	BreakglassJustification string `json:"breakglassJustification,omitempty"`
}
