package main

import (
	"log"

	"github.com/boltdb/bolt"
)

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
