package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortCodeLength = 6

func GenerateShortCode() (string, error) {
	result := make([]byte, shortCodeLength)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func GetBaseURL(host string, isTLS bool) string {
	scheme := "http"
	if isTLS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}
