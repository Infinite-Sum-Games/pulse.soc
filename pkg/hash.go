package pkg

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	cost := 14
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	hashedPassword := string(hash)
	return hashedPassword, nil
}
