package common

import (
	"crypto/sha1"
	"fmt"
)

type SHA1HashCalculator struct{}

func Calculate(data string) string {
	hash := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}
