package service

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/service/traffic"
)

// Service represents a Cloud Run service.
type Service struct {
	Name                  string                         `json:"name"`
	Description           string                         `json:"description,omitempty"`
	URI                   string                         `json:"uri"`
	CreateTime            time.Time                      `json:"createTime"`
	UpdateTime            time.Time                      `json:"updateTime"`
	DeleteTime            time.Time                      `json:"deleteTime"`
	ExpireTime            time.Time                      `json:"expireTime"`
	Creator               string                         `json:"creator,omitempty"`
	LastModifier          string                         `json:"lastModifier,omitempty"`
	Reconciling           bool                           `json:"reconciling"`
	TrafficStatuses       []*traffic.TrafficTargetStatus `json:"trafficStatuses,omitempty"`
	LatestReadyRevision   string                         `json:"latestReadyRevision,omitempty"`
	LatestCreatedRevision string                         `json:"latestCreatedRevision,omitempty"`
	TerminalCondition     *condition.Condition           `json:"terminalCondition,omitempty"`
	Conditions            []*condition.Condition         `json:"conditions,omitempty"`
	Etag                  string                         `json:"etag,omitempty"`
}
