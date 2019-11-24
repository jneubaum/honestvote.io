package random

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	math_rand "math/rand"
)

func init() {
	var b [16]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := math_rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
