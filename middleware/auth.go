package middleware

import (
	"fmt"
	"net/http"

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

	if !validIssuer || !validTokenType || !validEmail {
		cmd.Log.Warn(fmt.Sprintf("Tampered token detected at %s %s. Issuer: %v, Type: %v, Email: %v",
			c.Request.Method, c.FullPath(), validIssuer, validTokenType, validEmail))

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Server refused to process the request",
		})
		return
	}

	c.Set("email", claims.Email)
	c.Set("username", claims.GhUsername)
	c.Next()
}
