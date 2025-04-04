package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	cost int
}

func New(cost int) *Hasher {
	return &Hasher{cost: cost}
}

func (h Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h Hasher) Verify(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
