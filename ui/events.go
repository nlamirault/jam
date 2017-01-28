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

func (app *App) numAlb(k int) int {
	var i int
	for app.Artists[k] == "" {
		i++
		k--
	}
	return i

}
func (app *App) numSongs() int {
	j := 0
	for _, sngs := range app.Songs {
		j++
		for _ = range sngs {
			j++
		}
	}
	return j
}

func (app *App) pageUp() {
}

func (app *App) pageDn() {
}

func (app *App) scrollDown() {
	var length int
	if !app.Status.InTracks {
		length = len(app.Artists)
		app.Status.NumAlbum[false] = -1
	} else if app.Status.NumAlbum[false] == -1 {
		length = app.numSongs()
		app.Status.NumAlbum[true] = len(app.Songs) - 1
		app.Status.NumTrack = len(app.Songs[app.LastAlbum]) - 1
	} else {
		length = len(app.Songs[app.Albums[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]-1*(app.Status.NumAlbum[false]+1)]][app.Status.NumAlbum[true]]])
		app.Status.NumTrack = length - 1
	}
loop:
	for app.Status.CurPos[app.Status.InTracks]+app.Status.ScrOffset[app.Status.InTracks] < length {
		for app.Status.CurPos[app.Status.InTracks] < app.Height-3 {
			if app.Status.CurPos[app.Status.InTracks] == length {
				break loop
			}
			app.Status.CurPos[app.Status.InTracks]++
		}
		app.Status.ScrOffset[app.Status.InTracks]++
	}

}
func (app *App) scrollUp() {
	if !app.Status.InTracks {
		app.Status.CurPos[app.Status.InTracks] = 1
		app.Status.ScrOffset[app.Status.InTracks] = 0
		app.Status.NumAlbum[false] = -1
	} else if app.Status.NumAlbum[false] == -1 {
		app.Status.CurPos[app.Status.InTracks] = 2
		app.Status.ScrOffset[app.Status.InTracks] = 0
		app.Status.NumTrack = 0
		app.Status.NumAlbum[true] = 0
	} else {
		app.Status.CurPos[app.Status.InTracks] = 1
		app.Status.ScrOffset[app.Status.InTracks] = 0
		app.Status.NumTrack = 0
	}
}

func (app *App) toggleView() {
	app.Status.InTracks = !app.Status.InTracks

}

func (app *App) toggleAlbums() {
	app.Status.NumAlbum[false] = -1
	i := app.Status.CurPos[false] - 1 + app.Status.ScrOffset[false]
	app.ArtistsMap[app.Artists[i-app.numAlb(i)]] = !app.ArtistsMap[app.Artists[i-app.numAlb(i)]]

	if app.ArtistsMap[app.Artists[i-app.numAlb(i)]] {
		app.appendAlbums()
	} else {
		app.removeAlbums()
		app.Status.NumAlbum[true] = 0
		app.Status.NumTrack = 0
		app.Status.ScrOffset[true] = 0
		app.Status.CurPos[true] = 2
	}
}
func (app *App) appendAlbums() {
	var b []string
	boo := app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]] == app.Artists[len(app.Artists)-1]
	a := make([]string, len(app.Albums[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]]]))
	if !boo {
		b = make([]string, len(app.Artists[app.Status.CurPos[false]+app.Status.ScrOffset[false]:]))
		copy(b, app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]+1:])
	}
	app.Artists = append(app.Artists[:app.Status.CurPos[false]-1+app.Status.ScrOffset[false]+1], a...)
	if !boo {
		app.Artists = append(app.Artists, b...)
	}
}

func (app *App) removeAlbums() {
	i := app.Status.CurPos[false] - 1 + app.Status.ScrOffset[false]
	if app.Artists[i-app.numAlb(i)] != app.Artists[len(app.Artists)-1-len(app.Albums[app.Artists[i-app.numAlb(i)]])] {
		app.Status.CurPos[false] -= app.numAlb(i)
		app.Artists = append(app.Artists[:i-app.numAlb(i)+1], app.Artists[i-app.numAlb(i)+len(app.Albums[app.Artists[i-app.numAlb(i)]])+1:]...)
	} else {
		a := app.numAlb(i)
		app.Artists = app.Artists[:i-app.numAlb(i)+1]
		app.Status.ScrOffset[false] -= a
	}

}

