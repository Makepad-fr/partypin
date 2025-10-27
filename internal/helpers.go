package internal

import (
	"math/rand"
	"time"
)

func GenerateEventID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seeded := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seeded.Intn(len(charset))]
	}
	return string(b)
}
