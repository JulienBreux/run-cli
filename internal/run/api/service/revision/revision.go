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

package revision

import (
	"context"
	"strings"

	"cloud.google.com/go/run/apiv2/runpb"
	model_container "github.com/JulienBreux/run-cli/internal/run/model/common/container"
	model_resources "github.com/JulienBreux/run-cli/internal/run/model/common/resources"
	model "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
)

// List returns a list of revisions for the given service.
func List(project, region, service string) ([]model.Revision, error) {
	ctx := context.Background()
	pbRevisions, err := apiClient.ListRevisions(ctx, project, region, service)
	if err != nil {
		return nil, err
	}

	var revisions []model.Revision
	for _, resp := range pbRevisions {
		revisions = append(revisions, mapRevision(resp, service))
	}

	return revisions, nil
}

func mapRevision(resp *runpb.Revision, service string) model.Revision {
	nameParts := strings.Split(resp.Name, "/")
	name := nameParts[len(nameParts)-1]

	var cpuIdle, startupCpuBoost bool
	if len(resp.Containers) > 0 && resp.Containers[0].Resources != nil {
		cpuIdle = resp.Containers[0].Resources.CpuIdle
		startupCpuBoost = resp.Containers[0].Resources.StartupCpuBoost
	}

	var containers []*model_container.Container
	for _, c := range resp.Containers {
		var ports []*model_container.Port
		for _, p := range c.Ports {
			ports = append(ports, &model_container.Port{
				Name:          p.Name,
				ContainerPort: p.ContainerPort,
			})
		}

		var resources *model_resources.Resources
		if c.Resources != nil {
			resources = &model_resources.Resources{
				Limits:          c.Resources.Limits,
				CPUIdle:         c.Resources.CpuIdle,
				StartupCPUBoost: c.Resources.StartupCpuBoost,
			}
		}

		containers = append(containers, &model_container.Container{
			Name:      c.Name,
			Image:     c.Image,
			Command:   c.Command,
			Args:      c.Args,
			Ports:     ports,
			Resources: resources,
		})
	}

	var accelerator string
	if resp.NodeSelector != nil {
		accelerator = resp.NodeSelector.Accelerator
	}

	return model.Revision{
		Name:                          name,
		CreateTime:                    resp.CreateTime.AsTime(),
		UpdateTime:                    resp.UpdateTime.AsTime(),
		Service:                       service,
		Containers:                    containers,
		ExecutionEnvironment:          resp.ExecutionEnvironment.String(),
		MaxInstanceRequestConcurrency: resp.MaxInstanceRequestConcurrency,
		Timeout:                       resp.Timeout.AsDuration(),
		CpuIdle:                       cpuIdle,
		StartupCpuBoost:               startupCpuBoost,
		Accelerator:                   accelerator,
	}
}
