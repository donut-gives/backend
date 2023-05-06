package utils

import "golang.org/x/crypto/bcrypt"

func Hash(plainText string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func Verify(plainText string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText))
	if err != nil {
		return false
	}
	return true
}
