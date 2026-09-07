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

package instance

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	"github.com/JulienBreux/run-cli/internal/run/model/common/volume"
	model "github.com/JulienBreux/run-cli/internal/run/model/instance"
)

var apiClient Client = &GCPClient{}

// List returns a list of instances for the given project and region.
// If region is api_region.ALL, it lists instances from all supported Cloud Run regions.
func List(project, region string) ([]model.Instance, error) {
	if region == api_region.ALL {
		return listAllRegions(project)
	}

	ctx := context.Background()
	pbInstances, err := apiClient.ListInstances(ctx, project, region)
	if err != nil {
		return nil, err
	}

	var instances []model.Instance
	for _, resp := range pbInstances {
		instances = append(instances, mapInstance(resp, region, project))
	}

	return instances, nil
}

// Get returns an instance given project, region, and short name.
func Get(project, region, name string) (*model.Instance, error) {
	ctx := context.Background()
	fullName := formatInstanceName(project, region, name)
	resp, err := apiClient.GetInstance(ctx, fullName)
	if err != nil {
		return nil, err
	}
	inst := mapInstance(resp, region, project)
	return &inst, nil
}

// Start starts an instance.
func Start(project, region, name string) error {
	ctx := context.Background()
	fullName := formatInstanceName(project, region, name)
	_, err := apiClient.StartInstance(ctx, fullName)
	return err
}

// Stop stops an instance.
func Stop(project, region, name string) error {
	ctx := context.Background()
	fullName := formatInstanceName(project, region, name)
	_, err := apiClient.StopInstance(ctx, fullName)
	return err
}

// Delete deletes an instance.
func Delete(project, region, name string) error {
	ctx := context.Background()
	fullName := formatInstanceName(project, region, name)
	_, err := apiClient.DeleteInstance(ctx, fullName)
	return err
}

func formatInstanceName(project, region, name string) string {
	if strings.HasPrefix(name, "projects/") {
		return name
	}
	return fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, region, name)
}

func mapInstance(resp *runpb.Instance, region, project string) model.Instance {
	inst := model.Instance{
		Name:                       resp.Name,
		Description:                resp.Description,
		UID:                        resp.Uid,
		Generation:                 resp.Generation,
		Labels:                     resp.Labels,
		Annotations:                resp.Annotations,
		Creator:                    resp.Creator,
		LastModifier:               resp.LastModifier,
		Client:                     resp.Client,
		ClientVersion:              resp.ClientVersion,
		LaunchStage:                resp.LaunchStage.String(),
		ServiceAccount:             resp.ServiceAccount,
		EncryptionKey:              resp.EncryptionKey,
		ObservedGeneration:         resp.ObservedGeneration,
		LogURI:                     resp.LogUri,
		SatisfiesPZS:               resp.SatisfiesPzs,
		URLs:                       resp.Urls,
		Reconciling:                resp.Reconciling,
		ETag:                       resp.Etag,
		GpuZonalRedundancyDisabled: resp.GetGpuZonalRedundancyDisabled(),
		InvokerIamDisabled:         resp.InvokerIamDisabled,
		IapEnabled:                 resp.IapEnabled,
		Region:                     region,
		Project:                    project,
	}

	if resp.CreateTime != nil {
		inst.CreateTime = resp.CreateTime.AsTime()
	}
	if resp.UpdateTime != nil {
		inst.UpdateTime = resp.UpdateTime.AsTime()
	}
	if resp.DeleteTime != nil {
		inst.DeleteTime = resp.DeleteTime.AsTime()
	}
	if resp.ExpireTime != nil {
		inst.ExpireTime = resp.ExpireTime.AsTime()
	}

	if resp.TerminalCondition != nil {
		inst.TerminalCondition = &condition.Condition{
			Type:               resp.TerminalCondition.Type,
			State:              resp.TerminalCondition.State.String(),
			Message:            resp.TerminalCondition.Message,
			LastTransitionTime: resp.TerminalCondition.LastTransitionTime.AsTime(),
			Severity:           resp.TerminalCondition.Severity.String(),
			Reason:             resp.TerminalCondition.GetReason().String(),
		}
	}

	for _, c := range resp.Conditions {
		inst.Conditions = append(inst.Conditions, &condition.Condition{
			Type:               c.Type,
			State:              c.State.String(),
			Message:            c.Message,
			LastTransitionTime: c.LastTransitionTime.AsTime(),
			Severity:           c.Severity.String(),
			Reason:             c.GetReason().String(),
		})
	}

	for _, cs := range resp.ContainerStatuses {
		inst.ContainerStatuses = append(inst.ContainerStatuses, &model.ContainerStatus{
			Name:        cs.Name,
			ImageDigest: cs.ImageDigest,
		})
	}

	for _, c := range resp.Containers {
		inst.Containers = append(inst.Containers, &container.Container{
			Name:  c.Name,
			Image: c.Image,
		})
	}

	for _, v := range resp.Volumes {
		inst.Volumes = append(inst.Volumes, &volume.Volume{
			Name: v.Name,
		})
	}

	return inst
}

func listAllRegions(project string) ([]model.Instance, error) {
	var (
		mu        sync.Mutex
		instances []model.Instance
		wg        sync.WaitGroup
	)

	for _, region := range api_region.List() {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			if list, err := List(project, r); err == nil {
				mu.Lock()
				instances = append(instances, list...)
				mu.Unlock()
			}
		}(region)
	}

	wg.Wait()
	return instances, nil
}
