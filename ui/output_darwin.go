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

// +build darwin

package ui

import (
	// "fmt"

	"github.com/gordonklaus/portaudio"
)

const (
	bufferSize = 512
)

type darwinOutputStream struct {
	Stream *portaudio.Stream
	buffer []int32
}

func makeOutputStream() (OutputStream, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}
	return &darwinOutputStream{
		buffer: make([]int32, bufferSize),
		Stream: nil,
	}, nil
}

func (dos *darwinOutputStream) CloseStream() error {
	if err := dos.Stream.Close(); err != nil {
		return err
	}
	return dos.Stream.Stop()
}

func (dos *darwinOutputStream) Write(data []byte) (int, error) {
	var err error
	dos.Stream, err = portaudio.OpenDefaultStream(
		0, 1, 44100, len(dos.buffer), &dos.buffer)
	if err != nil {
		return 0, err
	}

	if err = dos.Stream.Start(); err != nil {
		return 0, err
	}

	dos.Stream.Write()
	return 0, nil
}
