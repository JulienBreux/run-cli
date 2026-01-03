package networking

// Networking represents the networking settings of a service.
type Networking struct {
	Ingress            string     `json:"ingress"`
	DefaultUriDisabled bool       `json:"defaultUriDisabled"`
	IapEnabled         bool       `json:"iapEnabled"`
	VpcAccess          *VpcAccess `json:"vpcAccess,omitempty"`
}

// VpcAccess represents the VPC Access settings.
type VpcAccess struct {
	Connector string `json:"connector,omitempty"`
	Egress    string `json:"egress,omitempty"`
}
