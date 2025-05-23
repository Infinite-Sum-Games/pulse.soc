package pkg

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/gin-gonic/gin"
)

func DbError(c *gin.Context, err error) {
	if err == context.DeadlineExceeded {
		cmd.Log.Warn(
			fmt.Sprintf("[CONTEXT-DEADLINE-EXCEEDED]: Server is experiencing delays at %s %s",
				c.Request.Method,
				c.FullPath(),
			),
		)
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": "The server is experiencing delays. Try again later.",
		})
	} else {
		cmd.Log.Fatal(
			fmt.Sprintf(
				"[INTERNAL-SERVER-ERROR]: DB Error at %s %s.\n",
				c.Request.Method,
				c.FullPath()),
			err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Oops! Something happened. Please try again later.",
		})
	}
}

func JSONUnmarshallError(c *gin.Context, err error) {
	cmd.Log.Error(
		fmt.Sprintf(
			"[REQUEST-ERROR] Unmarshalling failed at %s %s.\n",
			c.Request.Method,
			c.FullPath(),
		), err)
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "The request is malformed",
	})
	return
}

func RequestValidatorError(c *gin.Context, err error) {
	cmd.Log.Error(
		fmt.Sprintf(
			"[REQUEST-ERROR] Validation failed at %s %s",
			c.Request.Method,
			c.FullPath(),
		), err)
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "The request is malformed.",
	})
	return
}
