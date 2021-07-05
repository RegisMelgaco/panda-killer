package algorithms

import "golang.org/x/crypto/bcrypt"

type AccountSecurityAlgorithmsImpl struct{}

func (a AccountSecurityAlgorithmsImpl) GenerateSecretFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}

func (a AccountSecurityAlgorithmsImpl) CheckSecretAndPassword(secret, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(secret), []byte(password))
}
