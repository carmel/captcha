package captcha

import (
	"strings"

	"github.com/google/uuid"
)

func StringUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
