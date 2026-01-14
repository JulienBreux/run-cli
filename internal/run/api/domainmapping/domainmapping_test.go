package domainmapping

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/run/v1"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	ListDomainMappingsFunc func(ctx context.Context, project, region string) ([]*run.DomainMapping, error)
}

func (m *MockClient) ListDomainMappings(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
	if m.ListDomainMappingsFunc != nil {
		return m.ListDomainMappingsFunc(ctx, project, region)
	}
	return nil, nil
}

func TestMapDomainMapping(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	resp := &run.DomainMapping{
		Metadata: &run.ObjectMeta{
			Name:              "example.com",
			CreationTimestamp: now,
			Annotations: map[string]string{
				"serving.knative.dev/creator": "user@example.com",
			},
		},
		Spec: &run.DomainMappingSpec{
			RouteName: "my-service",
		},
		Status: &run.DomainMappingStatus{
			ResourceRecords: []*run.ResourceRecord{
				{
					Type:   "CNAME",
					Name:   "www",
					Rrdata: "ghs.googlehosted.com.",
				},
			},
			Conditions: []*run.GoogleCloudRunV1Condition{
				{
					Type:    "Ready",
					Status:  "True",
					Message: "Ready",
				},
			},
		},
	}

	result := mapDomainMapping(resp, "my-project", "us-central1")

	assert.Equal(t, "example.com", result.Name)
	assert.Equal(t, "my-service", result.RouteName)
	assert.Equal(t, "my-project", result.Project)
	assert.Equal(t, "user@example.com", result.Creator)
	assert.Equal(t, "us-central1", result.Region)
	assert.Len(t, result.Records, 1)
	assert.Equal(t, "CNAME", result.Records[0].Type)
	assert.Equal(t, "www", result.Records[0].Name)
	assert.Equal(t, "ghs.googlehosted.com.", result.Records[0].RRData)
	assert.Len(t, result.Conditions, 1)
	assert.Equal(t, "Ready", result.Conditions[0].Type)
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListDomainMappingsFunc = func(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
		return []*run.DomainMapping{
			{
				Metadata: &run.ObjectMeta{Name: "d1"},
				Spec:     &run.DomainMappingSpec{RouteName: "s1"},
			},
			{
				Metadata: &run.ObjectMeta{Name: "d2"},
				Spec:     &run.DomainMappingSpec{RouteName: "s2"},
			},
		}, nil
	}

	dms, err := List("p", "r")

	assert.NoError(t, err)
	assert.Len(t, dms, 2)
	assert.Equal(t, "d1", dms[0].Name)
	assert.Equal(t, "s1", dms[0].RouteName)
	assert.Equal(t, "d2", dms[1].Name)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListDomainMappingsFunc = func(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
		return nil, assert.AnError
	}

	dms, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, dms)
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListDomainMappingsFunc = func(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
		if region == "us-central1" {
			return []*run.DomainMapping{
				{Metadata: &run.ObjectMeta{Name: "dm1"}},
			}, nil
		}
		return nil, nil
	}

	dms, err := List("p", "all")
	assert.NoError(t, err)
	
	// Since api_region.List() returns many regions, we just want to ensure we called List for them and aggregated results.
	// We mocked return for "us-central1".
	found := false
	for _, dm := range dms {
		if dm.Name == "dm1" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find dm1 from us-central1")
}

// --- GCPClient Tests ---

type MockDomainMappingsClientWrapper struct {
	ListFunc func(parent string, pageToken string) (*run.ListDomainMappingsResponse, error)
}

func (m *MockDomainMappingsClientWrapper) List(parent string, pageToken string) (*run.ListDomainMappingsResponse, error) {
	if m.ListFunc != nil {
		return m.ListFunc(parent, pageToken)
	}
	return nil, nil
}

func TestGCPClient_ListDomainMappings(t *testing.T) {
	// Mock dependencies
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createClient = func(ctx context.Context, creds *google.Credentials) (DomainMappingsClientWrapper, error) {
			return &MockDomainMappingsClientWrapper{
				ListFunc: func(parent string, pageToken string) (*run.ListDomainMappingsResponse, error) {
					return &run.ListDomainMappingsResponse{
						Items: []*run.DomainMapping{{Metadata: &run.ObjectMeta{Name: "dm1"}}},
					}, nil
				},
			}, nil
		}

		c := &GCPClient{}
		dms, err := c.ListDomainMappings(context.Background(), "p", "r")
		assert.NoError(t, err)
		assert.Len(t, dms, 1)
		assert.Equal(t, "dm1", dms[0].Metadata.Name)
	})

	t.Run("Pagination", func(t *testing.T) {
		createClient = func(ctx context.Context, creds *google.Credentials) (DomainMappingsClientWrapper, error) {
			return &MockDomainMappingsClientWrapper{
				ListFunc: func(parent string, pageToken string) (*run.ListDomainMappingsResponse, error) {
					if pageToken == "" {
						return &run.ListDomainMappingsResponse{
							Items: []*run.DomainMapping{{Metadata: &run.ObjectMeta{Name: "dm1"}}},
							Metadata: &run.ListMeta{Continue: "next-page"},
						}, nil
					}
					if pageToken == "next-page" {
						return &run.ListDomainMappingsResponse{
							Items: []*run.DomainMapping{{Metadata: &run.ObjectMeta{Name: "dm2"}}},
						}, nil
					}
					return nil, errors.New("unexpected page token")
				},
			}, nil
		}

		c := &GCPClient{}
		dms, err := c.ListDomainMappings(context.Background(), "p", "r")
		assert.NoError(t, err)
		assert.Len(t, dms, 2)
		assert.Equal(t, "dm1", dms[0].Metadata.Name)
		assert.Equal(t, "dm2", dms[1].Metadata.Name)
	})

	t.Run("AuthError", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}
		c := &GCPClient{}
		_, err := c.ListDomainMappings(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find default credentials")
	})

	t.Run("ClientCreationError", func(t *testing.T) {
        // Reset auth mock for this test
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createClient = func(ctx context.Context, creds *google.Credentials) (DomainMappingsClientWrapper, error) {
			return nil, errors.New("creation failed")
		}
		c := &GCPClient{}
		_, err := c.ListDomainMappings(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create domain mappings client")
	})

	t.Run("ListError", func(t *testing.T) {
        // Reset auth mock for this test
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createClient = func(ctx context.Context, creds *google.Credentials) (DomainMappingsClientWrapper, error) {
			return &MockDomainMappingsClientWrapper{
				ListFunc: func(parent string, pageToken string) (*run.ListDomainMappingsResponse, error) {
					return nil, errors.New("list failed")
				},
			}, nil
		}
		c := &GCPClient{}
		_, err := c.ListDomainMappings(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list domain mappings")
	})
}
