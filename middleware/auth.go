package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/IAmRiteshKoushik/pulse/pkg"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	if !strings.HasPrefix(c.Request.RequestURI, "/api/v1") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Server could not understand the URL",
		})
		return
	}

	RefreshCookie, refErr := c.Cookie("refresh_token")
	if refErr == http.ErrNoCookie {
		cmd.Log.Error(
			fmt.Sprintf("[AUTH-ERROR]: RefreshToken not found at %s %s",
				c.Request.Method,
				c.FullPath(),
			), refErr)
		pkg.NullifyCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization denied",
		})
		return
	}

	_, authErr := c.Cookie("access_token")
	if authErr == http.ErrNoCookie {
		token, err := pkg.CheckForRefreshToken(RefreshCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization denied",
			})
			return
		}
		claims := token.Claims()
		_ = pkg.CreateAuthToken(
			// TODO: Type assertion to get around the error for now.
			// Need a more robust solution later on.
			claims["audience"].(string),
			claims["jti"].(string),
			claims["USER-ROLE"].(bool),
			claims["HOST-ROLE"].(bool),
			claims["STAFF-ROLE"].(bool),
		)
		// TODO: After creating new tokens, setup the tokens in the cookies
	}

	c.Next()
}
