package scaling

// Scaling represents the scaling configuration of a Cloud Run worker pool.
type Scaling struct {
	ManualInstanceCount int32 `json:"manualInstanceCount,omitempty"`
}
