package pkg

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	db "github.com/IAmRiteshKoushik/pulse/db/gen"
	"github.com/gin-gonic/gin"
)

func SetAuthCookie(c *gin.Context, authTokenString string) {
	c.SetCookie(
		"access_token",           // key
		authTokenString,          // value
		3600*24*90,               // maxAge (90 days)
		"/",                      // path
		cmd.EnvVars.Domain,       // domain
		cmd.EnvVars.CookieSecure, // secure
		true,                     // httpOnly
	)
}

func SetRefreshCookie(c *gin.Context, refreshTokenString string) {
	c.SetCookie(
		"refresh_token",          // key
		refreshTokenString,       // value
		3600*24*90,               // maxAge (90 days)
		"/",                      // path
		cmd.EnvVars.Domain,       // domain
		cmd.EnvVars.CookieSecure, // secure
		true,                     // httpOnly
	)
}

func SetCsrfCookie(c *gin.Context, csrfTokenString string, path string) {
	c.SetCookie(
		"csrf_token",             // key
		csrfTokenString,          // value
		300,                      // maxAge (5 minutes)
		path,                     // path to be constructed for restriction
		cmd.EnvVars.Domain,       // domain
		cmd.EnvVars.CookieSecure, // secure
		true,                     // httpOnly
	)
}

func NullifyCookies(c *gin.Context) {

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refesh_token", "", -1, "/", "", false, true)
	c.SetCookie("csrf_token", "", -1, "/", "", false, true)

	// If there is an error saying that there is no cookie then we are good
	// Otherwise we are in problem because nullification failed
	_, err := c.Cookie("access_token")
	if err != http.ErrNoCookie {
		cmd.Log.Fatal("[AUTH-ERROR] Failed to Nullify Cookies", err)
		return
	}
	_, err = c.Cookie("refresh_token")
	if err != http.ErrNoCookie {
		cmd.Log.Fatal("[AUTH-ERROR] Failed to Nullify Cookies", err)
		return
	}
	_, err = c.Cookie("csrf_token")
	if err != http.ErrNoCookie {
		cmd.Log.Fatal("[AUTH-ERROR] Failed to Nullify Cookies", err)
		return
	}

	email, exists := c.Get("email")
	if exists != true {
		cmd.Log.Fatal("[AUTH-ERROR] Failed to Revoke Refresh Token in DB",
			fmt.Errorf("Email fetch failed from gin.Context"))
		return
	}
	RevokeRefreshToken(email.(string))
	return
}

/*
 * We are revoking the refresh-token so that you cannot use it to get any more
 * Auth-Rokens in-case you have managed to steal the token from the browser
 * and kept it somewhere
 */
func RevokeRefreshToken(email string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := cmd.DBPool.Begin(ctx)
	if err != nil {
		cmd.Log.Fatal("[AUTH-ERROR] Failed to Revoke Refresh Token in DB", err)
		return
	}
	defer tx.Rollback(ctx)

	q := db.New()
	result, err := q.RevokeRefreshTokenQuery(ctx, tx, email)
	if err != nil || result.String != "" {
		cmd.Log.Fatal("[AUTH-ERROR] Failed to Revoke Refresh Token in DB", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		cmd.Log.Fatal("[AUTH-ERROR] Failed to Revoke Refresh Token in DB", err)
		return
	}
	cmd.Log.Info("[AUTH-INFO] Successfully revoked Refresh Token in DB")
	return
}
