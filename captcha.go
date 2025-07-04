package captcha

import (
	"image/color"
	"strings"
	"time"
)

// Captcha represents a captcha service.
type Captcha struct {
	store           *MapStore
	Width           int     `yaml:"width"`
	Height          int     `yaml:"height"`
	Length          int     `yaml:"length"`
	Noises          int     `yaml:"noises"`
	Expiration      int     `yaml:"expiration"`       // minutes
	CleanupInterval int     `yaml:"cleanup-interval"` // minutes
	MaxSkew         float64 `yaml:"max-skew"`
	ColorPalette    color.Palette
}

func (c *Captcha) Init() {
	if c.store == nil {
		c.store = NewMapStore(time.Duration(c.Expiration)*time.Minute, time.Duration(c.CleanupInterval)*time.Minute)
	}
}

// create a new captcha id
func (c *Captcha) Generate() (string, string) {
	id := StringUUID()
	chars := RandomStr(c.Length)
	c.store.Set(id, chars)
	img := Image{
		Chars:   chars,
		Width:   c.Width,
		Height:  c.Height,
		Noises:  c.Noises,
		MaxSkew: c.MaxSkew,
	}
	return id, img.Base64()
}

// verify from a request
// func (c *Captcha) VerifyReq(req *http.Request) bool {
// 	_ = req.ParseForm()
// 	return c.Verify(req.Form.Get(c.FieldIdName), req.Form.Get(c.FieldCaptchaName))
// }

// direct verify id and challenge string
func (c *Captcha) Verify(id string, challenge string) bool {
	if len(challenge) == 0 || len(id) == 0 {
		return false
	}

	val := c.store.Get(id)
	if val == nil {
		return false
	}
	return strings.EqualFold(challenge, val.(string))
}
