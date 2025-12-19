package log

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// StreamLogs streams logs for a given project and filter to the provided channel.
// It first sends the last 50 logs, then polls for new ones.
func StreamLogs(ctx context.Context, projectID, filter string, logChan chan<- string) error {
	// Explicitly find default credentials
	creds, err := google.FindDefaultCredentials(ctx, logging.ReadScope)
	if err != nil {
		return fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := logadmin.NewClient(ctx, projectID, option.WithCredentials(creds))
	if err != nil {
		return fmt.Errorf("failed to create logging client: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			logChan <- fmt.Sprintf("failed to close logging client: %v", err)
		}
	}()

	// 1. Fetch Initial Backlog (Last 50)
	// We use NewestFirst to get the 50 most recent, but we need to reverse them for display.
	iter := client.Entries(ctx, logadmin.Filter(filter), logadmin.NewestFirst())
	var backlog []*logging.Entry
	for len(backlog) < 50 {
		entry, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		backlog = append(backlog, entry)
	}

	var lastTimestamp time.Time

	// Send backlog (Reverse order: Oldest -> Newest)
	for i := len(backlog) - 1; i >= 0; i-- {
		entry := backlog[i]
		sendEntry(logChan, entry)
		if entry.Timestamp.After(lastTimestamp) {
			lastTimestamp = entry.Timestamp
		}
	}

	// 2. Poll for new logs
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// RFC3339Nano is important for precision
			tsFilter := fmt.Sprintf(`timestamp > "%s"`, lastTimestamp.Format(time.RFC3339Nano))
			newFilter := fmt.Sprintf("%s AND %s", filter, tsFilter)

			// Fetch new logs (OldestFirst is default and correct here)
			iter := client.Entries(ctx, logadmin.Filter(newFilter))
			for {
				entry, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					// On error, maybe just log to channel or ignore transient errors
					// For now, let's just break this poll loop and try again next tick
					break 
				}
				
				sendEntry(logChan, entry)
				if entry.Timestamp.After(lastTimestamp) {
					lastTimestamp = entry.Timestamp
				}
			}
		}
	}
}

func sendEntry(ch chan<- string, entry *logging.Entry) {
	payload := fmt.Sprintf("%v", entry.Payload)
	ts := entry.Timestamp.Format("15:04:05")
	ch <- fmt.Sprintf("[%s] %s", ts, payload)
}
