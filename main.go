package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	cmd "github.com/IAmRiteshKoushik/pulse/cmd"
	c "github.com/IAmRiteshKoushik/pulse/controllers"
	mw "github.com/IAmRiteshKoushik/pulse/middleware"
	pkg "github.com/IAmRiteshKoushik/pulse/pkg"
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

	// Initialize Server
	server := InitServer()
	port := strconv.Itoa(cmd.EnvVars.Port)
	cmd.Log.Info("[OK]: Server configured and starting on PORT " + port)
	err = server.Run(":" + port)
	if err != nil {
		panic(fmt.Errorf(failMsg, err))
	}
}

func InitServer() *gin.Engine {

	// Gin Configurations for Logging, Monitoring and Panic-Recovery
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

	v1 := router.Group("/api")
	Routes(v1)
	return router
}

func Routes(engine *gin.RouterGroup) {

	engine.POST("/auth/github", c.InitiateGitHubOAuth)
	engine.POST("/auth/github/callback", c.CompleteGitHubOAuth)
	engine.POST("/auth/register", c.RegisterUserAccount)
	engine.POST("/auth/register/otp/verify", c.RegisterUserOtpVerify)
	engine.POST("/auth/register/otp/resend", c.RegisterUserOtpResend)

	engine.GET("/profile", c.FetchUserAccount)
	engine.GET("/leaderboard", c.FetchLeaderboard)
	engine.GET("/projects", c.FetchProjects)
	engine.GET("/issues/:projectId", c.FetchIssues)
	engine.GET("/updates/live", c.FetchLiveUpdates)
}

func main() {
	StartApp()
}
