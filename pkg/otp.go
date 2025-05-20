package pkg

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP() (string, error) {
	const otpLength = 6
	const digits = "0123456789"
	otp := make([]byte, otpLength)

	for i := range otpLength {
		randomIndex, err := rand.Int(
			rand.Reader,
			big.NewInt(
				int64(len(digits)),
			),
		)

		if err != nil {
			return "", err
		}
		otp[i] = digits[randomIndex.Int64()]
	}

	return string(otp), nil
}
