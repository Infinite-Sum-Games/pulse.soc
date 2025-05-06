package controllers

import (
	pkg "github.com/IAmRiteshKoushik/pulse/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginUserCsrf(c *gin.Context) {

	token, err := pkg.NewCsrfToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		return
	}

	// TODO: Setup cookie

	c.JSON(http.StatusOK, gin.H{
		"message": "Token generated successfully.",
		"token":   token,
	})
	return
}

func LoginUser(c *gin.Context) {

}

func RegisterUserAccountCsrf(c *gin.Context) {

	token, err := pkg.NewCsrfToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		return
	}

	// TODO: Setup cookie

	c.JSON(http.StatusOK, gin.H{
		"message": "Token generated successfully.",
		"token":   token,
	})
	return
}

func RegisterUserAccount(c *gin.Context) {

}

func ResendUserOtpCsrf(c *gin.Context) {

	token, err := pkg.NewCsrfToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		return
	}

	// TODO: Setup cookie

	c.JSON(http.StatusOK, gin.H{
		"message": "Token generated successfully.",
		"token":   token,
	})
	return
}

func ResendUserOtp(c *gin.Context) {

}

func UserSession(c *gin.Context) {

}

func LogoutUser(c *gin.Context) {

}
