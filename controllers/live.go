package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/IAmRiteshKoushik/pulse/pkg"
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
	// Setup SSE-specific headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	keepAlive := time.NewTicker(15 * time.Second)
	defer keepAlive.Stop()

	c.JSON(http.StatusOK, gin.H{
		"message": "LIVE updates are a work-in-progress",
	})
	cmd.Log.Info(fmt.Sprintf(
		"[SUCCESS]: Processed request at %s %s",
		c.Request.Method, c.FullPath(),
	))
}
