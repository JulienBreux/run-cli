package condition

import "time"

// Condition represents a resource's status.
type Condition struct {
	Type                string    `json:"type,omitempty"`
	State               string    `json:"state,omitempty"`
	Message             string    `json:"message,omitempty"`
	LastTransitionTime  time.Time `json:"lastTransitionTime"`
	Severity            string    `json:"severity,omitempty"`
	Reason              string    `json:"reason,omitempty"`
	RevisionGeneration  int64     `json:"revisionGeneration,omitempty"`
	ObservedGeneration  int64     `json:"observedGeneration,omitempty"`
	DomainMappingReason string    `json:"domainMappingReason,omitempty"`
}
