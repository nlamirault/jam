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

	"github.com/budkin/jam/lastfm"
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

func CheckCreds(db *bolt.DB, lastFM bool) (*gmusic.GMusic, *lastfm.Client, string, error) {
	gm, err := loginFromDatabase(db)
	if err != nil {
		gm, err = authenticate()
		if err != nil {
			return nil, nil, "", err
		}
		err = storage.WriteCredentials(db, gm.Auth, gm.DeviceID)
		if err != nil {
			return nil, nil, "", err
		}
		fmt.Println("Syncing database is in order, may take a few seconds")
		err = music.RefreshLibrary(db, gm)
	}
	if err != nil {
		return nil, nil, "", err
	}
	if !lastFM {
		return gm, nil, "", nil
	}

	lmclient := lastfm.New(
		string([]byte{0x62, 0x39, 0x30, 0x36, 0x65, 0x62, 0x63, 0x35, 0x39, 0x35, 0x34, 0x63, 0x37, 0x65, 0x63, 0x39, 0x66, 0x39, 0x65, 0x63, 0x64, 0x32, 0x66, 0x66, 0x35, 0x63, 0x30, 0x62, 0x65, 0x33, 0x64, 0x34}),
		string([]byte{0x39, 0x36, 0x66, 0x63, 0x63, 0x33, 0x33, 0x33, 0x33, 0x61, 0x39, 0x61, 0x30, 0x33, 0x37, 0x66, 0x63, 0x65, 0x35, 0x31, 0x65, 0x63, 0x33, 0x62, 0x37, 0x62, 0x34, 0x37, 0x66, 0x66, 0x62, 0x37}))

	sk, err := storage.ReadLastFM(db)
	if err != nil {

		token, err := lmclient.GetToken()
		if err != nil {
			return nil, nil, "", err
		}

		if err := lmclient.LoginWithToken(token); err != nil {
			sk = "None"
		} else {
			sk = lmclient.Api.GetSessionKey()
		}

		err = storage.WriteLastFM([]byte(sk), db)
		if err != nil {
			return nil, nil, "", err
		}
	}

	if sk != "None" {
		lmclient.Api.SetSession(sk)
	}

	return gm, lmclient, sk, nil
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
