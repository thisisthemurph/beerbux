package otp

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghjkmnpqrstuvwxyz0123456789"

func Generate(length int) (string, error) {
	value := make([]byte, length)
	for i := range value {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		value[i] = charset[num.Int64()]
	}
	return string(value), nil
}
