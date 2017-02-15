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

// +build linux

package ui

import (
	"fmt"

	pulse "github.com/mesilliac/pulse-simple"
)

type linuxOutputStream struct {
	Stream *pulse.Stream
}

func makeOutputStream() (OutputStream, error) {
	ss := pulse.SampleSpec{pulse.SAMPLE_S16LE, 44100, 2}
	stream, err := pulse.Playback("jam", "jam", &ss)
	if err != nil {
		return nil, fmt.Errorf("Can't playback from Pulse: %s", err)
	}
	return &linuxOutputStream{
		Stream: stream,
	}, nil
}

func (los *linuxOutputStream) CloseStream() error {
	los.Stream.Free()
	if _, err := los.Stream.Drain(); err != nil {
		return err
	}
	return nil
}

func (los *linuxOutputStream) Write(data []byte) (int, error) {
	return los.Stream.Write(data)
}
