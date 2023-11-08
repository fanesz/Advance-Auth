package utils

import (
	"math/rand"
	"time"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers          = "0123456789"
)

func GenerateRandomString(length int, isCapital, isWithNumbers bool) string {
	rand.Seed(time.Now().UnixNano())
	var charSet string
	if isCapital {
		charSet += uppercaseLetters
	}
	charSet += lowercaseLetters
	if isWithNumbers {
		charSet += numbers
	}
	result := make([]byte, length)
	for i := range result {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}
