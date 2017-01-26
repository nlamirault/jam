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

package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/budkin/gmusic"
	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
)

var (
	artistsMap = make(map[string]bool) // map of artists, so they are shown in random order
	// the bool is whether the artist is unfolded in the artists view (albums are shown)
	artists       = make([]string, 0)         // slice of artists, populated from artistsMap
	songs         = make(map[string][]string) //map of currently selected artist's albums and songs in it
	albums        = make(map[string][]string) // map of artists and their albums
	width, height int                         // width and height of the window
	lastAlbum     string
	hlStyle       = tcell.StyleDefault.Bold(true).Reverse(true) // highlight tcell style
	barStyle      = tcell.StyleDefault.Reverse(true)            // tcell style used for bars
	alStyle       = tcell.StyleDefault.Bold(true)               // tcell style usef for albums in tracks view
	scrOffset     = map[bool]int{                               // how far we are from the first shown artist or track
		false: 0, // in artists view
		true:  0, // in tracks view
	}
	offset = 0             // horizontal offset when printing lines, e.g. 2 prints "  Artist", instead of "Artist"
	curPos = map[bool]int{ //cursor position on Y axis
		false: 1, // same as in
		true:  2, // scrOffset
	}
	numAlbum = map[bool]int{ // number of currently selected album
		false: -1, // same as in scrOffset. -1 is because the artist is unfolded (yet)
		true:  0,
	}
	inTracks = false // whether cursor is actvie in the tracks or artists view
	inSearch = false // whether user is searching for an artist

	numTrack  = 0                  // number of track in tracks view
	curArtist = make(chan string)  // song artists are not stored in the db along with song names, so this is a workaround
	queue     [][]*bTrack          // playlist, updated on each movement of cursor in artists view
	query     = make([]rune, 0)    // search query
	db        = new(bolt.DB)       // bolt database
	gm        = new(gmusic.GMusic) // gmusic instance
	state     = make(chan int)     // player's state: play, pause, stop, etc

	defDur    = time.Duration(0) // this and below are used for when nothing is playing
	defTrack  = &bTrack{}
	defArtist = ""
)

const (
	dfStyle = tcell.StyleDefault
)

const (
	play = iota
	stop
	pause
	next
	prev
)

func initUI() {
	s, e := tcell.NewScreen()
	checkErr(e)
	e = s.Init()
	checkErr(e)
	width, height = s.Size()
	defer s.Fini()
	populateArtists()

	go player(s)
	mainLoop(s)
}

func mainLoop(s tcell.Screen) {
	for {
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			width, height = s.Size()
			updateUI(s)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyPgDn:
				pageDn()
			case tcell.KeyPgUp:
				pageUp()
			case tcell.KeyEnd:
				scrollDown()
			case tcell.KeyHome:
				scrollUp()
			case tcell.KeyTab:
				toggleView(s)
			case tcell.KeyUp:
				upEntry(s)
			case tcell.KeyDown:
				downEntry(s)
			case tcell.KeyEnter:
				state <- play
				i := curPos[false] - 1 + scrOffset[false]
				curArtist <- artists[i-numAlb(i)]
			}
			switch ev.Rune() {
			case '/':
				search(s)
			case 'q':
				return
			case 'j':
				downEntry(s)
			case 'k':
				upEntry(s)
			case ' ':
				toggleAlbums(s)
			case 'u':
				refreshLibrary()
				populateArtists()
			case 'x':
				state <- play
				i := curPos[false] - 1 + scrOffset[false]
				curArtist <- artists[i-numAlb(i)]
			case 'v':
				state <- stop
			case 'c':
				state <- pause
			case 'b':
				state <- next
			case 'z':
				state <- prev
			}
		}
		updateUI(s)
	}
}

func search(s tcell.Screen) {
	inTracks = false
	inSearch = true
	numTrack = 0
	curPos[true] = 2
	scrOffset[true] = 0
	for {
		printStatus(s)
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				query = append(query, ev.Rune())
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(query) > 0 {
					s.SetContent(len(query), height-1, ' ', nil, dfStyle)
					query = query[:len(query)-1]
				} else {
					inSearch = false
					return
				}
			case tcell.KeyEnter:
				query = []rune{}
				inSearch = false
				return
			}
		}
		searchQuery(s)
	}
}

