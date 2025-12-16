package traffic

// TrafficTargetStatus represents the actual traffic allocated to a revision.
type TrafficTargetStatus struct {
	Type           string `json:"type,omitempty"`
	Revision       string `json:"revision,omitempty"`
	Percent        int32  `json:"percent,omitempty"`
	LatestRevision bool   `json:"latestRevision,omitempty"`
	URI            string `json:"uri,omitempty"`
}
