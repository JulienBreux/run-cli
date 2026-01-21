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

package log

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

var pollInterval = 2 * time.Second

// StreamLogs streams logs for a given project and filter to the provided channel.
// It first sends the last 50 logs, then polls for new ones.
func StreamLogs(ctx context.Context, projectID, filter string, logChan chan<- string) error {
	client, err := clientFactory(ctx, projectID)
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
	ticker := time.NewTicker(pollInterval)
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
	ch <- formatEntry(entry)
}

func formatEntry(entry *logging.Entry) string {
	payload := fmt.Sprintf("%v", entry.Payload)
	ts := entry.Timestamp.Format("15:04:05")
	return fmt.Sprintf("[%s] %s", ts, payload)
}
