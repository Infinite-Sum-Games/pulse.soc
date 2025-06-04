package pkg

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GrabUsername(c *gin.Context) (string, bool) {
	username, ok := c.Get("username")
	if !ok {
		return "", ok
	}
	if fmt.Sprintf("%T", username) == "string" {
		return fmt.Sprintf("%s", username), true
	}
	fmt.Printf("%T is the type of %s", username, username)
	return "", false
}