func searchQuery(s tcell.Screen) {
	if len(query) > 0 {
		for i := scrOffset[false] + curPos[false]; i < len(artists); i++ {
			if strings.HasPrefix(strings.ToLower(artists[i]), strings.ToLower(string(query))) {
				if i > 2 {
					scrOffset[false] = i - 2
					curPos[false] = 3
				} else {
					scrOffset[false] = 0
					curPos[false] = i + 1
				}
				updateUI(s)
				return
			}
		}
	}

}

func pageUp() {
}

func pageDn() {
}

func printStatus(s tcell.Screen) {
	if inSearch {
		s.SetContent(0, height-1, '/', nil, dfStyle)
	}
	for i, r := range query {
		s.SetContent(i+1, height-1, r, nil, dfStyle)
	}
}

func printBar(s tcell.Screen, dur time.Duration, track *bTrack, artist string) {
	strdur := ""
	str := fmt.Sprintf(" %02v:%02v %v - %v ", int(dur.Minutes()), int(dur.Seconds())%60, artist, track.Title)
	lenstr := 0
	for _, r := range str {
		lenstr += runewidth.RuneWidth(r)
	}
	print(s, 0, height-2, barStyle, str)
	leng := width - lenstr
	durat, _ := strconv.Atoi(track.DurationMillis)
	dura := time.Duration(durat) * time.Millisecond
	for i := 0.0; i < float64(leng)/dura.Seconds()*dur.Seconds(); i += 1.0 {
		strdur += "—"
	}
	print(s, lenstr, height-2, barStyle, strdur)
	s.Show()

}

func scrollDown() {
	var length int
	if !inTracks {
		length = len(artists)
		numAlbum[false] = -1
	} else if numAlbum[false] == -1 {
		length = numSongs()
		numAlbum[true] = len(songs) - 1
		numTrack = len(songs[lastAlbum]) - 1
	} else {
		length = len(songs[albums[artists[curPos[false]-1+scrOffset[false]-1*(numAlbum[false]+1)]][numAlbum[true]]])
		numTrack = length - 1
	}
loop:
	for curPos[inTracks]+scrOffset[inTracks] < length {
		for curPos[inTracks] < height-3 {
			if curPos[inTracks] == length {
				break loop
			}
			curPos[inTracks]++
		}
		scrOffset[inTracks]++
	}

}

func scrollUp() {
	if !inTracks {
		curPos[inTracks] = 1
		scrOffset[inTracks] = 0
		numAlbum[false] = -1
	} else if numAlbum[false] == -1 {
		curPos[inTracks] = 2
		scrOffset[inTracks] = 0
		numTrack = 0
		numAlbum[true] = 0
	} else {
		curPos[inTracks] = 1
		scrOffset[inTracks] = 0
		numTrack = 0
	}
}

func toggleView(s tcell.Screen) {
	inTracks = !inTracks

}

func toggleAlbums(s tcell.Screen) {
	numAlbum[false] = -1
	i := curPos[false] - 1 + scrOffset[false]
	artistsMap[artists[i-numAlb(i)]] = !artistsMap[artists[i-numAlb(i)]]

	if artistsMap[artists[i-numAlb(i)]] {
		appendAlbums()
	} else {
		removeAlbums()
		numAlbum[true] = 0
		numTrack = 0
		scrOffset[true] = 0
		curPos[true] = 2
	}
}

func appendAlbums() {
	var b []string
	boo := artists[curPos[false]-1+scrOffset[false]] == artists[len(artists)-1]
	a := make([]string, len(albums[artists[curPos[false]-1+scrOffset[false]]]))
	if !boo {
		b = make([]string, len(artists[curPos[false]+scrOffset[false]:]))
		copy(b, artists[curPos[false]-1+scrOffset[false]+1:])
	}
	artists = append(artists[:curPos[false]-1+scrOffset[false]+1], a...)
	if !boo {
		artists = append(artists, b...)
	}
}

