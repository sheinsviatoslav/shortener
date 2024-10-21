package hash

import (
	"math/rand"
	"strings"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Generator(n int) string {
	hashSlice := make([]string, n)
	for i := range hashSlice {
		hashSlice[i] = string(chars[rand.Intn(len(chars))])
	}

	return strings.Join(hashSlice, "")
}
