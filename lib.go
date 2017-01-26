package main

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
)

type bTrack struct {
	//AlbumArtist    string
	DiscNumber     uint8
	TrackNumber    uint8
	DurationMillis string
	EstimatedSize  string
	ID             string
	PlayCount      uint32
	Title          string
	Year           int
}

func refreshLibrary() {
	//db, err := bolt.Open(fullDbPath(), 0600, nil)
	//checkErr(err)
	//defer db.Close()

	tracks, err := gm.ListTracks()
	checkErr(err)
	err = db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("Library"))

		lib, err := tx.CreateBucketIfNotExists([]byte("Library"))
		checkErr(err)
		for _, t := range tracks {
			artist, err := lib.CreateBucketIfNotExists([]byte(t.Artist))
			checkErr(err)
			if t.Album == "" {
				t.Album = "Unknown Album"
			}
			album, err := artist.CreateBucketIfNotExists([]byte(t.Album))
			checkErr(err)

			//id, _ := album.NextSequence()

			bt := bTrack{t.DiscNumber, t.TrackNumber, t.DurationMillis,
				t.EstimatedSize, t.ID, t.PlayCount, t.Title, t.Year}
			//trackNumber, _ := album.NextSequence()
			buf, err := json.Marshal(bt)
			checkErr(err)
			var key string
			if t.TrackNumber < 10 {
				key = strconv.Itoa(int(t.DiscNumber)) + "0" + strconv.Itoa(int(t.TrackNumber))
			} else {
				key = strconv.Itoa(int(t.DiscNumber)) + strconv.Itoa(int(t.TrackNumber))
			}

			err = album.Put([]byte(key), buf)
			checkErr(err)

		}

		return nil
	})

}

func itob(v uint64) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(v))
	return b
}
