package forum

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CompareHashPassword(hash []byte, password string) error {
        err := bcrypt.CompareHashAndPassword(hash, []byte(password))
        return err
}
