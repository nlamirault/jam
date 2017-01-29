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

package main

import (
	"io"
	"time"

	"github.com/gdamore/tcell"
	"github.com/korandiz/mpa"

	pulse "github.com/mesilliac/pulse-simple"
)

func player(s tcell.Screen) {
	stop := make(chan bool)
	pause := make(chan bool)
	playing := false
	paused := false
	next := false
	prev := false
	pauseDur := time.Duration(0)

	ss := pulse.SampleSpec{pulse.SAMPLE_S16LE, 44100, 2}
	stream, err := pulse.Playback("jam", "jam", &ss)
	checkErr(err)
	defer stream.Free()
	defer stream.Drain()

	//var d mpa.Decoder
	var r *mpa.Reader
	data := make([]byte, 1024*8)

	//var buff [2][]float32
	for {
		switch <-state {
		case 0:
			if paused {
				pause <- true
			}
			if playing {
				stop <- true
			}

			album := numAlbum[true]
			ntrack := numTrack
			queueTemp := make([][]*bTrack, len(queue))
			copy(queueTemp, queue)

			track := queueTemp[album][ntrack]
			song, err := gm.GetStream(track.ID)
			checkErr(err)
			//d = mpa.Decoder{Input: song.Body}
			r = &mpa.Reader{Decoder: &mpa.Decoder{Input: song.Body}}
			defer song.Body.Close()
			timer := time.Now()
			go func() {
				for {
					select {
					case <-pause:
						pauseDur = defDur
						paused = true
						for {
							if <-pause {
								timer = time.Now()
								paused = false
								break
							}
						}
					case <-stop:
						playing = false
						pauseDur = time.Duration(0)
						return
					default:
						defer func() {
							playing = false
							defDur = time.Duration(0)
							defTrack = &bTrack{}
						}()
						playing = true

						defDur = time.Since(timer) + pauseDur
						defTrack = track
						printBar(s, defDur, defTrack)

						//buf := new(bytes.Buffer)

						i, err := r.Read(data)
						if err == io.EOF || i == 0 || next || prev {
							if next {
								next = false
							}

							switch err.(type) {
							case mpa.MalformedStream:
								continue
							}
							if !prev {
								if ntrack < len(queueTemp[album])-1 {
									ntrack++
								} else if album < len(queueTemp)-1 {
									album++
									ntrack = 0
								} else {
									return
								}
							} else {
								if ntrack > 0 {
									ntrack--
								} else if album > 0 && len(queueTemp[album-1]) > 0 {
									album--
									ntrack = len(queueTemp[album]) - 1
								} else {
									ntrack = 0
								}
								prev = false
							}

							track = queueTemp[album][ntrack]
							song, err = gm.GetStream(track.ID)
							checkErr(err)
							//d = mpa.Decoder{Input: song.Body}
							r = &mpa.Reader{Decoder: &mpa.Decoder{Input: song.Body}}
							checkErr(err)
							pauseDur = time.Duration(0)
							defDur = time.Duration(0)
							defTrack = &bTrack{}
							updateUI(s)

							timer = time.Now()
							continue
						}

						i, err = stream.Write(data)
						checkErr(err)

					}
				}
			}()
		case 1:
			if playing {
				stop <- true
			}
		case 2:
			if playing {
				pause <- true
			}
		case 3:
			if playing {
				next = true
			}
			if paused {
				pause <- true
				next = true
			}
		case 4:
			if playing {
				prev = true
			}
			if paused {
				pause <- true
				prev = true
			}
		}
	}
}
