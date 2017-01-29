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
	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	dfStyle = tcell.StyleDefault
)

var (
	styleNormal = tcell.StyleDefault.
			Foreground(tcell.ColorSilver).
			Background(tcell.ColorBlack)
	styleGood = tcell.StyleDefault.
			Foreground(tcell.ColorGreen).
			Background(tcell.ColorBlack)
	styleWarn = tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorBlack)
	styleError = tcell.StyleDefault.
			Foreground(tcell.ColorMaroon).
			Background(tcell.ColorBlack)

	// highlight tcell style
	hlStyle = tcell.StyleDefault.Bold(true).Reverse(true)
	// tcell style used for bars
	barStyle = tcell.StyleDefault.Reverse(true)
	// tcell style usef for albums in tracks view
	alStyle = tcell.StyleDefault.Bold(true)

	// hlStyle = tcell.StyleDefault.
	// 	Bold(true).
	// 	Foreground(tcell.ColorYellow)

	// alStyle = tcell.StyleDefault.
	// 	Bold(true).
	// 	Foreground(tcell.ColorGreen)
)

func fill(screen tcell.Screen, x, y, w, h int, v rune, stl tcell.Style) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			screen.SetContent(lx+x, ly+y, v, nil, stl)
		}
	}
}

func print(screen tcell.Screen, x, y int, stl tcell.Style, l string) {
	for _, v := range l {
		screen.SetContent(x, y, v, nil, stl)
		x += runewidth.RuneWidth(v)
	}
}

func printSingleItem(s tcell.Screen, x, y int, sty tcell.Style, l string, offset int, artist bool, width int) {
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
