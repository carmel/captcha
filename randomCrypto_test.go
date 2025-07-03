package captcha

import (
	"testing"
)

func TestRandomId(t *testing.T) {
	for range 10 {
		t.Log(randomDigits(5))
	}
}
func TestRandomDigits(t *testing.T) {
	for range 10 {
		t.Log(randomDigits(5))
	}
}
func TestParseDigitsToString(t *testing.T) {
	for range 10 {
		byss := randomDigits(5)
		t.Log(byss)
		bsssstring := parseDigitsToString(byss)
		t.Log(bsssstring)
	}
}
