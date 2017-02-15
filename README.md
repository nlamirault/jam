# jam

Master : [![Circle CI](https://circleci.com/gh/nlamirault/jam/tree/ci.svg?style=svg)](https://circleci.com/gh/nlamirault/jam/tree/ci)

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

| Key                                 | Action                                                                       |
|-------------------------------------|------------------------------------------------------------------------------|
| <kbd>return</kbd>, <kbd>x</kbd>     | play currently selected artist, album or song                                |
| <kbd>c</kbd>                        | pause                                                                        |
| <kbd>v</kbd>                        | stop                                                                         |
| <kbd>b</kbd>                        | next track                                                                   |
| <kbd>z</kbd>                        | previous track                                                               |
| <kbd>u</kbd>                        | synchronize the database (in case you added some songs in the web interface) |
| <kbd>/</kbd>                        | search artists                                                               |
| <kbd>n</kbd>                        | next search result                                                           |
| <kbd>tab</kbd>                      | toggle artists/tracks view                                                   |
| <kbd>escape</kbd>, <kbd>q</kbd>     | quit                                                                         |
| <kbd>up arrow</kbd>, <kbd>k</kbd>   | scroll up                                                                    |
| <kbd>down arrow</kbd>, <kbd>j</kbd> | scroll down                                                                  |
| <kbd>Home</kbd>                     | scroll to top                                                                |
| <kbd>End</kbd>                      | scroll to bottom                                                             |
| <kbd>space</kbd>                    | toggle albums                                                                |
| <kbd>R</kbd>                        | randomize artists                                                            |



[1]: https://github.com/mjibson/moggio


TODO
- make the interface detachable (like MOC)
- make the binary able to receive comand line arguments for controlling playback
  (next track, pause, etc)
- implement search within the GPM global database
- feature requests are welcome as well
