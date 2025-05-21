package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	cmd "github.com/IAmRiteshKoushik/pulse/cmd"
	c "github.com/IAmRiteshKoushik/pulse/controllers"
	mw "github.com/IAmRiteshKoushik/pulse/middleware"
	pkg "github.com/IAmRiteshKoushik/pulse/pkg"
	"github.com/gin-gonic/gin"
)

func init() {
	failMsg := "Could not initialize app\n%w"

	// Initialize global environment variables
	env, err := cmd.NewEnvConfig()
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.EnvVars = env
	log.Println("[OK]: Environment variables configured successfully")

	// Initialize logger
	f, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	defer f.Close()
	cmd.Log = cmd.NewLoggerService(cmd.EnvVars.Environment, f)
	cmd.Log.Info("[OK]: Logging service configured successfully.")

	// Initialize Regex
	if err := pkg.InitRegex(); err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.Log.Info("[OK]: All Regular Expressions created successfully.")

	// Initialize database connection pool
	cmd.DBPool, err = cmd.InitDB()
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.Log.Info("[OK]: DB initialized successfully.")
}

func setup() *gin.Engine {
	ginLogs, err := os.Create("gin.log")
	if err != nil {
		cmd.Log.Fatal("Error creating log file for Gin", err)
		return nil
	}
	defer ginLogs.Close()
	multiWriter := io.MultiWriter(os.Stdout, ginLogs)
	gin.DefaultWriter = multiWriter
	gin.DefaultErrorWriter = multiWriter
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(mw.RecoveryMiddleware)
	router.Use(gin.Logger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is LIVE",
		})
		return
	})

	v1 := router.Group("/api/v1")

	v1.POST("/auth/github", c.InitiateGitHubOAuth)
	v1.POST("/auth/github/callback", c.CompleteGitHubOAuth)
	v1.POST("/auth/register", c.RegisterUserAccount)
	v1.POST("/auth/register/otp/verify", mw.Auth, c.RegisterUserOtpVerify)
	v1.GET("/auth/register/otp/resend", mw.Auth, c.RegisterUserOtpResend)
	v1.GET("/auth/refresh", c.RegenerateToken)

	v1.GET("/profile", mw.Auth, c.FetchUserAccount)
	v1.GET("/leaderboard", mw.Auth, c.FetchLeaderboard)
	v1.GET("/projects", mw.Auth, c.FetchProjects)
	v1.GET("/issues/:projectId", mw.Auth, c.FetchIssues)
	v1.GET("/updates/live", mw.Auth, c.FetchLiveUpdates)

	return router
}

func main() {
	server := setup()
	port := strconv.Itoa(cmd.EnvVars.Port)
	cmd.Log.Info("[OK]: Server configured and starting on PORT " + port)
	err := server.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
