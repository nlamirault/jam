package waveout

type WaveFormatEx struct {
	FormatTag      uint16
	Channels       uint16
	SamplesPerSec  uint32
	AvgBytesPerSec uint32
	BlockAlign     uint16
	BitsPerSample  uint16
	Size           uint16
}

type WaveHdr struct {
	Data          *byte
	BufferLength  uint32
	BytesRecorded uint32
	User          uintptr
	Flags         uint32
	Loops         uint32
	Next          *WaveHdr
	Reserved      uintptr
}

const (
	WAVE_MAPPER = uint32(0xFFFFFFFF)

	CALLBACK_NULL     = uint32(0x00000000)
	CALLBACK_WINDOW   = uint32(0x00010000)
	CALLBACK_TASK     = uint32(0x00020000)
	CALLBACK_FUNCTION = uint32(0x00030000)
	CALLBACK_THREAD   = uint32(0x00020000)
	CALLBACK_EVENT    = uint32(0x00050000)

	WAVE_FORMAT_QUERY                        = uint32(0x00000001)
	WAVE_ALLOWSYNC                           = uint32(0x00000002)
	WAVE_MAPPED                              = uint32(0x00000004)
	WAVE_FORMAT_DIRECT                       = uint32(0x00000008)
	WAVE_FORMAT_DIRECT_QUERY                 = WAVE_FORMAT_QUERY | WAVE_FORMAT_DIRECT
	WAVE_MAPPED_DEFAULT_COMMUNICATION_DEVICE = uint32(0x00000010)

	WAVE_FORMAT_PCM = 0x0001

	WHDR_DONE      = uint32(0x00000001)
	WHDR_PREPARED  = uint32(0x00000002)
	WHDR_BEGINLOOP = uint32(0x00000004)
	WHDR_ENDLOOP   = uint32(0x00000008)
	WHDR_INQUEUE   = uint32(0x00000010)

	WOM_OPEN  = uint32(0x000003BB)
	WOM_CLOSE = uint32(0x000003BC)
	WOM_DONE  = uint32(0x000003BD)
)
