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
	stream, _ := pulse.Playback("jam", "jam", &ss)
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
			//d = mpa.Decoder{Input: song.Body}
			r = &mpa.Reader{Decoder: &mpa.Decoder{Input: song.Body}}
			defer song.Body.Close()
			artist := <-curArtist
			checkErr(err)
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
						pauseDur = time.Duration(0)
						return
					default:
						defer func() {

							playing = false
							defDur = time.Duration(0)
							defTrack = &bTrack{}
							defArtist = ""
						}()
						playing = true

						defDur = time.Since(timer) + pauseDur
						defTrack = track
						defArtist = artist
						printBar(s, defDur, defTrack, defArtist)

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
							defArtist = ""
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
