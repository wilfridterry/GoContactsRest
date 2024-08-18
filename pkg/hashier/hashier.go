package hashier

import (
	"crypto/sha256"
	"fmt"
)

type Hashier struct {
	salt string
}

func (h *Hashier) Hash(password string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", nil
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}

func NewHashier(salt string) *Hashier {
	return &Hashier{salt}
} 