func (app *App) upEntry() {
	switch app.Status.InTracks {
	case false:
		app.Status.NumTrack = 0
		app.Status.CurPos[true] = 2
		app.Status.ScrOffset[true] = 0
	case true:
		if app.Status.NumAlbum[false] == -1 {
			if app.Status.NumTrack > 0 {
				app.Status.NumTrack--
			} else if app.Status.NumTrack == 0 && app.Status.NumAlbum[true] > 0 {
				app.Status.NumAlbum[true]--
				if app.Status.CurPos[true] > 3 {
					app.Status.CurPos[true]--
				} else {
					app.Status.ScrOffset[true]--
				}
				songs := app.Songs[app.Albums[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]]][app.Status.NumAlbum[true]]]
				app.Status.NumTrack = len(songs) - 1
			} else if app.Status.CurPos[true] < 3 {
				return
			}
		} else {
			if app.Status.NumTrack > 0 {
				app.Status.NumTrack--
			}
		}
	}

	if app.Status.CurPos[app.Status.InTracks] > 3 {
		app.Status.CurPos[app.Status.InTracks]--
	} else if app.Status.ScrOffset[app.Status.InTracks]+app.Status.CurPos[app.Status.InTracks] > 3 {
		app.Status.ScrOffset[app.Status.InTracks]--
	} else if app.Status.ScrOffset[app.Status.InTracks]+app.Status.CurPos[app.Status.InTracks] > 1 {
		app.Status.CurPos[app.Status.InTracks]--
	}

	if !app.Status.InTracks {
		if app.ArtistsMap[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]]] {
			app.Status.NumAlbum[false] = -1
			app.Status.NumAlbum[true] = 0
		} else if app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]] == "" {
			curPosTemp := app.Status.CurPos[false] - 1
			for app.Artists[curPosTemp+app.Status.ScrOffset[false]] == "" {
				curPosTemp--
			}
			app.Status.NumAlbum[false] = app.Status.CurPos[false] - curPosTemp - 1 - 1
			if app.Status.NumAlbum[false] > -1 {

				app.Status.NumAlbum[true] = app.Status.NumAlbum[false]
			}
			app.Status.CurPos[true] = 1

		}
	}

}

func (app *App) downEntry() {
	var high int
	switch app.Status.InTracks {
	case true:
		high = 0
		if app.Status.NumAlbum[false] == -1 {
			for _, sngs := range app.Songs {
				high++
				for _ = range sngs {
					high++
				}
			}
		} else {
			high = len(app.Songs[app.Albums[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]-1*(app.Status.NumAlbum[false]+1)]][app.Status.NumAlbum[true]]])
		}

		if app.Status.NumTrack < len(app.Songs[app.Albums[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]-1*(app.Status.NumAlbum[false]+1)]][app.Status.NumAlbum[true]]])-1 {
			app.Status.NumTrack++
		} else if app.Status.NumAlbum[true] < len(app.Albums[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]-1*(app.Status.NumAlbum[false]+1)]])-1 && app.Status.NumAlbum[false] == -1 {
			app.Status.NumAlbum[true]++
			if app.Status.CurPos[true] > app.Height-6 {
				app.Status.ScrOffset[true]++
			} else {
				app.Status.CurPos[true]++
			}
			app.Status.NumTrack = 0
		}
	case false:
		high = len(app.Artists)
		if app.ArtistsMap[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]]] ||
			app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]] == "" {
			if app.Status.NumAlbum[false] < len(app.Albums[app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]-1*(app.Status.NumAlbum[false]+1)]])-1 {
				app.Status.NumAlbum[false]++
				app.Status.NumAlbum[true] = app.Status.NumAlbum[false]

				app.Status.CurPos[true] = 1
			}

		} else {
			app.Status.NumAlbum[false] = -1
			app.Status.NumAlbum[true] = 0
		}

		app.Status.NumTrack = 0
		app.Status.CurPos[true] = 2
		app.Status.ScrOffset[true] = 0

	}

	if app.Status.CurPos[app.Status.InTracks]+app.Status.ScrOffset[app.Status.InTracks] < high {
		if app.Status.CurPos[app.Status.InTracks] > app.Height-6 && app.Status.CurPos[app.Status.InTracks]+app.Status.ScrOffset[app.Status.InTracks] < high-2 {
			app.Status.ScrOffset[app.Status.InTracks]++
		} else if app.Status.CurPos[app.Status.InTracks] == app.Height-3 && app.Status.CurPos[app.Status.InTracks]+app.Status.ScrOffset[app.Status.InTracks] < high {
			app.Status.ScrOffset[app.Status.InTracks]++
		} else {
			app.Status.CurPos[app.Status.InTracks]++
		}
	}
	if !app.Status.InTracks {
		if app.Artists[app.Status.CurPos[false]-1+app.Status.ScrOffset[false]] != "" {
			app.Status.NumAlbum[false] = -1
		} else {
			app.Status.CurPos[true] = 1
		}
	}
}
