package log

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// MockClient is a mock implementation of Client.
type MockClient struct {
	EntriesFunc func(ctx context.Context, opts ...interface{}) EntryIterator
	CloseFunc   func() error
}

func (m *MockClient) Entries(ctx context.Context, opts ...interface{}) EntryIterator {
	if m.EntriesFunc != nil {
		return m.EntriesFunc(ctx, opts...)
	}
	return &MockEntryIterator{}
}

func (m *MockClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockEntryIterator is a mock implementation of EntryIterator.
type MockEntryIterator struct {
	Items []*logging.Entry
	Index int
	Err   error
}

func (m *MockEntryIterator) Next() (*logging.Entry, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Index >= len(m.Items) {
		return nil, iterator.Done
	}
	item := m.Items[m.Index]
	m.Index++
	return item, nil
}

func TestFormatEntry(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:00Z")
	entry := &logging.Entry{
		Timestamp: ts,
		Payload:   "Log message",
	}

	result := formatEntry(entry)
	assert.Equal(t, "[10:00:00] Log message", result)
}

func TestStreamLogs(t *testing.T) {
	// Restore clientFactory after test
	origFactory := clientFactory
	defer func() { clientFactory = origFactory }()

	t.Run("Backlog and Polling", func(t *testing.T) {
		// Speed up polling
		origInterval := pollInterval
		pollInterval = 10 * time.Millisecond
		defer func() { pollInterval = origInterval }()

		// Prepare mock data
		ts1, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:00Z")
		ts2, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:01Z")
		ts3, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:02Z")

		// Backlog: Newest First (so 2 then 1)
		backlogLogs := []*logging.Entry{
			{Timestamp: ts2, Payload: "Log 2"},
			{Timestamp: ts1, Payload: "Log 1"},
		}

		// New logs: Oldest First (so 3)
		newLogs := []*logging.Entry{
			{Timestamp: ts3, Payload: "Log 3"},
		}

		// Refined Mock Logic
		callCount := 0
		clientFactory = func(ctx context.Context, projectID string) (Client, error) {
			return &MockClient{
				EntriesFunc: func(ctx context.Context, opts ...interface{}) EntryIterator {
					callCount++
					if callCount == 1 {
						return &MockEntryIterator{Items: backlogLogs}
					}
					if callCount == 2 {
						return &MockEntryIterator{Items: newLogs}
					}
					return &MockEntryIterator{Items: nil} // No more logs
				},
			}, nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		logChan := make(chan string, 10)

		// Run StreamLogs in a goroutine
		errChan := make(chan error)
		go func() {
			errChan <- StreamLogs(ctx, "project", "filter", logChan)
		}()

		// Allow some time for processing
		time.Sleep(100 * time.Millisecond) // Wait for backlog

		// Check Backlog (Should be reversed: Log 1, Log 2)
		select {
		case msg := <-logChan:
			assert.Contains(t, msg, "Log 1")
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for Log 1")
		}

		select {
		case msg := <-logChan:
			assert.Contains(t, msg, "Log 2")
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for Log 2")
		}

		// Check Polling (Log 3)
		select {
		case msg := <-logChan:
			assert.Contains(t, msg, "Log 3")
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for Log 3")
		}

		cancel() // Stop the loop
		<-errChan
	})

	t.Run("Client Creation Error", func(t *testing.T) {
		expectedErr := errors.New("client error")
		clientFactory = func(ctx context.Context, projectID string) (Client, error) {
			return nil, expectedErr
		}

		logChan := make(chan string)
		err := StreamLogs(context.Background(), "p", "f", logChan)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestStreamLogs_PollingError(t *testing.T) {
	// Restore clientFactory after test
	origFactory := clientFactory
	defer func() { clientFactory = origFactory }()

	// Speed up polling
	origInterval := pollInterval
	pollInterval = 10 * time.Millisecond
	defer func() { pollInterval = origInterval }()

	callCount := 0
	clientFactory = func(ctx context.Context, projectID string) (Client, error) {
		return &MockClient{
			EntriesFunc: func(ctx context.Context, opts ...interface{}) EntryIterator {
				callCount++
				if callCount == 1 {
					// Backlog: empty
					return &MockEntryIterator{Items: nil}
				}
				// Polling: Error
				return &MockEntryIterator{Err: errors.New("polling error")}
			},
		}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logChan := make(chan string)

	// StreamLogs should not return error, but just retry or stop?
	// The code breaks the inner loop on error, then waits for next tick.
	// So it should continue running.

	// We want to verify it DOES NOT exit StreamLogs.
	done := make(chan error)
	go func() {
		done <- StreamLogs(ctx, "p", "f", logChan)
	}()

	select {
	case <-done:
		t.Fatal("StreamLogs exited unexpectedly on polling error")
	case <-time.After(100 * time.Millisecond):
		// passed
	}
}

// --- Mocks for GCPClient testing ---

type MockLogAdminClientWrapper struct {
	EntriesFunc func(ctx context.Context, opts ...logadmin.EntriesOption) EntryIterator
	CloseFunc   func() error
}

func (m *MockLogAdminClientWrapper) Entries(ctx context.Context, opts ...logadmin.EntriesOption) EntryIterator {
	if m.EntriesFunc != nil {
		return m.EntriesFunc(ctx, opts...)
	}
	return &MockEntryIterator{}
}

func (m *MockLogAdminClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestGCPClient(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createLogAdminClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createLogAdminClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("NewGCPClient_Success", func(t *testing.T) {
		createLogAdminClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (LogAdminClientWrapper, error) {
			return &MockLogAdminClientWrapper{
				EntriesFunc: func(ctx context.Context, opts ...logadmin.EntriesOption) EntryIterator {
					return &MockEntryIterator{
						Items: []*logging.Entry{{Payload: "log1"}},
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		client, err := NewGCPClient(context.Background(), "project")
		assert.NoError(t, err)

		it := client.Entries(context.Background())
		entry, err := it.Next()
		assert.NoError(t, err)
		assert.Equal(t, "log1", entry.Payload)
		
		assert.NoError(t, client.Close())
	})

	t.Run("NewGCPClient_AuthError", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}
		
		_, err := NewGCPClient(context.Background(), "project")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find default credentials")
	})

	t.Run("NewGCPClient_ClientCreationError", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createLogAdminClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (LogAdminClientWrapper, error) {
			return nil, errors.New("creation failed")
		}
		
		_, err := NewGCPClient(context.Background(), "project")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "creation failed")
	})
}

func TestWrappers_Delegation(t *testing.T) {
	// Expect panics because nil clients are used
	
	t.Run("RealLogAdminClient", func(t *testing.T) {
		w := &RealLogAdminClient{client: nil}
		assert.Panics(t, func() { w.Entries(context.Background()) })
		assert.Panics(t, func() { w.Close() })
	})
	
	t.Run("GCPEntryIterator", func(t *testing.T) {
		it := &GCPEntryIterator{it: nil}
		assert.Panics(t, func() { it.Next() })
	})
}