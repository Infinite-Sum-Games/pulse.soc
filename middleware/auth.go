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
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		cmd.Log.Warn(fmt.Sprintf("Authorization failed at %s %s", c.Request.Method, c.FullPath()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Authorization header required",
		})
		return
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[0:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		cmd.Log.Warn(fmt.Sprintf("Authorization failed at %s %s", c.Request.Method, c.FullPath()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Authorization header format",
		})
		return
	}

	claims, err := pkg.VerifyToken(tokenString)
	if err != nil {
		cmd.Log.Error(fmt.Sprintf("Authorization failed at %s %s", c.Request.Method, c.FullPath()),
			err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid or expired token",
		})
		return
	}

	validIssuer := claims.Issuer == "api.season-of-code"
	validTokenType := claims.TokenType == "access_token" || claims.TokenType == "temp_token"
	validEmail := claims.Email != ""
	validGithub := claims.GhUsername != ""

	if !validIssuer || !validTokenType || !validEmail {
		cmd.Log.Warn(fmt.Sprintf("Tampered token detected at %s %s. Issuer: %v, Type: %v, Email: %v",
			c.Request.Method, c.FullPath(), validIssuer, validTokenType, validEmail))

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Server refused to process the request",
		})
		return
	}
	// Allowing the GitHub linking endpoint even if the user has no Github Username yet
	// Check if the request path contains "/github"
	isLinkingEndpoint := strings.Contains(c.Request.URL.Path, "/github")

	if !validGithub && !isLinkingEndpoint {
		cmd.Log.Warn(
			fmt.Sprintf("User with incomplete profile attempting to use %s %s",
				c.Request.Method,
				c.FullPath(),
			),
		)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Please link github before proceeding",
		})
		return
	}

	c.Set("email", claims.Email)
	c.Set("username", claims.GhUsername)
	c.Next()
}
