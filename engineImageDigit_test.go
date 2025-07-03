package captcha

import (
	"os"
	"testing"
)

func TestEngineDigitsCreate(t *testing.T) {
	td, _ := os.MkdirTemp("", "audio")
	defer os.Remove(td)
	for i := 0; i < 14; i++ {
		idKey := randomId()
		im := EngineDigitsCreate(idKey, configD)
		err := CaptchaWriteToFile(im, td, idKey, "png")
		if err != nil {
			t.Error(err)
		}
	}
}
