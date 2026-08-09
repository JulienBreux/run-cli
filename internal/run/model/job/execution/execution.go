/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package execution

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
)

// Execution represents a Cloud Run job execution.
type Execution struct {
	Name              string                 `json:"name"`
	Job               string                 `json:"job"`
	CreateTime        time.Time              `json:"createTime"`
	StartTime         time.Time              `json:"startTime"`
	CompletionTime    time.Time              `json:"completionTime"`
	DeleteTime        time.Time              `json:"deleteTime"`
	ExpireTime        time.Time              `json:"expireTime"`
	TaskCount         int32                  `json:"taskCount"`
	SucceededCount    int32                  `json:"succeededCount"`
	FailedCount       int32                  `json:"failedCount"`
	RunningCount      int32                  `json:"runningCount"`
	CancelledCount    int32                  `json:"cancelledCount"`
	RetriedCount      int32                  `json:"retriedCount"`
	LogURI            string                 `json:"logUri"`
	Region            string                 `json:"region"`
	Conditions        []*condition.Condition `json:"conditions"`
	TerminalCondition *condition.Condition   `json:"terminalCondition"`
}