func removeAlbums() {
	i := curPos[false] - 1 + scrOffset[false]
	if artists[i-numAlb(i)] != artists[len(artists)-1-len(albums[artists[i-numAlb(i)]])] {
		curPos[false] -= numAlb(i)
		artists = append(artists[:i-numAlb(i)+1], artists[i-numAlb(i)+len(albums[artists[i-numAlb(i)]])+1:]...)
	} else {

		a := numAlb(i)
		artists = artists[:i-numAlb(i)+1]
		scrOffset[false] -= a
	}

}

func updateUI(s tcell.Screen) {
	s.Clear()
	fill(s, width/3, 1, 1, height-3, '│', dfStyle)

	printHeader(s)
	if len(artists) > 0 {
		printArtists(s, scrOffset[false], height+scrOffset[false]-3)
		printSongs(s, scrOffset[true], height+scrOffset[true]-3)
		for curPos[false] > height-3 {
			curPos[false]--
		}
		for curPos[true] > height-3 {
			curPos[true] = 2
			numTrack = 0
			numAlbum[true] = 0
		}
		hlEntry(s)
	}
	printStatus(s)
	printBar(s, defDur, defTrack, defArtist)
}

func printHeader(s tcell.Screen) {
	fill(s, 0, 0, width, 1, ' ', hlStyle)
	print(s, 1, 0, hlStyle, "Artist / Album")
	print(s, width/3+2, 0, hlStyle, "Track")
	print(s, width-8, 0, hlStyle, "Library")

	fill(s, 0, height-2, width, 1, ' ', hlStyle)
}

func hlEntry(s tcell.Screen) {
	i := curPos[false] - 1 + scrOffset[false]
	if artists[i] != "" {
		printSingleItem(s, 0, curPos[false], hlStyle, artists[i], 1, true)
	} else {
		printSingleItem(s, 0, curPos[false], hlStyle, albums[artists[i-numAlb(i)]][numAlb(i)-1], 3, true)

	}

	if inTracks {
		song := songs[albums[artists[i-numAlb(i)]][numAlbum[true]]][numTrack]
		js := new(bTrack)
		json.Unmarshal([]byte(song), js)
		printSingleItem(s, width/3+2, curPos[inTracks], hlStyle, makeSongLine(js), 0, false)
	}
}

func downEntry(s tcell.Screen) {
	var high int
	switch inTracks {
	case true:
		high = 0
		if numAlbum[false] == -1 {
			for _, sngs := range songs {
				high++
				for _ = range sngs {
					high++
				}
			}
		} else {
			high = len(songs[albums[artists[curPos[false]-1+scrOffset[false]-1*(numAlbum[false]+1)]][numAlbum[true]]])
		}

		if numTrack < len(songs[albums[artists[curPos[false]-1+scrOffset[false]-1*(numAlbum[false]+1)]][numAlbum[true]]])-1 {
			numTrack++
		} else if numAlbum[true] < len(albums[artists[curPos[false]-1+scrOffset[false]-1*(numAlbum[false]+1)]])-1 && numAlbum[false] == -1 {
			numAlbum[true]++
			if curPos[true] > height-6 {
				scrOffset[true]++
			} else {
				curPos[true]++
			}
			numTrack = 0
		}
	case false:
		high = len(artists)
		if artistsMap[artists[curPos[false]-1+scrOffset[false]]] ||
			artists[curPos[false]-1+scrOffset[false]] == "" {
			if numAlbum[false] < len(albums[artists[curPos[false]-1+scrOffset[false]-1*(numAlbum[false]+1)]])-1 {
				numAlbum[false]++
				numAlbum[true] = numAlbum[false]

				curPos[true] = 1
			}

		} else {
			numAlbum[false] = -1
			numAlbum[true] = 0
		}

		numTrack = 0
		curPos[true] = 2
		scrOffset[true] = 0

	}

	if curPos[inTracks]+scrOffset[inTracks] < high {
		if curPos[inTracks] > height-6 && curPos[inTracks]+scrOffset[inTracks] < high-2 {
			scrOffset[inTracks]++
		} else if curPos[inTracks] == height-3 && curPos[inTracks]+scrOffset[inTracks] < high {
			scrOffset[inTracks]++
		} else {
			curPos[inTracks]++
		}
	}
	if !inTracks {
		if artists[curPos[false]-1+scrOffset[false]] != "" {
			numAlbum[false] = -1
		} else {
			curPos[true] = 1
		}
	}
}

