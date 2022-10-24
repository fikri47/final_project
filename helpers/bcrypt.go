package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(p string) string {
	salt := 8
	password := []byte(p)
	hash, err := bcrypt.GenerateFromPassword(password, salt)

	if err != nil {
		fmt.Println("Error hashing password:", err)
		return ""
	}

	return string(hash)
}

func ComparePass(h, p []byte) bool {
	hash, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err == nil
}
