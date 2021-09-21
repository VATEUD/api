package utils

import "math/rand"

func RandomString(length int) string {

	var letters = []byte(alphabet + numerals)
	var result = make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
