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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/boltdb/bolt"
	"github.com/budkin/gmusic"
	"github.com/howeyc/gopass"
)

func fullDbPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "jamdb")
	}
	return filepath.Join(os.Getenv("HOME"), ".local/share/jamdb")
}

func checkCreds() {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("#AuthDetails"))
		if b == nil {
			return errors.New("no bucket")
		}
		gm = &gmusic.GMusic{
			Auth:     string(b.Get([]byte("Auth"))),
			DeviceID: string(b.Get([]byte("DeviceID"))),
		}
		return nil
	})

	if err != nil {
		err = authenticate()
		checkErr(err)
		refreshLibrary()
		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("#AuthDetails"))
			checkErr(err)

			err = b.Put([]byte("Auth"), []byte(gm.Auth))
			err = b.Put([]byte("DeviceID"), []byte(gm.DeviceID))

			return err
		})
	}
}

func authenticate() error {
	email := askForEmail()
	password := askForPassword()
	var err error
	gm, err = gmusic.Login(email, string(password))
	return err
}

func askForEmail() string {
	var email string
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	return email
}

func askForPassword() []byte {
	var password []byte
	fmt.Print("Password: ")
	password, err := gopass.GetPasswdMasked()
	checkErr(err)
	return password
}
