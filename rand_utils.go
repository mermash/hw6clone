package main

import (
	"fmt"
	"math/rand"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// from somewhere
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandBytesHex(n int) string {
	return fmt.Sprintf("%x", RandBytes(n))
}

func RandBytes(n int) []byte {
	res := make([]byte, n)
	rand.Read(res)
	return res
}
