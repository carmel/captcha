// Copyright 2017 Eric Zhou. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package captcha supports digits, numbers,alphabet, arithmetic, audio and digit-alphabet captcha.
// captcha is used for fast development of RESTful APIs, web apps and backend services in Go. give a string identifier to the package and it returns with a base64-encoding-png-string
package captcha

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	globalStore = NewMapStore(1*time.Minute, 10*time.Minute)
)

// CaptchaInterface captcha interface for captcha engine to to write staff
type CaptchaInterface interface {
	// BinaryEncodeing covert to bytes
	BinaryEncodeing() []byte
	// WriteTo output captcha entity
	WriteTo(w io.Writer) (n int64, err error)
}

// WriteToBase64Encoding converts captcha to base64 encoding string.
// mimeType is one of "audio/wav" "image/png".
func WriteToBase64Encoding(cap CaptchaInterface) string {
	binaryData := cap.BinaryEncodeing()
	var mimeType string
	if _, ok := cap.(*Audio); ok {
		mimeType = MimeTypeCaptchaAudio
	} else {
		mimeType = MimeTypeCaptchaImage
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(binaryData))
}

// WriteToFile output captcha to file.
// fileExt is one of "png","wav"
func WriteToFile(cap CaptchaInterface, outputDir, fileName, fileExt string) error {
	filePath := filepath.Join(outputDir, fileName+"."+fileExt)
	if !pathExists(outputDir) {
		_ = os.MkdirAll(outputDir, os.ModePerm)
	}
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("%s is invalid path.error:%v", filePath, err)
		return err
	}
	defer file.Close()
	_, err = cap.WriteTo(file)
	return err
}

// CaptchaItem captcha basic information.
type CaptchaItem struct {
	// Content captcha entity content.
	Content string
	// VerifyValue captcha verify value.
	VerifyValue string
	// ImageWidth image width pixel.
	ImageWidth int
	// ImageHeight image height pixel.
	ImageHeight int
}

// Verify by given id key and remove the captcha value in store, return boolean value.
// 验证图像验证码,返回boolean.
func Verify(identifier, verifyValue string) bool {
	// return VerifyAndIsClear(identifier, verifyValue, true)

	if verifyValue == "" {
		return false
	}
	storeValue := globalStore.Get(identifier)
	if storeValue == nil {
		return false
	}
	return strings.EqualFold(storeValue.(string), verifyValue)
}

// VerifyAndIsClear verify captcha, return boolean value.
// identifier is the captcha id,
// verifyValue is the captcha image value,
// isClear is whether to clear the value in store.
// 验证图像验证码,返回boolean.
// func VerifyAndIsClear(identifier, verifyValue string, isClear bool) bool {
// 	if verifyValue == "" {
// 		return false
// 	}
// 	storeValue := globalStore.Get(identifier, isClear)
// 	return strings.EqualFold(storeValue, verifyValue)
// }

// Generate create captcha by config struct and id.
// idkey can be an empty string, base64 will create a unique id four you.
// if idKey is a empty string, the package will generate a random unique identifier for you.
// configuration struct should be one of those struct ConfigAudio, ConfigCharacter, ConfigDigit.
//
// Example Code
//
//	//config struct for digits
//	var configD = captcha.ConfigDigit{
//		Height:     80,
//		Width:      240,
//		MaxSkew:    0.7,
//		DotCount:   80,
//		CaptchaLen: 5,
//	}
//	//config struct for audio
//	var configA = captcha.ConfigAudio{
//		CaptchaLen: 6,
//		Language:   "zh",
//	}
//	//config struct for Character
//	var configC = captcha.ConfigCharacter{
//		Height:             60,
//		Width:              240,
//		//const ModeNumber:数字,ModeAlphabet:字母,ModeArithmetic:算术,ModeNumberAlphabet:数字字母混合.
//		Mode:               captcha.ModeNumber,
//		ComplexOfNoiseText: captcha.ComplexLower,
//		ComplexOfNoiseDot:  captcha.ComplexLower,
//		IsUseSimpleFont:    true,
//		IsShowHollowLine:   false,
//		IsShowNoiseDot:     false,
//		IsShowNoiseText:    false,
//		IsShowSlimeLine:    false,
//		IsShowSineLine:     false,
//		CaptchaLen:         6,
//	}
//	//create a audio captcha.
//	//Generate first parameter is empty string,so the package will generate a random uuid for you.
//	idKeyA,capA := captcha.Generate("",configA)
//	//write to base64 string.
//	//Generate first parameter is empty string,so the package will generate a random uuid for you.
//	base64stringA := captcha.WriteToBase64Encoding(capA)
//	//create a characters captcha.
//	//Generate first parameter is empty string,so the package will generate a random uuid for you.
//	idKeyC,capC := captcha.Generate("",configC)
//	//write to base64 string.
//	base64stringC := captcha.WriteToBase64Encoding(capC)
//	//create a digits captcha.
//	idKeyD,capD := captcha.Generate("",configD)
//	//write to base64 string.
//	base64stringD := captcha.WriteToBase64Encoding(capD)
func Generate(configuration any) (id string, captchaInstance CaptchaInterface) {
	id = StringUUID()
	var verifyValue string
	switch config := configuration.(type) {
	case ConfigAudio:
		audio := EngineAudioCreate(id, config)
		verifyValue = audio.VerifyValue
		captchaInstance = audio

	case ConfigCharacter:
		char := EngineCharCreate(config)
		verifyValue = char.VerifyValue
		captchaInstance = char

	case ConfigDigit:
		dig := EngineDigitsCreate(id, config)
		verifyValue = dig.VerifyValue
		captchaInstance = dig

	default:
		log.Fatal("config type not supported", config)
	}

	globalStore.Set(id, verifyValue)

	return id, captchaInstance
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
