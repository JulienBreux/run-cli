package execution

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
)

// Execution represents a Cloud Run job execution.
type Execution struct {
	Name              string               `json:"name"`
	Job               string               `json:"job"`
	CreateTime        time.Time            `json:"createTime"`
	StartTime         time.Time            `json:"startTime"`
	CompletionTime    time.Time            `json:"completionTime"`
	DeleteTime        time.Time            `json:"deleteTime"`
	ExpireTime        time.Time            `json:"expireTime"`
	TaskCount         int32                `json:"taskCount"`
	SucceededCount    int32                `json:"succeededCount"`
	FailedCount       int32                `json:"failedCount"`
	RunningCount      int32                `json:"runningCount"`
	CancelledCount    int32                `json:"cancelledCount"`
	RetriedCount      int32                `json:"retriedCount"`
	LogURI            string               `json:"logUri"`
	Region            string               `json:"region"`
	Conditions        []*condition.Condition `json:"conditions"`
	TerminalCondition *condition.Condition `json:"terminalCondition"`
}
