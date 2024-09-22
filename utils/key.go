package utils

import "math/rand"

func Key(limit int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	key := make([]rune, limit)
	for i := range key {
		key[i] = chars[rand.Intn(len(chars))]
	}

	return string(key)
}
