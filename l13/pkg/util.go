package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
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

type Logger struct {
	timestampString string
}

func NewLogger() *Logger {
	return &Logger{
		timestampString: time.Now().Format(time.RFC3339),
	}
}

func (l *Logger) LogAndWriteStatus(filename, message string) {
	fullMessage := fmt.Sprintf("%s - %s", l.timestampString, message)

	log.Println(fullMessage)

	if filename != "" {
		if err := os.WriteFile(filename, []byte(fullMessage+"\n"), 0644); err != nil {
			log.Printf("Failed to write to status file (%s): %v", filename, err)
		}
	}
}
