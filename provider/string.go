package provider

import (
	"math/rand"
	"time"
)

type String struct {
}

func RandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	count := len(letterRunes)
	for i := range b {
		b[i] = letterRunes[r.Int63n(int64(count))]
	}
	return string(b)
}

func RandomInt(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	count := len(letterRunes)
	for i := range b {
		b[i] = letterRunes[r.Int63n(int64(count))]
	}
	return string(b)
}