func upEntry(s tcell.Screen) {

	switch inTracks {
	case false:
		numTrack = 0
		curPos[true] = 2
		scrOffset[true] = 0
	case true:
		if numAlbum[false] == -1 {
			if numTrack > 0 {
				numTrack--
			} else if numTrack == 0 && numAlbum[true] > 0 {
				numAlbum[true]--
				if curPos[true] > 3 {
					curPos[true]--
				} else {
					scrOffset[true]--
				}
				songs := songs[albums[artists[curPos[false]-1+scrOffset[false]]][numAlbum[true]]]
				numTrack = len(songs) - 1
			} else if curPos[true] < 3 {
				return
			}
		} else {
			if numTrack > 0 {
				numTrack--
			}
		}
	}

	if curPos[inTracks] > 3 {
		curPos[inTracks]--
	} else if scrOffset[inTracks]+curPos[inTracks] > 3 {
		scrOffset[inTracks]--
	} else if scrOffset[inTracks]+curPos[inTracks] > 1 {
		curPos[inTracks]--
	}

	if !inTracks {
		if artistsMap[artists[curPos[false]-1+scrOffset[false]]] {
			numAlbum[false] = -1
			numAlbum[true] = 0
		} else if artists[curPos[false]-1+scrOffset[false]] == "" {
			curPosTemp := curPos[false] - 1
			for artists[curPosTemp+scrOffset[false]] == "" {
				curPosTemp--
			}
			numAlbum[false] = curPos[false] - curPosTemp - 1 - 1
			if numAlbum[false] > -1 {

				numAlbum[true] = numAlbum[false]
			}
			curPos[true] = 1

		}
	}

}

func print(s tcell.Screen, x, y int, stl tcell.Style, l string) {
	for _, v := range l {
		s.SetContent(x, y, v, nil, stl)
		x += runewidth.RuneWidth(v)
	}
}

func fill(s tcell.Screen, x, y, w, h int, v rune, stl tcell.Style) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			s.SetContent(lx+x, ly+y, v, nil, stl)
		}
	}
}

func printArtists(s tcell.Screen, beg, end int) {
	for len(artists) < end {
		end--
	}
	j := 1
	for beg < end {
		if artists[beg] != "" {
			printSingleItem(s, 0, j, dfStyle, artists[beg], 1, true)

			j++
			beg++
		} else {
			printSingleItem(s, 0, j, dfStyle, albums[artists[beg-numAlb(beg)]][numAlb(beg)-1], 3, true)
			j++
			beg++
		}

	}
}

func numAlb(k int) int {
	var i int
	for artists[k] == "" {
		i++
		k--
	}
	return i

}
func numSongs() int {
	j := 0
	for _, sngs := range songs {
		j++
		for _ = range sngs {
			j++
		}
	}
	return j
}

