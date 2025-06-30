package pkg

import (
	"fmt"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(ghUsername, email, tokenType string) (string, error) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			ID:        email,
			Audience:  []string{ghUsername},
			Issuer:    "api.season-of-code",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiryAt),
			Subject:   tokenType,
		})

	tokenString, err := token.SignedString([]byte(cmd.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
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
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, fmt.Errorf("token expired")
		}
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims type")
}
