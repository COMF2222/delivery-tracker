package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateTrackNumber() (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	length := 12

	result := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))

		if err != nil {
			return "", fmt.Errorf("failed to generate track number: %w", err)
		}

		randomIndex := n.Int64()

		result[i] = letters[int(randomIndex)]
	}
	return string(result), nil
}
