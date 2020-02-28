package random

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"io"
	"math/rand"
	"strings"
)

// IntFromRange returns random int within range from min to max
func IntFromRange(min, max int) int {
	return rand.Intn((max-min)+1) + min
}

// Token creates a random string of given length
func Token(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$_")

	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// SecureToken creates a secure token of given length
func SecureToken(n int) string {
	b := make([]byte, n)
	if _, err := io.ReadFull(cryptoRand.Reader, b); err != nil {
		panic(err.Error())
	}
	return removePadding(base64.URLEncoding.EncodeToString(b))
}

func removePadding(token string) string {
	return strings.TrimRight(token, "=")
}
