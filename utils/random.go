package utils

import "math/rand"

const (
	alphabetLowercase = "abcdefghijklmnopqrstuvwxyz"
	alphabetUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numerals          = "0123456789"
)

func RandomString(length int) string {

	var letters = []byte(alphabetUppercase + alphabetLowercase + numerals)
	var result = make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
