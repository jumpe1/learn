package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

func jqueryNow() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func hexSha256(s string) string {
	h := sha256.Sum256([]byte(s))
	return strings.ToUpper(hex.EncodeToString(h[:]))
}

func passwordAlgorithmsCookie(s string) string {
	return hexSha256(s)
}
