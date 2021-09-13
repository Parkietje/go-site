package main

import (
    "golang.org/x/crypto/sha3"
	"encoding/hex"
)

func hash(password string) string{
    h := sha3.New512()
    h.Write([]byte(password))
    sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
