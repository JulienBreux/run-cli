package project

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	model "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"google.golang.org/api/iterator"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// List returns a list of projects for the current user.
func List() ([]model.Project, error) {
	ctx := context.Background()

	// Explicitly find default credentials
	creds, err := google.FindDefaultCredentials(ctx, resourcemanager.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w. Tip: Try running 'gcloud auth application-default login' to authenticate the Go client", err)
	}

	c, err := resourcemanager.NewProjectsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = c.Close()
	}()

	req := &resourcemanagerpb.SearchProjectsRequest{
		// Query: "", // Empty query lists all projects
	}

	var projects []model.Project
	it := c.SearchProjects(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
				return nil, fmt.Errorf("authentication failed: %w. Tip: Ensure your 'gcloud auth application-default login' is valid and has permissions", err)
			}
			return nil, err
		}

		// Parse Project Number from Name "projects/123456"
		parts := strings.Split(resp.Name, "/")
		numberStr := parts[len(parts)-1]
		number, _ := strconv.Atoi(numberStr)

		projects = append(projects, model.Project{
			Name:   resp.ProjectId,
			Number: number,
		})
	}

	return projects, nil
}