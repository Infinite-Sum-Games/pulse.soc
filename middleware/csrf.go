package middleware

import (
	"fmt"
	"net/http"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/gin-gonic/gin"
)

func VerifyCsrf(c *gin.Context) {
	csrfFromHeader := c.GetHeader("X-Csrf-Token")
	if csrfFromHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing security token.",
		})
		return
	}

	csrfFromCookie, csrfErr := c.Cookie("csrf_token")
	if csrfErr == http.ErrNoCookie {
		cmd.Log.Error("[]", csrfErr)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing security token.",
		})
		return
	}

	if csrfFromCookie != csrfFromHeader {
		cmd.Log.Error(
			"[AUTH-ERROR]: Mismatching CSRF tokens",
			fmt.Errorf("ERR: CSRF tokens' do not match. Cookie: %s, Header: %s",
				csrfFromCookie,
				csrfFromHeader,
			))
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Security tokens do not match.",
		})
		return
	}
	c.Next()
}
