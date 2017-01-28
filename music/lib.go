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

package music

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/budkin/gmusic"
)

type BTrack struct {
	//AlbumArtist    string
	DiscNumber     uint8
	TrackNumber    uint32
	DurationMillis string
	EstimatedSize  string
	ID             string
	PlayCount      uint32
	Title          string
	Year           int
}

func RefreshLibrary(db *bolt.DB, gm *gmusic.GMusic) error {
	tracks, err := gm.ListTracks()
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_ = tx.DeleteBucket([]byte("Library"))
		lib, err := tx.CreateBucketIfNotExists([]byte("Library"))
		if err != nil {
			return err
		}
		for _, t := range tracks {
			artist, err := lib.CreateBucketIfNotExists([]byte(t.Artist))
			if err != nil {
				return err
			}
			if t.Album == "" {
				t.Album = "Unknown Album"
			}
			album, err := artist.CreateBucketIfNotExists([]byte(t.Album))
			if err != nil {
				return err
			}

			//id, _ := album.NextSequence()

			bt := BTrack{t.DiscNumber, t.TrackNumber, t.DurationMillis, t.EstimatedSize, t.ID, t.PlayCount, t.Title, t.Year}
			//trackNumber, _ := album.NextSequence()
			buf, err := json.Marshal(bt)
			if err != nil {
				return err
			}
			var key string
			if t.TrackNumber < 10 {
				key = strconv.Itoa(int(t.DiscNumber)) + "0" + strconv.Itoa(int(t.TrackNumber))
			} else {
				key = strconv.Itoa(int(t.DiscNumber)) + strconv.Itoa(int(t.TrackNumber))
			}

			err = album.Put([]byte(key), buf)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return nil
}

func itob(v uint64) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(v))
	return b
}
