package authentication

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"
)

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		zap.S().Errorf("Bad password: %v", err)
		return "", err
	}

	return fmt.Sprintf("%s", hash), err
}
