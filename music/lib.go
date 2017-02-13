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
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/budkin/gmusic"
)

type BTrack struct {
	Artist         string
	DiscNumber     uint8
	TrackNumber    uint32
	DurationMillis string
	EstimatedSize  string
	ID             string
	PlayCount      uint32
	Title          string
	Year           uint32
}

func RefreshLibrary(db *bolt.DB, gm *gmusic.GMusic) error {
	//db, err := bolt.Open(fullDbPath(), 0600, nil)
	//checkErr(err)
	//defer db.Close()
	var artist *bolt.Bucket
	var err error
	var key string
	var temp int
	var mixedAlbum bool
	var buf []byte
	tracks, err := gm.ListTracks()
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("Library"))

		lib, err := tx.CreateBucketIfNotExists([]byte("Library"))
		if err != nil {
			return err
		}
		for i := 0; i < len(tracks); i++ {
			temp = i + 1
			for mixedAlbum = false; temp < len(tracks) && tracks[temp].Album == tracks[i].Album; temp++ {
				if tracks[temp].Artist != tracks[i].Artist {
					mixedAlbum = true
				}
			}

			if mixedAlbum {
				artist, err = lib.CreateBucketIfNotExists([]byte("Various Artists"))
				if err != nil {
					return err
				}

				if tracks[i].Album == "" {
					tracks[i].Album = "Unknown Album"
				}
				album, err := artist.CreateBucketIfNotExists([]byte(tracks[i].Album))
				if err != nil {
					return err
				}
				for i < temp {
					bt := BTrack{tracks[i].Artist, tracks[i].DiscNumber, tracks[i].TrackNumber, tracks[i].DurationMillis,
						tracks[i].EstimatedSize, tracks[i].ID, tracks[i].PlayCount, tracks[i].Title, tracks[i].Year}
					buf, err = json.Marshal(bt)
					if err != nil {
						return err
					}

					if tracks[i].TrackNumber < 10 {
						key = strconv.Itoa(int(tracks[i].DiscNumber)) + "0" + strconv.Itoa(int(tracks[i].TrackNumber))
					} else {
						key = strconv.Itoa(int(tracks[i].DiscNumber)) + strconv.Itoa(int(tracks[i].TrackNumber))
					}

					err = album.Put([]byte(key), buf)
					if err != nil {
						return err
					}
					i++
				}
				continue
			} else {
				if tracks[i].Artist == "" {
					if tracks[i].AlbumArtist == "" {
						tracks[i].Artist = "Unknown Artist"
					} else {
						tracks[i].Artist = tracks[i].AlbumArtist
					}
				}
				artist, err = lib.CreateBucketIfNotExists([]byte(tracks[i].Artist))
				if err != nil {
					return err
				}
			}
			if tracks[i].Album == "" {
				tracks[i].Album = "Unknown Album"
			}
			album, err := artist.CreateBucketIfNotExists([]byte(tracks[i].Album))
			if err != nil {
				return err
			}

			bt := BTrack{tracks[i].Artist, tracks[i].DiscNumber, tracks[i].TrackNumber, tracks[i].DurationMillis,
				tracks[i].EstimatedSize, tracks[i].ID, tracks[i].PlayCount, tracks[i].Title, tracks[i].Year}
			buf, err = json.Marshal(bt)
			if err != nil {
				return err
			}

			if tracks[i].TrackNumber < 10 {
				key = strconv.Itoa(int(tracks[i].DiscNumber)) + "0" + strconv.Itoa(int(tracks[i].TrackNumber))
			} else {
				key = strconv.Itoa(int(tracks[i].DiscNumber)) + strconv.Itoa(int(tracks[i].TrackNumber))
			}

			err = album.Put([]byte(key), buf)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}
