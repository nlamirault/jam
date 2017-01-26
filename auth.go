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
	return filepath.Join(os.Getenv("HOME"), ".jamdb")
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
