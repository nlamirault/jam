# jam

This is my first Go program, I wanted to listen to Google Play Music on console,
so I wrote a player. It is inspired by Matt Jibson's [Moggio](https://github.com/mjibson/moggio/) and uses one of
his libraries. You can see it in action if you follow this link:
https://dl.dropboxusercontent.com/u/3651269/out-2.ogv

The features it has are:

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

The binary I released is not static, it depends on the following libraries on my
system

```bash
 $ ldd ./jam 
	linux-vdso.so.1 (0x00007fff8e3e7000)
	libpulse-simple.so.0 => /usr/lib/libpulse-simple.so.0 (0x00007fe7eca1e000)
	libpulse.so.0 => /usr/lib/libpulse.so.0 (0x00007fe7ec7cd000)
	libpthread.so.0 => /usr/lib/libpthread.so.0 (0x00007fe7ec5b0000)
	libc.so.6 => /usr/lib/libc.so.6 (0x00007fe7ec212000)
	libpulsecommon-10.0.so => /usr/lib/pulseaudio/libpulsecommon-10.0.so (0x00007fe7ebf8d000)
	libdbus-1.so.3 => /usr/lib/libdbus-1.so.3 (0x00007fe7ebd3d000)
	libdl.so.2 => /usr/lib/libdl.so.2 (0x00007fe7ebb39000)
	libm.so.6 => /usr/lib/libm.so.6 (0x00007fe7eb835000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fe7ecc23000)
	libxcb.so.1 => /usr/lib/libxcb.so.1 (0x00007fe7eb60c000)
	libsystemd.so.0 => /usr/lib/libsystemd.so.0 (0x00007fe7ecd8f000)
	libsndfile.so.1 => /usr/lib/libsndfile.so.1 (0x00007fe7eb394000)
	libasyncns.so.0 => /usr/lib/libasyncns.so.0 (0x00007fe7eb18e000)
	librt.so.1 => /usr/lib/librt.so.1 (0x00007fe7eaf86000)
	libXau.so.6 => /usr/lib/libXau.so.6 (0x00007fe7ead82000)
	libXdmcp.so.6 => /usr/lib/libXdmcp.so.6 (0x00007fe7eab7c000)
	libresolv.so.2 => /usr/lib/libresolv.so.2 (0x00007fe7ea965000)
	libcap.so.2 => /usr/lib/libcap.so.2 (0x00007fe7ea761000)
	liblzma.so.5 => /usr/lib/liblzma.so.5 (0x00007fe7ea53b000)
	liblz4.so.1 => /usr/lib/liblz4.so.1 (0x00007fe7ea327000)
	libgcrypt.so.20 => /usr/lib/libgcrypt.so.20 (0x00007fe7ea018000)
	libgpg-error.so.0 => /usr/lib/libgpg-error.so.0 (0x00007fe7e9e04000)
	libFLAC.so.8 => /usr/lib/libFLAC.so.8 (0x00007fe7e9b8e000)
	libogg.so.0 => /usr/lib/libogg.so.0 (0x00007fe7e9987000)
	libvorbis.so.0 => /usr/lib/libvorbis.so.0 (0x00007fe7e975a000)
	libvorbisenc.so.2 => /usr/lib/libvorbisenc.so.2 (0x00007fe7e94a7000)
```

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
| space         | toggle albums                                                                |
| R             | randomize artists                                                            |

[1]: https://github.com/mjibson/moggio



TODO
- make the interface detachable (like MOC)
- make the binary able to receive comand line arguments for controlling playback
  (next track, pause, etc)
- implement search within the GPM global database
- feature requests are welcome as well

