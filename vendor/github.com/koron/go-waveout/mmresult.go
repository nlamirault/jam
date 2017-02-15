// +build windows

package waveout

import "fmt"

type MMRESULT uint32

const (
	MMSYSERR_NOERROR      = MMRESULT(0)
	MMSYSERR_ERROR        = MMRESULT(1)
	MMSYSERR_BADDEVICEID  = MMRESULT(2)
	MMSYSERR_NOTENABLED   = MMRESULT(3)
	MMSYSERR_ALLOCATED    = MMRESULT(4)
	MMSYSERR_INVALHANDLE  = MMRESULT(5)
	MMSYSERR_NODRIVER     = MMRESULT(6)
	MMSYSERR_NOMEM        = MMRESULT(7)
	MMSYSERR_NOTSUPPORTED = MMRESULT(8)
	MMSYSERR_BADERRNUM    = MMRESULT(9)
	MMSYSERR_INVALFLAG    = MMRESULT(10)
	MMSYSERR_INVALPARAM   = MMRESULT(11)
	MMSYSERR_HANDLEBUSY   = MMRESULT(12)
	MMSYSERR_INVALIDALIAS = MMRESULT(13)
	MMSYSERR_BADDB        = MMRESULT(14)
	MMSYSERR_KEYNOTFOUND  = MMRESULT(15)
	MMSYSERR_READERROR    = MMRESULT(16)
	MMSYSERR_WRITEERROR   = MMRESULT(17)
	MMSYSERR_DELETEERROR  = MMRESULT(18)
	MMSYSERR_VALNOTFOUND  = MMRESULT(19)
	MMSYSERR_NODRIVERCB   = MMRESULT(20)
	MMSYSERR_MOREDATA     = MMRESULT(21)
	MMSYSERR_LASTERROR    = MMRESULT(21)

	WAVERR_BADFORMAT    = MMRESULT(32)
	WAVERR_STILLPLAYING = MMRESULT(33)
	WAVERR_UNPREPARED   = MMRESULT(34)
	WAVERR_SYNC         = MMRESULT(35)
	WAVERR_LASTERROR    = MMRESULT(35)
)

func (r MMRESULT) Error() string {
	switch r {
	case MMSYSERR_NOERROR:
		return "no error"
	case MMSYSERR_ERROR:
		return "unspecified error"
	case MMSYSERR_BADDEVICEID:
		return "device ID out of range"
	case MMSYSERR_NOTENABLED:
		return "driver failed enable"
	case MMSYSERR_ALLOCATED:
		return "device already allocated"
	case MMSYSERR_INVALHANDLE:
		return "device handle is invalid"
	case MMSYSERR_NODRIVER:
		return "no device driver present"
	case MMSYSERR_NOMEM:
		return "memory allocation error"
	case MMSYSERR_NOTSUPPORTED:
		return "function isn't supported"
	case MMSYSERR_BADERRNUM:
		return "error value out of range"
	case MMSYSERR_INVALFLAG:
		return "invalid flag passed"
	case MMSYSERR_INVALPARAM:
		return "invalid parameter passed"
	case MMSYSERR_HANDLEBUSY:
		return "handle being used simultaneously on another thread"
	case MMSYSERR_INVALIDALIAS:
		return "specified alias not found"
	case MMSYSERR_BADDB:
		return "bad registry database"
	case MMSYSERR_KEYNOTFOUND:
		return "registry key not found"
	case MMSYSERR_READERROR:
		return "registry read error"
	case MMSYSERR_WRITEERROR:
		return "registry write error"
	case MMSYSERR_DELETEERROR:
		return "registry delete error"
	case MMSYSERR_VALNOTFOUND:
		return "registry value not found"
	case MMSYSERR_NODRIVERCB:
		return "driver does not call DriverCallback"
	case MMSYSERR_MOREDATA:
		return "more data to be returned"
	case WAVERR_BADFORMAT:
		return "unsupported wave format"
	case WAVERR_STILLPLAYING:
		return "still something playing"
	case WAVERR_UNPREPARED:
		return "header not prepared"
	case WAVERR_SYNC:
		return "device is synchronous"
	default:
		return fmt.Sprintf("unknown error: %d", r)
	}
}
