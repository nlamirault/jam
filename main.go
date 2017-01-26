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
	vers bool
)

func init() {
	// parse flags
	flag.BoolVar(&vers, "version", false, "print version and exit")

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
		log.Fatal(e)
	}
}
