package pkg

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncodeToHexString(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashByteSlice := hash.Sum(nil)
	hashString := hex.EncodeToString(hashByteSlice)

	return hashString
}
