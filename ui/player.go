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

package ui

import (
	"io"
	"log"
	"strconv"
	"time"

	"github.com/korandiz/mpa"

	"github.com/budkin/jam/music"
)

// OutputStream define an output stream
type OutputStream interface {
	CloseStream() error

	Write(data []byte) (int, error)
}

func (app *App) player() {
	stop := make(chan bool)
	pause := make(chan bool)
	playing := false
	paused := false
	next := false
	prev := false
	trackScrobbled := false
	pauseDur := time.Duration(0)
	var pauseTimer time.Time
	var songDur time.Duration

	stream, err := makeOutputStream()
	if err != nil {
		log.Fatalf("Can't playback from Pulse: %s", err)
	}
	defer stream.CloseStream()

	//var d mpa.Decoder
	var r *mpa.Reader
	data := make([]byte, 1024*8)

	//var buff [2][]float32
	for {
		switch <-app.Status.State {
		case 0:
			if paused {
				pause <- true
			}
			if playing {
				stop <- true
			}

			album := app.Status.NumAlbum[true]
			ntrack := app.Status.NumTrack
			queueTemp := make([][]*music.BTrack, len(app.Status.Queue))
			copy(queueTemp, app.Status.Queue)

			track := queueTemp[album][ntrack]

			song, err := app.GMusic.GetStream(track.ID)
			if err != nil {
				log.Fatalf("Can't play stream: %s", err)
			}

			defDur = time.Duration(0)
			defTrack = &music.BTrack{}
			app.printBar(defDur, defTrack)
			temp, _ := strconv.Atoi(track.DurationMillis)
			songDur = time.Duration(temp) * time.Millisecond

			//d = mpa.Decoder{Input: song.Body}
			r = &mpa.Reader{Decoder: &mpa.Decoder{Input: song.Body}}
			defer song.Body.Close()
			timer := time.Now()
			if app.Status.LastFM && app.LastFM != nil {
				go app.LastFM.NowPlaying(track.Title, track.Artist)
			}
			go func() {
				for {
					select {
					case <-pause:
						pauseTimer = time.Now()
						paused = true
					loop:
						for {
							select {
							case <-stop:
								pauseDur = time.Duration(0)
								paused = false
								return
							case <-pause:
								pauseDur = pauseDur + time.Since(pauseTimer)
								paused = false
								break loop
							}
						}
					case <-stop:
						pauseDur = time.Duration(0)
						return
					default:
						defer func() {
							if app.Status.LastFM && (defDur > songDur/2 ||
								defDur > 4*time.Minute) && !trackScrobbled &&
								songDur > 30*time.Second {
								trackScrobbled = true
								if app.LastFM != nil {
									go app.LastFM.Scrobble(track.Artist, track.Title, timer.Unix())
								}
							}
							trackScrobbled = false
							playing = false
							defDur = time.Duration(0)
							defTrack = &music.BTrack{}
							app.printBar(defDur, defTrack)
						}()
						playing = true

						defDur = time.Since(timer) - pauseDur
						defTrack = track
						app.printBar(defDur, defTrack)

						//buf := new(bytes.Buffer)

						i, err := r.Read(data)
						if err == io.EOF || i == 0 || next || prev {
							if next {
								next = false
							}

							if app.Status.LastFM && (defDur > songDur/2 ||
								defDur > 4*time.Minute) && !trackScrobbled &&
								songDur > 30*time.Second {
								trackScrobbled = true
								if app.LastFM != nil {
									go app.LastFM.Scrobble(track.Artist, track.Title, timer.Unix())
								}
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
							song, err = app.GMusic.GetStream(track.ID)
							if err != nil {
								log.Fatalf("Can't get stream: %s", err)
							}
							//d = mpa.Decoder{Input: song.Body}
							r = &mpa.Reader{Decoder: &mpa.Decoder{Input: song.Body}}
							pauseDur = time.Duration(0)
							defDur = time.Duration(0)
							defTrack = &music.BTrack{}
							app.printBar(defDur, defTrack)

							temp, _ := strconv.Atoi(track.DurationMillis)
							songDur = time.Duration(temp) * time.Millisecond
							timer = time.Now()
							trackScrobbled = false
							if app.Status.LastFM && app.LastFM != nil {
								go app.LastFM.NowPlaying(track.Title, track.Artist)
							}
							continue
						}

						i, err = stream.Write(data)
						if err != nil {
							log.Fatalf("Can't write stream: %s", err)
						}

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
