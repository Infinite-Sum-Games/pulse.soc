package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			cmd.Log.Fatal(
				fmt.Sprintf("[PANIC-RECOVERED]: Panic occured at %s %s", c.Request.Method, c.FullPath()),
				fmt.Errorf("%v", err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Oops! Something happened. Please try again later.",
			})
		}
	}()
	c.Next()
}

// For goroutines which might panic, there is a requirement to restart them
// as they are critical to the application
func RecoverRoutine(name string, routine func()) {
	go func() {
		for {
			// Panic recovery block is wrapping the goroutine
			// In-case there is an exit, this function is triggered
			// and the panic is recoverred and logged.
			defer func() {
				if err := recover(); err != nil {
					cmd.Log.Fatal(name+" panicked.",
						fmt.Errorf("%v", err),
					)
					// Adding a delay to prevent a tight loop
					time.Sleep(5 * time.Second)
				}
			}()
			// Run the goroutine
			routine()
			// If the goroutine exists normally, then break out of the loop
			break
		}
	}()
}
