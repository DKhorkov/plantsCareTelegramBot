package utils

import (
	"math/rand"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Глобальный источник случайности, инициализируемый один раз.
	src = rand.NewSource(time.Now().UnixNano())
	rng = rand.New(src)
)

// The RandSequence function that generates a random string sequence.
func RandSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rng.Intn(len(letterRunes))]
	}

	return string(b)
}
