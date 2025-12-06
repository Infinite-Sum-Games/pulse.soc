package pkg

import (
	"fmt"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/golang-jwt/jwt/v5"
)
type CustomClaims struct {
	Email      string `json:"email"`
	GhUsername string `json:"github_username"`
	TokenType  string `json:"token_type"`
	jwt.RegisteredClaims
}

func CreateToken(GhUsername, email, tokenType string) (string, error) {
	var expiryAt time.Time
	switch tokenType {
	case "temp_token":
		expiryAt = time.Now().Add(25 * time.Minute)
	case "access_token":
		expiryAt = time.Now().Add(4 * time.Hour)
	case "refresh_token":
		expiryAt = time.Now().Add(90 * 24 * time.Hour)
	default:
		return "", fmt.Errorf("invalid tokenType provided. valid types: %s, %s or %s",
			"temp_token", "access_token", "refresh_token")
	}

	claims := CustomClaims{
		Email:      email,
		GhUsername: GhUsername,
		TokenType:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "api.season-of-code",
			ExpiresAt: jwt.NewNumericDate(expiryAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cmd.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cmd.AppConfig.JWTSecret), nil
		})

	if err != nil {
		return nil, fmt.Errorf("token parsing error: %s", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims, ok := token.Claims.(*CustomClaims); ok {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, fmt.Errorf("token expired")
		}
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims type")
}
