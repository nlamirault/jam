# jam

This is my first Go program, I wanted to listen to Google Play Music on console,
so I wrote a player. It is inspired by Matt Jibson's [Moggio](https://github.com/mjibson/moggio/) and uses one of
his libraries. You can see it in action if you follow this link:
https://dl.dropboxusercontent.com/u/3651269/out-2.ogv

The features it has are:

- Last.fm scrobbling (use -lastfm flag)
- populating a local database with the artists and albums you saved through the
  web interface (or by any other means)
- searching within artists in the database
- playing, pausing (buggy, I need help with it) , stopping, previous track, next
  track
- the interface is Cmus rip off, I've only added a progress bar
- this player no longer lists artists in random order - if you want to randomize
  them press R


If you use 2-factor authorisation with your Google account, you will have to
generate an app password, follow this link 
https://security.google.com/settings/security/apppasswords

The linux binary I release is not static, it depends on pulseaudio, if you want
to build it from source, you are going to need the pulseaudio development package
installed.
Windows users are all set



If you have an x86 system, you'll have to compile it yourself, sorry

Contributions are welcome!

The keybindins are mostly the same as in Cmus:

| Key           | Action                                                                       |
|---------------|------------------------------------------------------------------------------|
| return, x     | play currently selected artist, album or song                                |
| c             | pause                                                                        |
| v             | stop                                                                         |
| b             | next track                                                                   |
| z             | previous track                                                               |
| u             | synchronize the database (in case you added some songs in the web interface) |
| /             | search artists                                                               |
| n             | next search result                                                           |
| tab           | toggle artists/tracks view                                                   |
| escape, q     | quit                                                                         |
| up arrow, k   | scroll up                                                                    |
| down arrow, j | scroll down                                                                  |
| Home          | scroll to top                                                                |
| End           | scroll to bottom                                                             |
| space         | toggle albums                                                                |
| R             | randomize artists                                                            |

[1]: https://github.com/mjibson/moggio



TODO
- make the interface detachable (like MOC)
- make the binary able to receive comand line arguments for controlling playback
  (next track, pause, etc)
- implement search within the GPM global database
- feature requests are welcome as well

