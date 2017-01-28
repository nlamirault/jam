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
	// "encoding/json"
	// "fmt"
	"log"
	"math/rand"
	"sort"
	// "strconv"
	"strings"

	// "time"

	"github.com/boltdb/bolt"
	"github.com/budkin/gmusic"
	"github.com/gdamore/tcell"
	// runewidth "github.com/mattn/go-runewidth"

	"github.com/budkin/jam/music"
)

const (
	play = iota
	stop
	pause
	next
	prev
)

// type Database struct {
// 	DB         *bolt.DB
// 	ArtistsMap map[string]bool
// 	Artists    sort.StringSlice
// 	Songs      map[string][]string
// 	Albums     map[string][]string
// 	LastAlbum  string
// }

type Status struct {
	ScrOffset map[bool]int
	Offset    int
	CurPos    map[bool]int
	NumAlbum  map[bool]int
	InTracks  bool
	InSearch  bool
	NumTrack  int
	CurArtist chan string
	Queue     [][]*music.BTrack // playlist, updated on each movement of cursor in artists view
	Query     []rune            // search query

	State chan int // player's state: play, pause, stop, etc

}

// App define the UI application
type App struct {
	Screen tcell.Screen
	Width  int
	Height int

	GMusic *gmusic.GMusic

	// Better:
	// Database *Database
	DB         *bolt.DB
	ArtistsMap map[string]bool
	Artists    sort.StringSlice
	Songs      map[string][]string
	Albums     map[string][]string

	LastAlbum string
	Status    *Status
}

// New creates a new UI
func New(gmusic *gmusic.GMusic, db *bolt.DB) (*App, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	err = screen.Init()
	if err != nil {
		return nil, err
	}
	width, height := screen.Size()
	return &App{
		Screen:     screen,
		Width:      width,
		Height:     height,
		GMusic:     gmusic,
		DB:         db,
		ArtistsMap: make(map[string]bool),
		Artists:    sort.StringSlice{},
		Songs:      make(map[string][]string),
		Albums:     make(map[string][]string),
		Status: &Status{
			ScrOffset: map[bool]int{
				false: 0, // in artists view
				true:  0, // in tracks view
			},
			Offset: 0,
			CurPos: map[bool]int{
				false: -1, // same as in scrOffset. -1 is because the artist is unfolded (yet)
				true:  0,
			},
			NumAlbum: map[bool]int{
				false: -1, // same as in scrOffset. -1 is because the artist is unfolded (yet)
				true:  0,
			},
			InTracks:  false,
			InSearch:  false,
			NumTrack:  0,
			CurArtist: make(chan string),
			Queue:     make([][]*music.BTrack, 0),
		},
	}, nil
}

func (app *App) Run() {
	app.populateArtists()
	go app.player()
	app.mainLoop()
}

func (app *App) populateArtists() {
	app.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Library"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if !app.ArtistsMap[string(k)] {
				app.ArtistsMap[string(k)] = false
			}
			if v == nil {
				if err := b.Bucket(k).ForEach(func(kk []byte, vv []byte) error {
					app.Albums[string(k)] = append(app.Albums[string(k)], string(kk))

					return nil
				}); err != nil {
					log.Fatalf("Can't populate artists: %s", err)
				}
			}

		}
		for k := range app.ArtistsMap {
			app.Artists = append(app.Artists, k)
		}
		app.Artists.Sort()

		return nil
	})
}

func (app *App) populateSongs() {
	if err := app.DB.View(func(tx *bolt.Tx) error {
		l := tx.Bucket([]byte("Library"))
		i := app.Status.CurPos[false] - 1 + app.Status.ScrOffset[false]
		b := l.Bucket([]byte(app.Artists[i-app.numAlb(i)]))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil {
				cc := b.Bucket(k).Cursor()
				for kk, vv := cc.First(); kk != nil; kk, vv = cc.Next() {
					app.Songs[string(k)] = append(app.Songs[string(k)], string(vv))
				}
			}

		}

		return nil

	}); err != nil {
		log.Fatalf("Can't populate songs: %s", err)
	}

}

func (app *App) search() {
	app.Status.InTracks = false
	app.Status.InSearch = true
	app.Status.NumTrack = 0
	app.Status.CurPos[true] = 2
	app.Status.ScrOffset[true] = 0
	app.Status.Query = []rune{}
	for {
		app.printStatus()
		app.Screen.Show()
		ev := app.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				app.Status.Query = append(app.Status.Query, ev.Rune())
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(app.Status.Query) > 0 {
					app.Screen.SetContent(len(app.Status.Query), app.Height-1, ' ', nil, dfStyle)
					app.Status.Query = app.Status.Query[:len(app.Status.Query)-1]
				} else {
					app.Status.InSearch = false
					return
				}
			case tcell.KeyEnter:
				app.Status.InSearch = false
				return
			}
		}
		app.searchQuery()
	}
}

func (app *App) searchQuery() {
	var i int
	if !app.Status.InSearch {
		i = app.Status.ScrOffset[false] + app.Status.CurPos[false]
	}
	if len(app.Status.Query) > 0 {
		for i < len(app.Artists) {
			if strings.HasPrefix(strings.ToLower(app.Artists[i]), strings.ToLower(string(app.Status.Query))) {
				if i > 2 {
					app.Status.ScrOffset[false] = i - 2
					app.Status.CurPos[false] = 3
				} else {
					app.Status.ScrOffset[false] = 0
					app.Status.CurPos[false] = i + 1
				}
				app.updateUI()
				return
			}
			i++
		}
	}

}

func (app *App) randomizeArtists() {
	var temp = make(sort.StringSlice, len(app.Artists))
	perm := rand.Perm(len(app.Artists))
	for i, v := range perm {
		temp[v] = app.Artists[i]
	}

	app.Artists = temp
	app.updateUI()

}

func (app *App) mainLoop() {
	for {
		app.Screen.Show()
		ev := app.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			width, height := app.Screen.Size()
			app.Width = width
			app.Height = height
			// updateUI(s)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyPgDn:
				app.pageDn()
			case tcell.KeyPgUp:
				app.pageUp()
			case tcell.KeyEnd:
				app.scrollDown()
			case tcell.KeyHome:
				app.scrollUp()
			case tcell.KeyTab:
				app.toggleView()
			case tcell.KeyUp:
				app.upEntry()
			case tcell.KeyDown:
				app.downEntry()
			case tcell.KeyEnter:
				// app.Status.State <- play
				// i := app.Status.CurPos[false] - 1 + app.Status.ScrOffset[false]
				// app.Status.CurArtist <- app.Artists[i-app.Status.NumAlbum(i)]
			}
			switch ev.Rune() {
			case '/':
				app.search()
			case 'q':
				return
			case 'j':
				app.downEntry()
			case 'k':
				app.upEntry()
			case ' ':
				app.toggleAlbums()
			case 'u':
				music.RefreshLibrary(app.DB, app.GMusic)
				app.populateArtists()
			case 'x':
				// app.Status.State <- play
				// i := app.Status.CurPos[false] - 1 + app.Status.ScrOffset[false]
				// app.Status.CurArtist <- app.Artists[i-app.Status.NumAlbum(i)]
			case 'v':
				app.Status.State <- stop
			case 'c':
				app.Status.State <- pause
			case 'b':
				app.Status.State <- next
			case 'z':
				app.Status.State <- prev
			case 'n':
				app.searchQuery()
			case 'R':
				app.randomizeArtists()
			}
		}
		app.updateUI()
	}
}
