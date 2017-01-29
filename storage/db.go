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

package storage

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/boltdb/bolt"
)

func fullDbPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "jamdb")
	}
	return filepath.Join(os.Getenv("HOME"), ".local/share/jamdb")
}

func Open() (*bolt.DB, error) {
	return bolt.Open(fullDbPath(), 0600, nil)
}

func ReadCredentials(db *bolt.DB) ([]byte, []byte, error) {
	var auth []byte
	var deviceID []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("#AuthDetails"))
		if b == nil {
			return errors.New("No bucket into database")
		}
		auth = b.Get([]byte("Auth"))
		deviceID = b.Get([]byte("DeviceID"))
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return auth, deviceID, nil
}

func WriteCredentials(db *bolt.DB, auth string, deviceID string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("#AuthDetails"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("Auth"), []byte(auth))
		err = b.Put([]byte("DeviceID"), []byte(deviceID))
		return err
	})
}
