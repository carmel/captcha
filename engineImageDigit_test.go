package captcha

import (
	"os"
	"testing"
)

func TestEngineDigitsCreate(t *testing.T) {
	td, _ := os.MkdirTemp("", "audio")
	defer os.Remove(td)
	for range 14 {
		idKey := randomId()
		im := EngineDigitsCreate(idKey, configD)
		err := WriteToFile(im, td, idKey, "png")
		if err != nil {
			t.Error(err)
		}
	}
}
