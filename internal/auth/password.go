package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type passwordVerifier struct {
}

func NewPasswordVerifier() *passwordVerifier {
	return &passwordVerifier{}
}

func (passwordVerifier *passwordVerifier) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
