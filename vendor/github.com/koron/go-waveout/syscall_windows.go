package waveout

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zsyscall_windows.go syscall_windows.go

//sys   Open(handle *syscall.Handle, deviceID uint32, waveFormat *WaveFormatEx, callback uintptr, inst uint32, flag uint32) (result MMRESULT) = winmm.waveOutOpen
//sys   Close(handle syscall.Handle) (result MMRESULT) = winmm.waveOutClose
//sys   GetVolume(handle syscall.Handle, volume *uint32) (result MMRESULT) = winmm.waveOutGetVolume
//sys   SetVolume(handle syscall.Handle, volume uint32) (result MMRESULT) = winmm.waveOutSetVolume
//sys   PrepareHeader(handle syscall.Handle, header *WaveHdr, size uint32) (result MMRESULT) = winmm.waveOutPrepareHeader
//sys   UnprepareHeader(handle syscall.Handle, header *WaveHdr, size uint32) (result MMRESULT) = winmm.waveOutUnprepareHeader
//sys   Write(handle syscall.Handle, header *WaveHdr, size uint32) (result MMRESULT) = winmm.waveOutWrite
//sys   Pause(handle syscall.Handle) (result MMRESULT) = winmm.waveOutPause
//sys   Restart(handle syscall.Handle) (result MMRESULT) = winmm.waveOutRestart
//sys   Reset(handle syscall.Handle) (result MMRESULT) = winmm.waveOutReset
//sys   BreakLoop(handle syscall.Handle) (result MMRESULT) = winmm.waveOutBreakLoop
