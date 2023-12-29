package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
		// return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}


func VerifyPassword(hashedPassword string, EnteredPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(EnteredPassword))
}

