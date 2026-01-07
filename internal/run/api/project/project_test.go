package project

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_WithoutCredentials(t *testing.T) {
	if hasGCloudCredentials() {
		t.Skip("Skipping test when credentials are available")
	}

	projects, err := List()

	assert.Error(t, err)
	assert.Nil(t, projects)
	assert.Contains(t, err.Error(), "failed to find default credentials")
}

func TestList_WithCredentials(t *testing.T) {
	if !hasGCloudCredentials() {
		t.Skip("Skipping test without gcloud credentials")
	}

	projects, err := List()

	if err != nil {
		assert.Contains(t, err.Error(), "authentication")
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, projects)
}

func TestList_ReturnsProjectsWithNameAndNumber(t *testing.T) {
	if !hasGCloudCredentials() {
		t.Skip("Skipping test without gcloud credentials")
	}

	projects, err := List()
	if err != nil {
		t.Skip("Skipping due to authentication error")
	}

	if len(projects) > 0 {
		project := projects[0]
		assert.NotEmpty(t, project.Name)
		assert.GreaterOrEqual(t, project.Number, 0)
	}
}

func hasGCloudCredentials() bool {
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath != "" {
		_, err := os.Stat(credPath)
		return err == nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	adcPath := homeDir + "/.config/gcloud/application_default_credentials.json"
	_, err = os.Stat(adcPath)
	return err == nil
}

