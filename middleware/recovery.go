package middleware

import (
	"fmt"
	"net/http"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			cmd.Log.Fatal(
				fmt.Sprintf("[PANIC-RECOVERED]: Panic occured at %s %s", c.Request.Method, c.FullPath()),
				fmt.Errorf("%v\n", err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Oops! Something happened. Please try again later.",
			})
		}
	}()
	c.Next()
}
