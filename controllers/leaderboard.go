package controllers

import (
	"fmt"
	"net/http"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/gin-gonic/gin"
)

func FetchLeaderboard(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Leaderboard WIP",
	})
	cmd.Log.Info(fmt.Sprintf(
		"[SUCCESS]: Processed request at %s %s",
		c.Request.Method, c.FullPath(),
	))
	return
}
