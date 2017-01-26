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
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"

	"github.com/budkin/jam/version"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = "Jam - %s\n"
)

var (
	vers  bool
	debug bool
)

func init() {
	// parse flags
	flag.BoolVar(&vers, "version", false, "print version and exit")
	flag.BoolVar(&debug, "debug", false, "debug")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.Version))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vers {
		fmt.Printf("%s\n", version.Version)
		os.Exit(0)
	}
}

func main() {

	var err error
	db, err = bolt.Open(fullDbPath(), 0600, nil)
	checkErr(err)
	defer db.Close()
	checkCreds()
	//refreshLibrary()
	initUI()
	//play(g)

}

func checkErr(e error) {
	if e != nil {
		if debug {
			panic(e)
		}
		log.Fatal(e)
	}
}
