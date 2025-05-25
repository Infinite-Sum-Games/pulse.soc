package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	cmd "github.com/IAmRiteshKoushik/pulse/cmd"
	c "github.com/IAmRiteshKoushik/pulse/controllers"
	mw "github.com/IAmRiteshKoushik/pulse/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartApp() {
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

	// Initialize database connection pool
	cmd.DBPool, err = cmd.InitDB()
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
	cmd.Log.Info("[OK]: DB initialized successfully.")

	// Initialize oAuth
	cmd.OAuthInit()
	cmd.Log.Info("[OK]: GitHub oAuth configuration loaded successfully.")

	// Starting the server
	ginLogs, err := os.Create("gin.log")
	if err != nil {
		cmd.Log.Fatal("Error creating log file for Gin", err)
	}
	defer ginLogs.Close()
	multiWriter := io.MultiWriter(os.Stdout, ginLogs)
	gin.DefaultWriter = multiWriter
	gin.DefaultErrorWriter = multiWriter
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(mw.RecoveryMiddleware)
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cmd.EnvVars.FrontendUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is LIVE",
		})
		cmd.Log.Info(fmt.Sprintf(
			"[SUCCESS]: Processed request at %s %s",
			c.Request.Method, c.FullPath(),
		))
		return
	})

	v1 := router.Group("/api/v1")

	v1.GET("/auth/github", c.InitiateGitHubOAuth)
	v1.GET("/auth/github/callback", c.CompleteGitHubOAuth)
	v1.POST("/auth/register", c.RegisterUserAccount)
	v1.POST("/auth/register/otp/verify", mw.Auth, c.RegisterUserOtpVerify)
	v1.GET("/auth/register/otp/resend", mw.Auth, c.RegisterUserOtpResend)
	v1.GET("/auth/refresh", c.RegenerateToken)

	v1.GET("/profile", mw.Auth, c.FetchUserAccount)
	v1.GET("/leaderboard", mw.Auth, c.FetchLeaderboard)
	v1.GET("/projects", mw.Auth, c.FetchProjects)
	v1.GET("/issues/:projectId", mw.Auth, c.FetchIssues)
	v1.GET("/updates/live", mw.Auth, c.FetchLiveUpdates)

	port := strconv.Itoa(cmd.EnvVars.Port)
	cmd.Log.Info("[OK]: Server configured and starting on PORT " + port)
	err = router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

func main() {
	StartApp()
}
