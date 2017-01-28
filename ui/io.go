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
	"fmt"
	"strconv"
	"time"

	// "github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"

	"github.com/budkin/jam/music"
)

func makeSongLine(track *music.BTrack, width int) string {
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

func makeAlbumLine(str *string, width int) {
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
