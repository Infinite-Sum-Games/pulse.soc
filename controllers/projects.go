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
	"github.com/google/uuid"
)

func FetchProjects(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		pkg.DbError(c, err)
		return
	}
	defer conn.Release()

	q := db.New()
	results, err := q.FetchAllProjectsQuery(ctx, conn)
	if err != nil {
		pkg.DbError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Projects retrived successfully",
		"projects": results,
	})
	cmd.Log.Info(fmt.Sprintf(
		"[SUCCESS]: Processed request at %s %s",
		c.Request.Method, c.FullPath(),
	))
}

func FetchIssues(c *gin.Context) {
	projectIdParam := c.Param("projectId")
	projectId, err := uuid.Parse(projectIdParam)
	if err != nil {
		cmd.Log.Error(
			fmt.Sprintf("[INVALID-ID]: Given project-id is invalid UUID at %s %s",
				c.Request.Method, c.FullPath()), err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid project-id. Could not fetch issues.",
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
	ok, err := q.CheckIfProjectExistsQuery(ctx, conn, projectId)
	if err != nil {
		pkg.DbError(c, err)
		return
	}
	if !ok {
		cmd.Log.Error(
			fmt.Sprintf("[INVALID-ID]: No project with given project-id exists at %s %s",
				c.Request.Method, c.FullPath()), err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid project-id. Could not fetch issues.",
		})
		return
	}

	results, err := q.FetchAllIssuesByProjectIdQuery(ctx, conn, projectId)
	if err != nil {
		pkg.DbError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Issues retrived successfully",
		"issues":  results,
	})
	cmd.Log.Info(fmt.Sprintf(
		"[SUCCESS]: Processed request at %s %s",
		c.Request.Method, c.FullPath(),
	))
}
