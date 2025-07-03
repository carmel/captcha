package captcha

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func RandomStr(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for range l {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func StringUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
