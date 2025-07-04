// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"io"
	"os"
	"testing"
)

func BenchmarkNewAudio(b *testing.B) {
	b.StopTimer()
	d := randomDigits(DefaultLen)
	id := StringUUID()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		newAudio(id, d, "")
	}
}

func BenchmarkAudioWriteTo(b *testing.B) {
	b.StopTimer()
	d := randomDigits(DefaultLen)
	id := StringUUID()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := newAudio(id, d, "")
		n, _ := a.WriteTo(io.Discard)
		b.SetBytes(n)
	}
}

func TestEngineAudioCreate(t *testing.T) {
	ta, _ := os.MkdirTemp("", "audio")
	defer os.RemoveAll(ta)
	for i := 0; i < 10; i++ {
		idKey := StringUUID()
		au := EngineAudioCreate(idKey, configA)
		if err := WriteToFile(au, ta, idKey, "wav"); err != nil {
			t.Log(err)
		}
	}

}
