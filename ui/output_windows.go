// Copyright (c) 2016, 2017 Evgeny Badin

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// +build windows

package ui

import (
	"fmt"

	"github.com/koron/go-waveout"
)

type windowsOutputStream struct {
	Player *waveout.Player
}

func makeOutputStream() (OutputStream, error) {
	player, err := waveout.NewWithBuffers(2, 44100, 16, 8, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to create waveout: %s", err)
	}
	return &windowsOutputStream{
		Player: player,
	}, nil
}

func (wos *windowsOutputStream) CloseStream() error {
	return wos.Player.Close()
}

func (wos *windowsOutputStream) Write(data []byte) (int, error) {
	return wos.Player.Write(data)
}
