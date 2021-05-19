package random

import (
	"crypto/rand"
	"github.com/go-basic/uuid"
	"strings"
)

func Uuid(separator bool) string {
	uuid := uuid.New()
	if separator {
		return uuid
	}
	return strings.ReplaceAll(uuid, "-", "")
}

func Bytes(length int) []byte {
	tmp := make([]byte, length)
	_, err := rand.Read(tmp)
	if err != nil {
		return nil
	}
	return tmp
}
