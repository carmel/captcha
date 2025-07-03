package captcha

import (
	"os"
	"strings"
	"testing"
)

func TestEngineCharCreate(t *testing.T) {
	tc, _ := os.MkdirTemp("", "audio")
	defer os.Remove(tc)
	for i := range 16 {
		configC.Mode = i % 4
		boooo := i%2 == 0
		configC.IsUseSimpleFont = boooo
		configC.IsShowSlimeLine = boooo
		configC.IsShowNoiseText = boooo
		configC.IsShowHollowLine = boooo
		configC.IsShowSineLine = boooo
		configC.IsShowNoiseDot = boooo

		im := EngineCharCreate(configC)
		fileName := strings.Trim(im.Content, "/+-+=?")
		err := WriteToFile(im, tc, fileName, "png")
		if err != nil {
			t.Error(err)
		}
	}
}
func TestMath(t *testing.T) {
	for range 100 {
		q, r := randArithmetic()
		t.Log(q, "--->", r)
	}
}
