package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordGenerator struct {
}

func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

func (passwordGenerator *PasswordGenerator) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
