package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func generateToken(length int) string {
    data := make([]byte, length)
    rand.Read(data)
    return base64.URLEncoding.EncodeToString(data)
}
