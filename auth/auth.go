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

package auth

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/budkin/gmusic"
	"github.com/howeyc/gopass"

	"github.com/budkin/jam/music"
	"github.com/budkin/jam/storage"
)

func loginFromDatabase(db *bolt.DB) (*gmusic.GMusic, error) {
	auth, deviceID, err := storage.ReadCredentials(db)
	if err != nil {
		return nil, err
	}
	return &gmusic.GMusic{
		Auth:     string(auth),
		DeviceID: string(deviceID),
	}, nil
}

func CheckCreds(db *bolt.DB) (*gmusic.GMusic, error) {
	gm, err := loginFromDatabase(db)
	if err != nil {
		gm, err = authenticate()
		if err != nil {
			return nil, err
		}
		err = music.RefreshLibrary(db, gm)
	}
	if err != nil {
		return nil, err
	}
	err = storage.WriteCredentials(db, gm.Auth, gm.DeviceID)
	if err != nil {
		return nil, err
	}
	return gm, nil
}

func authenticate() (*gmusic.GMusic, error) {
	email := askForEmail()
	password, err := askForPassword()
	if err != nil {
		return nil, err
	}
	return gmusic.Login(email, string(password))
}

func askForEmail() string {
	var email string
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	return email
}

func askForPassword() ([]byte, error) {
	var password []byte
	fmt.Print("Password: ")
	password, err := gopass.GetPasswdMasked()
	if err != nil {
		return nil, err
	}
	return password, nil
}
