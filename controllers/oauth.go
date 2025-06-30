package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	db "github.com/IAmRiteshKoushik/pulse/db/gen"
	"github.com/IAmRiteshKoushik/pulse/pkg"
	"github.com/IAmRiteshKoushik/pulse/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func InitiateGitHubOAuth(c *gin.Context) {
	url := cmd.GithubOAuthConfig.AuthCodeURL("")
	c.Redirect(http.StatusTemporaryRedirect, url)
	cmd.Log.Info("Successfully initiated GitHub oAuth at GET /api/v1/auth/github")
}

func CompleteGitHubOAuth(c *gin.Context) {
	redirectUrl := cmd.AppConfig.FrontendURL + "/register"

	// Extract code from github oauth callback URL
	code := c.Query("code")
	if code == "" {
		cmd.Log.Warn(
			fmt.Sprintf("Missing authorization code in github oauth callback at %s %s",
				c.Request.Method, c.FullPath()))
		c.Redirect(http.StatusTemporaryRedirect, cmd.AppConfig.FrontendURL)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetching the github user
	token, err := cmd.GithubOAuthConfig.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("%v", cmd.GithubOAuthConfig)
		cmd.Log.Error(
			fmt.Sprintf("Failed to exchange code for token at %s %s",
				c.Request.Method, c.FullPath()), err)
		c.Redirect(http.StatusTemporaryRedirect, cmd.AppConfig.FrontendURL)
		return
	}

	client := cmd.GithubOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		cmd.Log.Warn(
			fmt.Sprintf("Failed to fetch user info from GitHub at %s %s",
				c.Request.Method, c.FullPath()))
		c.Redirect(http.StatusTemporaryRedirect, cmd.AppConfig.FrontendURL)
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "Oops! Something happened. Please try again later",
		// })
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cmd.Log.Warn(fmt.Sprintf("Failed to unmarshal github user info at %s %s",
			c.Request.Method, c.FullPath()))
		c.Redirect(http.StatusTemporaryRedirect, cmd.AppConfig.FrontendURL)
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "Oops! Something happened. Please try again later",
		// })
		return
	}
	// Extracting the github user
	var user types.GithubUser
	if err := json.Unmarshal(body, &user); err != nil {
		cmd.Log.Error(
			fmt.Sprintf("Failed to parse github user info at %s %s",
				c.Request.Method, c.FullPath()), err)
		c.Redirect(http.StatusTemporaryRedirect, cmd.AppConfig.FrontendURL)
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "Oops! Something happened. Please try again later",
		// })
		return
	}

	// Verifying the github account's presence against database to validate
	// post registration
	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		// pkg.DbError(c, err)
		c.Redirect(http.StatusTemporaryRedirect, cmd.AppConfig.FrontendURL)
		return
	}
	defer tx.Rollback(ctx)

	q := db.New()
	userExist, err := q.RetriveExistingUserQuery(ctx, tx, user.Username)
	if err != nil {
		// pkg.DbError(c, err)
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}

	// If the presence is verified, then generate access and refresh token
	// , add them in DB and respond back in request
	accessToken, err := pkg.CreateToken(userExist.Ghusername, userExist.Email, "access_token")
	if err != nil {
		cmd.Log.Error(
			fmt.Sprintf("Failed to create access token at %s %s", c.Request.Method, c.FullPath()),
			err)
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "Oops! Something happened. Please try again later",
		// })
		return
	}
	refreshToken, err := pkg.CreateToken(userExist.Ghusername, userExist.Email, "refresh_token")
	if err != nil {
		cmd.Log.Error(
			fmt.Sprintf("Failed to create token at %s %s", c.Request.Method, c.FullPath()),
			err)
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "Oops! Something happened. Please try again later",
		// })
		return
	}

	loginUser, err := q.AddRefreshTokenQuery(ctx, tx, db.AddRefreshTokenQueryParams{
		Ghusername:   userExist.Ghusername,
		RefreshToken: pgtype.Text{String: refreshToken, Valid: true},
	})
	if err != nil {
		cmd.Log.Error(
			fmt.Sprintf("Failed to process request at %s %s", c.Request.Method, c.FullPath()), err)
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		cmd.Log.Error(fmt.Sprintf("Failed to process request at %s %s", c.Request.Method, c.FullPath()), err)
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}

	// Redirect back to frontend with tokens
	frontendURL := fmt.Sprintf("%s/auth/callback", cmd.AppConfig.FrontendURL)
	redirectURL := fmt.Sprintf("%s?access_token=%s&refresh_token=%s&github_username=%s&email=%s&bounty=%d",
		frontendURL,
		accessToken,
		refreshToken,
		loginUser.Ghusername,
		loginUser.Email,
		loginUser.Bounty,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	cmd.Log.Info(fmt.Sprintf(
		"[SUCCESS]: Processed request at %s %s",
		c.Request.Method, c.FullPath(),
	))
}

func RegenerateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		cmd.Log.Warn(
			fmt.Sprintf("RefreshToken not sent as Authorization header at %s %s",
				c.Request.Method, c.FullPath()),
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Authorization header is missing in request",
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
		cmd.Log.Warn(
			fmt.Sprintf("Invalid refresh token at %s %s",
				c.Request.Method, c.FullPath()),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "The request is malformed",
		})
		return
	}

	validIssuer := claims.Issuer == "api.season-of-code"
	validSub := claims.Subject == "refresh_token"
	validAudience := len(claims.Audience) == 1
	if !validIssuer || !validSub || !validAudience {
		cmd.Log.Error(
			fmt.Sprintf("Tampered token sent at %s %s", c.Request.Method, c.FullPath()),
			err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Server refused to process the request",
		})
		return
	}

	// Actual controller
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		pkg.DbError(c, err)
		return
	}
	defer conn.Release()

	q := db.New()
	result, err := q.CheckRefreshTokenQuery(ctx, conn, db.CheckRefreshTokenQueryParams{
		Email:        claims.ID,
		RefreshToken: pgtype.Text{String: tokenString, Valid: true},
	})
	if err != nil {
		pkg.DbError(c, err)
		return
	}

	accessToken, err := pkg.CreateToken(result.Ghusername, result.Email, "access_token")
	if err != nil {
		cmd.Log.Error(
			fmt.Sprintf("Could not generate access token at %s %s", c.Request.Method, c.FullPath()),
			err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops! Something happened. Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Token refreshed successfully",
		"accessKey": accessToken,
	})
	cmd.Log.Info(fmt.Sprintf(
		"[SUCCESS]: Processed request at %s %s",
		c.Request.Method, c.FullPath(),
	))
}