func printSongs(s tcell.Screen, beg, end int) {
	queue = [][]*bTrack{}
	populateSongs(s)
	i, k := 0, 1
	if numAlbum[false] == -1 {

		j := numSongs()
		if j < end {
			end = j
		}
		keys := []string{}
		for a := range songs {
			keys = append(keys, a)
		}
		sort.Strings(keys)
		lastAlbum = keys[len(keys)-1]
		for _, key := range keys {
			que := []*bTrack{}
			if i >= beg && i < end {
				printAlbum(s, k, key)
				k++
			}
			i++
			for _, song := range songs[key] {
				js := new(bTrack)
				json.Unmarshal([]byte(song), js)
				if i >= beg && i < end {
					printSingleItem(s, width/3+2, k, dfStyle, makeSongLine(js), 0, false)
					que = append(que, js)
					k++
				}
				i++
			}
			queue = append(queue, que)

		}
	} else {
		que := []*bTrack{}
		j := curPos[false] - 1 + scrOffset[false]
		if len(songs[albums[artists[j-numAlb(j)]][numAlbum[false]]]) < end {
			end = len(songs[albums[artists[j-numAlb(j)]][numAlbum[false]]])
		}

		for _, song := range songs[albums[artists[j-numAlb(j)]][numAlbum[false]]] {
			js := new(bTrack)
			json.Unmarshal([]byte(song), js)
			if i >= beg {
				printSingleItem(s, width/3+2, k, dfStyle, makeSongLine(js), 0, false)
				que = append(que, js)
				k++
			}
			i++
		}
		for l := numAlb(j); l > 1; l-- {
			queue = append(queue, []*bTrack{})
		}
		queue = append(queue, que)
	}
}

func printAlbum(s tcell.Screen, y int, alb string) {
	makeAlbumLine(&alb)
	print(s, width/3+2, y, alStyle, alb)

}

func makeSongLine(track *bTrack) string {
	var res string
	length := 0

	res = fmt.Sprintf("%2v. %v", track.TrackNumber, track.Title)
	for _, c := range res {
		length += runewidth.RuneWidth(c)
	}

	for length < width-width/3+2-16 {
		res += " "
		length++
	}

	for length > width-width/3+2-16 {
		run := []rune(res)
		length -= runewidth.RuneWidth(run[len(run)-1])
		run = run[:len(run)-1]
		res = string(run)
	}

	du, _ := strconv.Atoi(track.DurationMillis)
	dur := time.Duration(time.Millisecond * time.Duration(du))
	res += fmt.Sprintf(" %4v %02v:%02v ", track.Year, int(dur.Minutes()), int(dur.Seconds())%60)
	return res
}

func makeAlbumLine(str *string) {
	len := 0
	for _, c := range *str {
		len += runewidth.RuneWidth(c)
	}
	*str += " "
	len++
	for len < (width/3)*2 {
		*str += "─"
		len++
	}
}

func populateSongs(s tcell.Screen) {
	songs = map[string][]string{}
	if err := db.View(func(tx *bolt.Tx) error {
		l := tx.Bucket([]byte("Library"))
		i := curPos[false] - 1 + scrOffset[false]
		b := l.Bucket([]byte(artists[i-numAlb(i)]))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil {
				cc := b.Bucket(k).Cursor()
				for kk, vv := cc.First(); kk != nil; kk, vv = cc.Next() {
					songs[string(k)] = append(songs[string(k)], string(vv))
				}
			}

		}

		return nil

	}); err != nil {
		checkErr(err)
	}

}

func printSingleItem(s tcell.Screen, x, y int, sty tcell.Style, l string, offset int, artist bool) {
	if artist {
		makeLine(&l, width/3, offset)
	}
	print(s, x, y, sty, l)
}

// this makes length of a string with unicode characters equal to width/3
func makeLine(l *string, strLen int, offset int) {

	//*l = " " + *l
	for i := 0; i < offset; i++ {
		*l = " " + *l
	}
	len := 0
	for _, c := range *l {
		len += runewidth.RuneWidth(c)
	}
	if len > strLen {
		*l = (*l)[:strLen]
		return
	}
	for len < strLen {
		*l += " "
		len++
	}
}

func populateArtists() {

	albums = make(map[string][]string)
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Library"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if !artistsMap[string(k)] {
				artistsMap[string(k)] = false
			}
			if v == nil {
				if err := b.Bucket(k).ForEach(func(kk []byte, vv []byte) error {
					albums[string(k)] = append(albums[string(k)], string(kk))

					return nil
				}); err != nil {
					checkErr(err)
				}
			}

		}
		artists = []string{}
		for k := range artistsMap {
			artists = append(artists, k)
		}

		return nil
	})
}
