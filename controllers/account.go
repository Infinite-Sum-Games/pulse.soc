package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	db "github.com/IAmRiteshKoushik/pulse/db/gen"
	"github.com/IAmRiteshKoushik/pulse/pkg"
	"github.com/gin-gonic/gin"
)

func FetchUserAccount(c *gin.Context) {
	username, ok := pkg.GrabUsername(c)
	if !ok {
		cmd.Log.Warn(
			fmt.Sprintf(
				"Username did not set in Gin-Context post Authentication at %s %s",
				c.Request.Method,
				c.FullPath(),
			),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		pkg.DbError(c, err)
		return
	}
	defer conn.Release()

	q := db.New()
	userProfile, err := q.FetchProfileQuery(ctx, conn, username)
	if err != nil {
		pkg.DbError(c, err)
		return
	}

	rank, err := pkg.GetParticipantRank(username)
	if err != nil {
		cmd.Log.Error(
			fmt.Sprintf("Failed to fetch rank at %s %s",
				c.Request.Method,
				c.FullPath(),
			),
			err,
		)
		return
	}

	cmd.Log.Info(
		fmt.Sprintf("Successfully retrived user profile at %s %s",
			c.Request.Method, c.FullPath(),
		),
	)
	c.JSON(http.StatusOK, gin.H{
		"message":             "User profile retrived successfully",
		"github_username":     userProfile.Ghusername,
		"first_name":          userProfile.FirstName,
		"middle_name":         userProfile.MiddleName.String,
		"last_name":           userProfile.LastName,
		"bounty":              userProfile.Bounty,
		"pull_request_count":  userProfile.SolutionCount,
		"pending_issue_count": userProfile.ActiveClaims,
		"rank":                rank,
		"badges":              userProfile.Badges,
	})
}
