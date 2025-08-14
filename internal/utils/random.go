package utils

import (
	"math/rand"
	"time"
)

// The RandSequence function that generates a random string sequence.
func RandSequence(n int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
