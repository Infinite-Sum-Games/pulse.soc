package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/IAmRiteshKoushik/pulse/sse"
	"github.com/redis/go-redis/v9"
)

func ReadStream(lv *sse.LiveServer) {
	// Read stream always begins from the latest set of events. If there is a
	// server restart or a panic then a recover-and-restart mechanism is
	// required inorder to continue operations. In between these, if there were
	// events in the stream that were not broadcasted, then the client must
	// refetch the latest 5 events and reconnect to the stream.

	// The stream is not responsible for transmitting lost events in-case of
	// failures. That is handled by the FetchLatestUpdates endpoint.
	lastID := "$"
	for {
		args := &redis.XReadArgs{
			Streams: []string{LiveUpdateStream, lastID},
			Count:   10,
			Block:   0, // indefinite block
		}
		streams, err := Valkey.XRead(context.Background(), args).Result()
		if err != nil {
			cmd.Log.Error("Failed to read from stream. Retrying in 2 seconds...", err)
			// Backoff and then retry mechanism
			time.Sleep(2 * time.Second)
			continue
		}

		// Extract and process stream entries
		for _, stream := range streams {
			for _, message := range stream.Messages {
				lastID = message.ID
				for _, val := range message.Values {
					data, ok := val.(string)
					if !ok {
						continue
					}
					liveUpdate := fmt.Sprintf("data: %s\n\n", data)
					lv.Broadcast <- liveUpdate
				}
			}
		}
	}
}
