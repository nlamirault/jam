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

package lastfm

import (
	"time"

	"github.com/shkh/lastfm-go/lastfm"
)

// Client provides a Last.FM API client.
type Client struct {
	Api      *lastfm.Api
	Username string
	Password string
}

func New(apiKey string, apiSecret string) *Client {
	return &Client{
		Api: lastfm.New(apiKey, apiSecret),
	}
}

func (client *Client) Login(username string, password string) error {
	return client.Api.Login(username, password)
}

func (client *Client) Scrobble(artist string, track string) error {
	p := lastfm.P{"artist": artist, "track": track}
	_, err := client.Api.Track.UpdateNowPlaying(p)
	if err != nil {
		return err
	}
	// log.Printf("Now-Playing.")
	start := time.Now().Unix()
	time.Sleep(35 * time.Second)
	p["timestamp"] = start
	_, err = client.Api.Track.Scrobble(p)
	if err != nil {
		return err
	}
	// log.Printf("Scrobbled.")
	return nil
}
