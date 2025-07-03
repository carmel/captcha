// Copyright 2014 The Macaron Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package captcha

import (
	"fmt"
	"testing"
)

func TestCaptcha(t *testing.T) {
	captcha := Captcha{Height: 35, Width: 120, Length: 5, Expiration: 1, CleanupInterval: 30, Noises: 20, MaxSkew: 0.7}
	captcha.Init()
	captcha.Generate()

	img := Image{
		Chars:   "aB1xY9",
		Width:   120,
		Height:  35,
		Noises:  20,
		MaxSkew: 0.3,
	}
	fmt.Println(img.Base64())
}
