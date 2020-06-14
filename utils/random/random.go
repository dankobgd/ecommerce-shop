package random

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"io"
	"math/rand"
	"strings"
)

var (
	numeric              = []rune("0123456789")
	alpha                = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alphanumeric         = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	alphanumericExtended = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$_")
)

// IntFromRange returns random int within range from min to max
func IntFromRange(min, max int) int {
	return rand.Intn((max-min)+1) + min
}

// Numeric returns the random numbers string
func Numeric(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = numeric[rand.Intn(10)]
	}
	return string(b)
}

// Alpha returns the random alphabetic string
func Alpha(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = alpha[rand.Intn(52)]
	}
	return string(b)
}

// AlphaNumeric returns the random alpha numeric string
func AlphaNumeric(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(62)]
	}
	return string(b)
}

// Token creates a random string of given length
func Token(length int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$_")

	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// SecureToken creates a secure url friendly token
func SecureToken(length int) string {
	b := make([]byte, length)
	if _, err := io.ReadFull(cryptoRand.Reader, b); err != nil {
		panic(err.Error())
	}
	return removePadding(base64.URLEncoding.EncodeToString(b))
}

func removePadding(token string) string {
	return strings.TrimRight(token, "=")
}
