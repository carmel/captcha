package captcha

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var configD = ConfigDigit{
	Height:     80,
	Width:      240,
	MaxSkew:    0.7,
	DotCount:   80,
	CaptchaLen: 5,
}

var configA = ConfigAudio{
	CaptchaLen: 6,
	Language:   "zh",
}

var configC = ConfigCharacter{
	Height:             60,
	Width:              240,
	Mode:               0,
	ComplexOfNoiseText: 0,
	ComplexOfNoiseDot:  0,
	IsUseSimpleFont:    false,
	IsShowHollowLine:   false,
	IsShowNoiseDot:     false,
	IsShowNoiseText:    false,
	IsShowSlimeLine:    false,
	IsShowSineLine:     false,
	CaptchaLen:         6,
}

func TestGenerateCaptcha(t *testing.T) {
	testDir, _ := os.MkdirTemp("", "")
	defer os.Remove(testDir)

	for idx, vv := range []any{configA, configD} {

		idkey, cap := Generate(vv)
		ext := "png"
		if idx == 0 {
			ext = "wav"
		}

		WriteToFile(cap, testDir, idkey, ext)
		WriteToFile(cap, testDir, idkey, ext)

		WriteToFile(cap, testDir, idkey, ext)

		// t.Log(idkey, globalStore.Get(idkey, false))

	}
	testDirAll, _ := os.MkdirTemp("", "all")
	defer os.RemoveAll(testDirAll)
	for i := range 16 {
		configC.Mode = i % 4
		idkey, cap := Generate(configC)
		ext := "png"
		err := WriteToFile(cap, testDirAll, "char_"+idkey, ext)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCaptchaWriteToBase64Encoding(t *testing.T) {
	_, cap := Generate(configD)
	base64string := WriteToBase64Encoding(cap)
	if !strings.Contains(base64string, MimeTypeCaptchaImage) {

		t.Error("encodeing base64 string failed.")
	}
	_, capA := Generate(configA)
	base64stringA := WriteToBase64Encoding(capA)
	if !strings.Contains(base64stringA, MimeTypeCaptchaAudio) {

		t.Error("encodeing base64 string failed.")
	}
}

func TestVerifyCaptcha(t *testing.T) {
	idkey, _ := Generate(configD)
	verifyValue := globalStore.Get(idkey)
	if Verify(idkey, verifyValue.(string)) {
		t.Log(idkey, verifyValue)
	} else {
		t.Error("verify captcha content is failed.")
	}

	Verify("", "")
	Verify("dsafasf", "ddd")

}

func TestPathExists(t *testing.T) {

	testDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(testDir)
	assert.True(t, pathExists(testDir))
	assert.False(t, pathExists(testDir+"/NotExistFolder"))
}

func TestCaptchaWriteToFileCreateDirectory(t *testing.T) {

	idKey, captcha := Generate(configD)
	testDir, _ := os.MkdirTemp("", "")
	defer os.Remove(testDir)
	assert.Nil(t, WriteToFile(captcha, testDir+"/NotExistFolder", idKey, "png"))
}

func TestCaptchaWriteToFileCreateFileFailed(t *testing.T) {

	var err error
	idKey, captcha := Generate(configD)
	testDir, _ := os.MkdirTemp("", "")
	defer os.Remove(testDir)
	noPermissionDirPath := testDir + "/NoPermission"

	err = os.Mkdir(noPermissionDirPath, os.ModeDir)
	assert.Nil(t, err)

	err = WriteToFile(captcha, noPermissionDirPath, idKey, "png")
	//has no permission must failed
	if runtime.GOOS == "windows" {
		assert.Nil(t, err)
	} else {
		assert.NotNil(t, err)
	}
}

func TestSetCustomStore(t *testing.T) {
	globalStore = NewMapStore(1*time.Second, 5*time.Second)
	verifyCode := "sddsffds"
	globalStore.Set("1", verifyCode)

	assert.Equal(t, globalStore.Get("1"), verifyCode)
}
