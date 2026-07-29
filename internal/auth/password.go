package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type passwordGenerator struct {
}

func NewPasswordGenerator() *passwordGenerator {
	return &passwordGenerator{}
}

func (passwordGenerator *passwordGenerator) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
