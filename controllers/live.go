package controllers

import (
	"fmt"
	"net/http"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/IAmRiteshKoushik/pulse/pkg"
	"github.com/IAmRiteshKoushik/pulse/sse"
	"github.com/gin-gonic/gin"
)

// Fetching the latest events before setting up a persistent uni-directional
// SSE connection
func FetchLatestUpdates(c *gin.Context) {
	updates, err := pkg.GetLatestLiveEvents(c)
	if err != nil {
		cmd.Log.Error(
			fmt.Sprintf("Failed to fetch latest updates at %s %s",
				c.Request.Method,
				c.FullPath(),
			),
			err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Latest updates fetches successfully",
		"updates":      updates,
		"update_count": len(updates),
	})
	cmd.Log.Info(fmt.Sprintf(
		"[SUCCESS]: Processed request at %s %s",
		c.Request.Method,
		c.FullPath(),
	))
}

// Handle server-sent-events with broadcast to multiple connections
// simulatenously from a Redis stream
func SetupLiveUpdates(c *gin.Context) {
	// NOTE: Setup SSE-specific headers
	// ================================
	// The text/event-stream tells the browser that the response is an
	// SSE stream with the SSE protocol. The specification is as follows
	// data: <message>\n\n
	// In the absence of this header, EventSource API in the browser
	// does not identify it and results in parser-failed to parse response.
	c.Writer.Header().Set("Content-Type", "text/event-stream")

	// As SSE streams are LIVE, their responses must not be cached by
	// the browser or proxies. Any caching would interfere with
	// the continuous delivery of new events.
	c.Writer.Header().Set("Cache-Control", "no-cache")

	// As the connection is a long-lived TCP connection, this header
	// tells the client and server to hold the HTTP connection open
	// for the entire duration of the SSE session. An absence would
	// lead to closing of connections by proxies, client or the server.
	c.Writer.Header().Set("Connection", "keep-alive")

	// This header disabled buffering by reverse proxies like Nginx
	// or Caddy. Proxies buffer responses to optimize the throughput of
	// the server but this goes orthogonal to SSE's streaming model.
	// Any buffering would delay / batch-events thus removing the
	// real-time nature of SSE. If it is missing then Nginx or Caddy
	// would deliver new evnets only when their buffer is full.
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// Creating a new client and adding it to the LiveServer
	// 100 is the number of un-read messages it can queue up
	// without requiring an immediate receiver. If it is unread
	// the further sends to this channel will be blocked.
	// There is a configuration in "sse/sse.go" where the client connection
	// would be forcibly dropped if it becomes unresponsive for 100
	// messages. This event can occur if there is network outage / slowdow
	// on the client side and it is not able to receive it at the servers'
	// broadcasting rate. If we had used an unbuffered channel then it would
	// have blocked the server from looping over the entire client-map to
	// broadcast messages to all clients due to a single unreponsive client.
	client := &sse.Client{Channel: make(chan string, 100)}
	sse.Live.Register <- client
	defer func() {
		// Clean up the connection in-case the client is terminated.

		// If this is not done then the clients-map grows infinitely.
		// Also, if this clean-up is not done then it would lead to
		// memory leaks and wasted CPU cycles because we would be
		// broadcasting messages and getting errors to
		// closed / unregistered client-connections.
		sse.Live.Unregister <- client
	}()

	// Blocking call using infinite-loop  for client-connection
	for {
		// Blockingly wait for multiple channel operations simulateneously
		// handling a new message from client.Channel or when the client
		// disconnects
		select {
		case msg, ok := <-client.Channel:
			// The message-type can be :
			// 1. data: <message>\n\n
			// 2. : keep-alive\n\n
			if !ok {
				// If a client disconnection happens then we close the
				// channel. In this case, ok gets the closing event
				// and the return is triggered to complete the
				// controller's task and trigger the clean up-logic
				// mentioned in "defer func(){}" above
				return
			}
			// Write the SSE-formatted JSON-string.
			_, err := c.Writer.WriteString(msg)
			if err != nil {
				// If the client connection is broken or closed while this
				// message is being transmitted then an error occurs and
				// return is triggered. Network fails, closing the browser
				// tab are common reasons for this to happen.
				return
			}
			// Once the writer has received a piece of data, it needs to be
			// pushed to the client. In-order to achieve this, we use Flush()
			// If it is not used, then the server's buffering mechanism would
			// cause a delay in delivery of the messages
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			// Any context.Context() in Golang has a Done() which is a channel
			// In this case, it is listening for client disconnections and
			// triggers the clean-up logic if it happens.
			return
		}
	}
}
