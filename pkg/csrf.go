package pkg

import (
	"crypto/rand"
	"encoding/base64"
)

func NewCsrfToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	csrfToken := base64.StdEncoding.EncodeToString(token)
	return csrfToken, nil
}
